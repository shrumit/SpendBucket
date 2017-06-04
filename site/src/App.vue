<template>
	<div id="app">
		<template v-if="loggedIn">
			<dashboard></dashboard>
		</template>

		<template v-else>
			<login></login>
		</template>
		<p>{{counter}}</p>

	</div>
</template>

<script>
import Login from './components/Login.vue'
import Dashboard from './components/Dashboard.vue'

export default {
	name: 'app',
	components: { Login, Dashboard},
	computed: {
		counter() {
			return this.$store.state.counter
		},
		loggedIn() {
			return this.$store.state.loggedIn	
		},
		token() {
			return this.$store.state.token	
		}
	}
}
</script>

<style lang="less">
body {
	background-color: #ebf1e4;
}

#app {
	font-family: 'Roboto', sans-serif;
}

h1 {
	font-family: 'Roboto Slab', serif;
	font-weight: normal;
	font-size: 3rem;
}

h2 {
	font-family: 'Roboto Slab', serif;
	font-size: 2.2rem;
	font-weight: normal;
	margin: 10px 0;
	margin-left: 20px;	
}

.mu-appbar {
	margin-bottom:20px;
}

.mu-appbar-title {
	font-family: 'Roboto Slab', serif;
}

/* Transaction and Person paper header spacing */
.section-header {
	margin-bottom: -28px;
	text-align:center;
}
.section-add {
	text-align:right;
	margin-right: 5%;
}
.section-table {
	margin-top: -15px;	
}

/* get rid of calendar controls */
.mu-calendar-landspace .mu-date-display {
	height:290px;
}
.mu-calendar-actions {
	display:none;
}

/*dialog shrink to fit*/
.mu-dialog {
	width: auto;
}

/* Dashboard stuff spacing */
.dash-content {
	max-width: 1200px;
	margin: 0 auto;
}

.dash-content .mu-paper {
	padding: 5px 0px;
	margin: 10px;
}

/* dash information text */

.dash-info {
	text-align: center;
	margin: 0 20px;
	white-space:pre-wrap;
}

</style>