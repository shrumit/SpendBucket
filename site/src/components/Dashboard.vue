<template>
<div>
	<toolbar></toolbar>

	<div class="dash-info">
		<h3>{{ infoText }}</h3>
	</div>

	<div v-if="!noGroups">
	<mu-row class="dash-content">
		<mu-col width="100" tablet="100" desktop="70">
			<transactions></transactions>
		</mu-col>
		<mu-col width="100" tablet="100" desktop="30">
			<persons></persons>
		</mu-col>
	</mu-row>
	</div>
</div>
</template>

<script>

import Toolbar from './Toolbar.vue'
import Transactions from './Transactions.vue'
import Persons from './Persons.vue'

export default {
	name: 'dashboard',
	components: { Toolbar, Transactions, Persons },
	data() {
		return {}
	},
	computed: {
		noGroups() {
			return (this.$store.state.groupsMetadata.length === 0)
		},
		infoText() {
			if (this.noGroups) {
				return 'Welcome to Spendbucket!\nBegin by making a new group or joining an existing group using the menu on the left.'
			} else if (this.$store.state.selectedGroupData.persons.length === 0) {
				return 'Start by making a new person. Persons can take part in transactions.'
			}
			else if (this.$store.state.selectedGroupData.transactions.length === 0)
				return 'Excellent! You can add more persons or add your first transaction.'
			else
				return ''

		}
	},
	methods: {
	}
}
</script>

<style>

</style>