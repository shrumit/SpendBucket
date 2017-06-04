<template>
<div>
	<mu-paper :zDepth=1 style="max-width: 400px;">

		<div class="section-header">
			<h2>Persons</h2>
			<mu-divider/>
		</div>
		<div class="section-add">
			<mu-float-button @click="openDialog" icon="add" secondary/>
		</div>

		<!-- new person dialog -->
		<mu-dialog :open="dialog" @close="closeDialog">
			<h3>New Person</h3>
			<mu-text-field v-model.trim="newName" label="Name" labelFloat/> <br/>
			<p v-if="newIsDuplicate">Duplicate names not allowed</p>
			<div style="text-align: center;">
				<mu-raised-button label="Submit" @click="newPerson" :disabled="!newValid" secondary/>
				<!-- <mu-flat-button label="Cancel" @click="closeDialog" id="newCloseButton"/> -->
			</div>
 		</mu-dialog>

 		<!-- persons table -->
		<mu-table :selectable=false :showCheckbox=false class="section-table">
			<mu-thead>
				<mu-tr>
					<mu-th>Person</mu-th>
					<mu-th>Balance</mu-th>
				</mu-tr>
			</mu-thead>
			<mu-tbody>
				<mu-tr v-for="p in persons" :key="p.personId">
					<mu-td>{{p.personName}}</mu-td>
					<mu-td>{{p.balanceText}}</mu-td>
				</mu-tr>
			</mu-tbody>
		</mu-table>

	</mu-paper>
</div>
</template>

<script>
export default {
	name: 'persons',
	data() {
		return {
			dialog: false,
			newName : '',
			newMsg : ''
		}
	},
	computed: {
		persons() {

			return this.$store.state.selectedGroupData.persons.map(function(elem){
				elem.balance = parseFloat(elem.balance)
				if (elem.balance < 0) {
					elem.balanceText = 'Owes $' + (elem.balance*-1).toFixed(2)
				} else {
					elem.balanceText = 'Is owed $' + elem.balance.toFixed(2)
				}
				return elem
			})
		},

		newIsDuplicate() {
			let n = this.newName
			let i = this.persons.findIndex(function(obj){
				return obj.personName === n
			})
			return (i > -1)
		},

		newValid() {
			return (this.newName && !this.newIsDuplicate)
		}
	},
	methods: {
		openDialog() {
			this.dialog = true
		},

		closeDialog() {
			this.dialog = false
		},

		newPerson() {
			if (!this.newValid)
				return
			
			// delay is to correct for a bug in muse-ui when closing the overlay
			setTimeout(this.closeDialog, 100)

			this.$store.dispatch('addPerson', this.newName)
			this.newName = ''
		}
	}
}
</script>

<style scoped>
/*.mu-paper {
	padding: 5px 0px;
	margin: 10px;
}*/
</style>