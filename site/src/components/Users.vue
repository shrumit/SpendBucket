<template>
<div>
	<mu-paper :zDepth="2">

		<div class="section-header">
			<h2>Users in this Group</h2>
			<mu-divider/>
		</div>

		<!-- username table -->
		<mu-table :selectable=false :showCheckbox=false >
			<mu-thead>
				<mu-tr>
					<mu-th>Username</mu-th>
				</mu-tr>
			</mu-thead>
			<mu-tbody>
				<mu-tr v-for="u in usernames" :key="u">
					<mu-td>{{ u }}</mu-td>
				</mu-tr>
			</mu-tbody>
		</mu-table>
	</mu-paper>
</div>
</template>

<script>
export default {
	name: 'users',
	data() {
		return {

		}
	},
	computed: {
		usernames() {
			return this.$store.state.selectedGroupData.users
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