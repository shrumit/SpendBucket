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

		<mu-table :selectable=false :showCheckbox=false class="section-table">
			<mu-thead>
				<mu-tr>
					<mu-th>Title</mu-th>
					<mu-th>Amount</mu-th>
					<mu-th>Date</mu-th>
					<mu-th>Paid By</mu-th>
					<mu-th>Shared By</mu-th>
					<mu-th class="mu-checkbox-col"></mu-th>
				</mu-tr>
			</mu-thead>
			<mu-tbody>
				<mu-tr v-for="t in trans" :key="t.transId">
					<mu-td>{{t.title}}</mu-td>
					<mu-td>{{t.amount}}</mu-td>
					<mu-td>{{t.transDate}}</mu-td>
					<mu-td>{{t.paidBy_display}}</mu-td>
					<mu-td>{{t.sharedBy_display}}</mu-td>
					<mu-td style="padding:0px;">
						<mu-icon-menu icon="edit" :anchorOrigin="leftTop" :targetOrigin="leftTop">
							<!-- <mu-menu-item @click="editTransaction(t.transId)" title="Edit" /> -->
							<mu-menu-item @click="deleteTransaction(t.transId)" title="Delete" />
						</mu-icon-menu>
					</mu-td>
				</mu-tr>
			</mu-tbody>
		</mu-table>

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