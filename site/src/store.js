import Vue from 'vue'
import Vuex from 'vuex'
var request = require('superagent')

Vue.use(Vuex)

const BASE_URL = (process.env.NODE_ENV === 'production') ? '' : 'http://localhost:8081'
const STORAGE_KEY = 'spendbucket_token'

var persist = false

var tokenStorage = {
	fetch: function() {
		return localStorage.getItem(STORAGE_KEY)
	},
	save: function(token) {
		localStorage.setItem(STORAGE_KEY, token)
	},
	delete: function() {
		localStorage.removeItem(STORAGE_KEY)
	}
}

export const store = new Vuex.Store({
	state: {
		loggedIn: false,
		groupsMetadata: [],
		selectedGroupId: 0,
		selectedGroupData: {
			persons: [],
			transactions: [],
			users: []
		}
	},

	mutations: {
		loggedIn(state, payload) {
			state.loggedIn = payload
		},
		groupsMetadata(state, payload) {
			state.groupsMetadata = payload
		},
		groupsMetadata_add(state, payload) {
			state.groupsMetadata.push(payload)
		},
		selectedGroupId(state, payload) {
			state.selectedGroupId = payload
		},
		selectedGroupData_persons(state, payload) {
			state.selectedGroupData.persons = payload
		},
		selectedGroupData_transactions(state, payload) {
			state.selectedGroupData.transactions = payload
		},
		selectedGroupData_users(state, payload) {
			state.selectedGroupData.users = payload
		},
		selectedGroupData_transactions_delete(state, payload) {
			let i = state.selectedGroupData.transactions.findIndex(function(obj){
				return obj.transId === payload
			})
			if (i > -1) {
				state.selectedGroupData.transactions.splice(i, 1)
			}
		},
		resetToDefaults(state) {
			state.loggedIn = false
			state.groupsMetadata = []
			state.selectedGroupId = 0
			state.selectedGroupData = {
				persons: [],
				transactions: [],
				users: []
			}
		}
	},

	actions: {
		initialize({ commit, dispatch, state }) {
			let token = tokenStorage.fetch()
			if (token === null)
				return

			request.post(BASE_URL+'/getGroups')
			.type('form')
			.send({token: token})
			.end(function (err, res) {
				if (!err && res.body.success) {
					if (res.body.groups != null) {
						commit('groupsMetadata', res.body.groups)
						commit('selectedGroupId', res.body.groups[0].groupId)
						dispatch('updateGroupData')
					}
					commit('loggedIn', true) // this positioning here is important
				}
			})
		},
		login({ dispatch }, payload) {
			// authenticate and get token
			request.post(BASE_URL+'/login')
			.type('form')
			.send(payload)
			.end(function (err, res) {
				if (!err && res.body.success) {
					tokenStorage.save(res.body.token)
					dispatch('initialize')
				}
			})
		},
		logout({ commit }){
			tokenStorage.delete()
			commit('resetToDefaults')
		},
		register({ dispatch }, payload) {
			request.post(BASE_URL+'/register')
			.type('form')
			.send(payload)
			.end(function (err, res) {
				if (!err && res.body.success) {
					tokenStorage.save(res.body.token)
					dispatch('initialize')
				}
			})			
		},
		changeSelectedGroup({ commit, dispatch }, groupId) {
			commit('selectedGroupId', groupId)
			dispatch('updateGroupData')
		},
		addTransaction({ dispatch, state }, transaction) {
			let token = tokenStorage.fetch()
			let data = {token: token, groupId: state.selectedGroupId, transaction: transaction}

			request.post(BASE_URL+'/addTransaction')
			.type('form')
			.send(data)
			.end(function (err, res) {
				if (!err && res.body.success) {
					dispatch('updateTransactions')
					dispatch('updatePersons')
				}
			})
		},
		deleteTransaction({ commit, dispatch, state }, transId) {
			request.post(BASE_URL+'/deleteTransaction')
			.type('form')
			.send({token: tokenStorage.fetch(), transId: transId, groupId: state.selectedGroupId})
			.end(function (err, res) {
				if (!err && res.body.success) {
					commit('selectedGroupData_transactions_delete', transId)
					dispatch('updatePersons')
				}
			})
		},
		addPerson({ dispatch, state }, personName) {
			request.post(BASE_URL+'/addPerson')
			.type('form')
			.send({token: tokenStorage.fetch(), groupId: state.selectedGroupId, personName: personName})
			.end(function (err, res) {
				if (!err && res.body.success) {
					dispatch('updatePersons')
				}
			})
		},
		makeGroup({ commit, dispatch }, groupName) {
			request.post(BASE_URL+'/createGroup')
			.type('form')
			.send({token: tokenStorage.fetch(), groupName: groupName})
			.end(function (err, res) {
				if (!err && res.body.success) {
					commit('groupsMetadata_add', res.body.group)
					commit('selectedGroupId', res.body.group.groupId)
					dispatch('updateGroupData')
				}
			})

		},
		joinGroup({ commit, dispatch }, inviteCode) {
			request.post(BASE_URL+'/enterGroupInvite')
			.type('form')
			.send({token: tokenStorage.fetch(), inviteCode: inviteCode})
			.end(function (err, res) {
				if (!err && res.body.success) {
					commit('groupsMetadata_add', res.body.group)
					commit('selectedGroupId', res.body.group.groupId)
					dispatch('updateGroupData')
				}
			})
		},

		// TODO: Make API method for POST requests (payload, endpoint, sendGid)

		/* Fetch group data */

		updateGroupData({ dispatch }) {		
			dispatch('updateTransactions')
			dispatch('updatePersons')
			dispatch('updateUsers')
		},
		updateTransactions({ commit, dispatch, state }) {
			request.post(BASE_URL+'/getTransactions')
			.type('form')
			.send({token: tokenStorage.fetch(), groupId: state.selectedGroupId})
			.end(function (err, res) {
				if (!err && res.body.success) {
					res.body.transactions = (res.body.transactions === null ? [] : res.body.transactions)
					commit('selectedGroupData_transactions', res.body.transactions)
				} else {
					dispatch('initialize')
				}
			})
		},
		updatePersons({ commit, state }) {
			request.post(BASE_URL+'/getPersons')
			.type('form')
			.send({token: tokenStorage.fetch(), groupId: state.selectedGroupId})
			.end(function (err, res) {
				if (!err && res.body.success) {
					res.body.persons = (res.body.persons === null ? [] : res.body.persons)
					commit('selectedGroupData_persons', res.body.persons)
				}
			})
		},
		updateUsers({ commit, state }) {
			request.post(BASE_URL+'/getUsernames')
			.type('form')
			.send({token: tokenStorage.fetch(), groupId: state.selectedGroupId})
			.end(function (err, res) {
				if (!err && res.body.success) {
					res.body.usernames = (res.body.usernames === null ? [] : res.body.usernames)
					commit('selectedGroupData_users', res.body.usernames)
				}
			})			
		}
	}
})