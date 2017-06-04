<template>
<div>
	<mu-paper :zDepth=1>

		<div class="section-header">
			<h2>Transactions</h2>
			<mu-divider/>
		</div>
		<div class="section-add">
			<mu-float-button :disabled="persons.length === 0" @click="openDialog" icon="add" secondary/>
		</div>

 		<!-- transaction table -->
		<md-table class="section-table">
			<md-table-header>
				<md-table-row>
					<md-table-head>Title</md-table-head>
					<md-table-head>Amount</md-table-head>
					<md-table-head>Date</md-table-head>
					<md-table-head>Paid By</md-table-head>
					<md-table-head>Shared By</md-table-head>
				</md-table-row>
			</md-table-header>
			<md-table-body>
				<md-table-row v-for="t in trans" :key="t.transId">
					<md-table-cell>{{t.title}}</md-table-cell>
					<md-table-cell md-numeric>{{t.amount}}</md-table-cell>
					<md-table-cell>{{t.transDate}}</md-table-cell>
					<md-table-cell>{{t.paidBy_display}}</md-table-cell>
					<md-table-cell>{{t.sharedBy_display}}</md-table-cell>
					<mu-icon-menu icon="edit" :anchorOrigin="leftTop" :targetOrigin="leftTop">
						<mu-menu-item @click="deleteTransaction(t.transId)" title="Delete" />
					</mu-icon-menu>
				</md-table-row>
			</md-table-body>
		</md-table>

		<!-- new transaction dialog -->
		<mu-dialog :open="dialog" @close="closeDialog">
			<h3>New Transaction</h3>
			<mu-text-field v-model.trim="newTran.title" label="Title" labelFloat/> <br/>
			<mu-text-field v-model.trim="newTran.amount" label="Amount" :errorText="newAmountErrorText" labelFloat/> <br/>
			<mu-date-picker v-model.trim="newTran.transDate" label="Date" :dateTimeFormat="enDateFormat" labelFloat mode="landscape" container="inline" autoOk/> <br/>

			<mu-select-field v-model="newTran.paidBy" label="Paid By" labelFloat>
				<mu-menu-item v-for="p in persons" :key="p.personId" :title=p.personName :value="p.personId" />
			</mu-select-field>
			<br/>
			<mu-select-field v-model="newTran.sharedBy" multiple label="Shared By" labelFloat>
				<mu-menu-item v-for="p in persons" :key="p.personId" :title=p.personName :value="p.personId" />
			</mu-select-field>

			<div style="text-align: center;">
				<mu-raised-button label="Submit" @click="newTransaction" :disabled="!newValid" secondary/>
				<!-- <mu-flat-button label="Cancel" @click="closeDialog" id="newCloseButton"/> -->
			</div>
 		</mu-dialog>


	</mu-paper>

</div>
</template>

<script>

// import Toolbar from './Toolbar.vue'

const dayAbbreviation = ['S', 'M', 'T', 'W', 'T', 'F', 'S']
const dayList = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
const monthList = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep',
  'Oct', 'Nov', 'Dec']
const monthLongList = ['January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December']

const enDateFormat = {
  formatDisplay (date) {
    return `${dayList[date.getDay()]}, ${monthList[date.getMonth()]} ${date.getDate()}`
  },
  formatMonth (date) {
    return `${monthLongList[date.getMonth()]} ${date.getFullYear()}`
  },
  getWeekDayArray (firstDayOfWeek) {
    let beforeArray = []
    let afterArray = []
    for (let i = 0; i < dayAbbreviation.length; i++) {
      if (i < firstDayOfWeek) {
        afterArray.push(dayAbbreviation[i])
      } else {
        beforeArray.push(dayAbbreviation[i])
      }
    }
    return beforeArray.concat(afterArray)
  }
}

export default {
	name: 'transactions',
	// components: { Toolbar },
	data() {
		return {
			enDateFormat,
			dialog: false,
			leftTop: {horizontal: 'left', vertical: 'top'},

			newTran: {
				transId: 0,
				title: '',
				amount: '',
				transDate: '',
				paidBy: '',
				sharedBy: []			
			}
		}
	},
	computed: {
		persons() {
			return this.$store.state.selectedGroupData.persons
		},
		personsMap() {
			let m = {}
			this.persons.forEach(function(p){
				m[p.personId] = p.personName
			})
			return m
		},
		trans() {
			// return this.$store.state.selectedGroupData.transactions
			let trans = this.$store.state.selectedGroupData.transactions
			let pMap = this.personsMap

			trans = trans.map(function(t) {
				t.sharedBy_display = ''
				t.sharedBy.forEach(function(s){
					t.sharedBy_display += pMap[s] + ', '
				})
				t.sharedBy_display = t.sharedBy_display.slice(0,-2);
				t.paidBy_display = pMap[t.paidBy]
				t.amount = parseFloat(t.amount).toFixed(2)
				return t
			})
			// console.log(transactions)
			return trans
		},
		newValid() {
			if ((this.newTran.title === '') || 
				(this.newTran.amount === '') || 
				(this.newTran.transDate === '') || 
				(this.newTran.paidBy === '') || 
				(this.newTran.sharedBy.length === 0))
				return false

			if (isNaN(this.newTran.amount)) {
				return false
			}

			return true
		},
		newAmountErrorText() {
			if (isNaN(this.newTran.amount))
				return 'Please enter a valid number!'
			else
				return ''
		}
	},
	methods: {
		openDialog() {
			this.dialog = true
		},

		closeDialog() {
			this.dialog = false
		},

		newTransaction() {
			if (!this.newValid)
				return
			
			// this is to correct for a bug in muse-ui when closing the overlay
			setTimeout(this.closeDialog, 100)

			this.$store.dispatch('addTransaction', this.newTran)
			this.newTran = {
				title: '',
				amount: '',
				transDate: '',
				paidBy: '',
				sharedBy: []			
			}
		},

		deleteTransaction(val) {
			console.log(val)
			this.$store.dispatch('deleteTransaction', val)
		}
	}
}
</script>

<style scoped>

</style>