<template>
	<div>
		<mu-appbar :zDepth=1>
			<mu-icon-button icon="menu" slot="left" @click="toggleDrawer"/>
			<h2 v-if="selectedGroupData">{{ selectedGroupData.groupName }}</h2>
			<h2 slot="right">
				<span class="invite-title">Invite Code: </span>{{selectedGroupData.inviteCode}}
			</h2>
		</mu-appbar>

		<!-- Drawer -->
		<mu-drawer :open="openDrawer" @close="toggleDrawer" class="drawer" :docked=false :zDepth=2>
			<mu-appbar title="SpendBucket" />
			<mu-list>

			<mu-sub-header>Users in this group</mu-sub-header>
			<mu-list-item v-for="u in usernames" :key="u" :title="u" disabled>
				<mu-icon slot="left" value="person"/>
			</mu-list-item>

			<mu-divider />
			<mu-sub-header>Your groups</mu-sub-header>
				<mu-list-item v-for="group in groupsMetadata" :key="group.groupId" :title="group.groupName" @click="changeSelectedGroup(group.groupId)" />
			
			<mu-divider />
			<mu-sub-header>Group actions</mu-sub-header>
			<mu-list-item title="Join a group" @click="toggleInvite">
				<mu-icon slot="left" value="person_add"/>
			</mu-list-item>
			<mu-list-item title="Make a new group" @click="toggleNew">
				<mu-icon slot="left" value="create_new_folder"/>
			</mu-list-item>
			
			<mu-divider />
			<mu-list-item title="Sign out" @click="logout">
			</mu-list-item>
			</mu-list>
		</mu-drawer>

		<!-- Enter invite dialog -->
		<mu-dialog :open="openInvite" @close="toggleInvite">
			<h3>Enter Invite</h3>
			<mu-text-field v-model.trim="invite" label="Invite Code" labelFloat/> <br/>
			<div style="text-align: center;">
				<mu-raised-button label="Ok" @click="joinGroup" :disabled="!inviteValid" secondary/>
			</div>
		</mu-dialog>

		<!-- Make new group dialog -->
		<mu-dialog :open="openNew" @close="toggleNew">
			<h3>New Group</h3>
			<mu-text-field v-model.trim="groupName" label="Group Name" labelFloat/> <br/>
			<div style="text-align: center;">
				<mu-raised-button label="Ok" @click="makeGroup" :disabled="!groupNameValid" secondary/>
			</div>
		</mu-dialog>


	</div>
</template>

<script>
export default {
	name: 'toolbar',
	data() {
		return {
			// this is in data so that only initial value is taken from store
			openDrawer: false,
			openInvite: false,
			openNew: false,
			invite: '',
			groupName: ''
		}
	},

	computed: {
		loggedIn() {
			return this.$store.state.loggedIn
		},
		groupsMetadata() {
			return this.$store.state.groupsMetadata
		},
		noGroups() {
			return (this.$store.state.groupsMetadata.length === 0)
		},
		selectedGroupData() {
			if (this.noGroups)
				return {groupName: '', inviteCode: ''}

			let id = this.$store.state.selectedGroupId
			return this.$store.state.groupsMetadata.find(function(elem) {
				if (elem.groupId === id)
					return elem
			})
		},
		inviteValid() {
			return (this.invite !== '')
		},
		groupNameValid() {
			return (this.groupName !== '')
		},
		usernames() {
			return this.$store.state.selectedGroupData.users
		}
	},

	methods: {
		// resetData() {
		// 	this.openDrawer = false
		// 	this.openInvite = false
		// 	this.openNew = false
		// 	this.invite = ''
		// 	this.groupName = ''
		// },
		changeSelectedGroup(val) {
			this.openDrawer = false
			this.$store.dispatch('changeSelectedGroup', val)
		},
		toggleDrawer() {
			this.openDrawer = !this.openDrawer
		},
		toggleInvite() {
			this.openInvite = !this.openInvite
		},
		toggleNew() {
			this.openNew = !this.openNew
		},
		joinGroup() {
			if (!this.inviteValid)
				return
			this.$store.dispatch('joinGroup', this.invite)
			this.openInvite = false
			this.invite = ''
		},
		makeGroup() {
			if (!this.groupNameValid)
				return
			this.$store.dispatch('makeGroup', this.groupName)
			this.openNew = false
			this.groupName = ''
		},
		logout() {
			// this.resetData()
			this.$store.dispatch('logout')
		}
	},

	mounted() {
		this.localSelectedGroupId = this.$store.state.selectedGroupId
	}
}
</script>

<style>
.drawer .mu-divider {
	margin-top: 30px;
}

.invite-title {
	color: grey;
}
</style>