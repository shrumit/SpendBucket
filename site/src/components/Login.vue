<template>
<div class="central">
	<mu-content-block class="intro">
		<p style="margin:0px; margin-bottom:-30px; text-align:right;">Beta</p>
		<h1>Spend Bucket</h1>
		<!-- <h4>Split expenses. Easily.</h4> -->
	</mu-content-block>

	<mu-paper class="card" :zDepth=2>
		<mu-tabs :value="activeTab" @change="handleTabChange">
			<mu-tab value="login" title="Login"/>
			<mu-tab value="register" title="Register"/>
			<mu-tab value="demo" title="Demo"/>
		</mu-tabs>

		<!-- Login -->
		<mu-content-block v-if="activeTab === 'login'">
			<mu-text-field v-model="username" label="Username" labelFloat fullWidth/>
			<mu-text-field v-model="password" type="password" @keyup.enter="login" label="Password" labelFloat fullWidth/>
			<!-- <mu-checkbox v-model="rememberMe" label="Remember?" /> -->
			<br />
			<mu-raised-button label="Login" @click="login" :disabled="!enableLogin" secondary fullWidth />
			<p>{{loginMsg}}</p>
		</mu-content-block>
		
		<!-- Register -->
		<mu-content-block v-if="activeTab === 'register'">
			<mu-text-field v-model="username" label="Username" labelFloat fullWidth/>
			<mu-text-field v-model="password" label="Password" type="password" labelFloat fullWidth/>
			<mu-text-field v-model="password2" label="Password" type="password" labelFloat fullWidth />
			<br />
			<mu-raised-button label="Register" @click="register" :disabled="!enableRegister" secondary fullWidth/>
			<p>{{registerMsg}}</p>
		</mu-content-block>
		
		<!-- Demo -->
		<mu-content-block v-if="activeTab === 'demo'">
			<p>Log into a dummy account.</p>
			<p>User data is reset to default every 15 minutes.</p>
			<mu-raised-button label="Enter Demo" @click="demo" secondary fullWidth/>
		</mu-content-block>
	</mu-paper>
</div>
</template>

<style lang="less" scoped>
@import '../assets/colors.less';

.intro {
	text-align: center;
}

.central {
	width: 95%;
	max-width: 300px;
	margin: auto auto;
	/*text-alig: center;*/
}

.mu-raised-button {
	margin: 10px auto;
}

p {
	margin-top: 5px;
	text-align: center;
	color: grey;
}

em {
	text-align: right;
}
</style>

<script>

	export default {
		name: 'login',
		data() {
			return {
				username: '',
				password: '',
				password2: '',
				rememberMe: true,
				activeTab: 'login',
				// loginMsg: '',
				// registerMsg: ''
			}
		},
		computed: {
			enableLogin: function() {
				return (this.username !== '' && this.password !== '')
			},
			enableRegister: function() {
				return ((this.username.length >= 4) && (this.password.length >= 6) && (this.password === this.password2))
			},
			registerMsg: function() {
				if (this.username.length < 4)
					return 'Username too short.'
				else if (this.password.length < 6)
					return 'Password too short.'
				else if (this.password !== this.password2)
					return 'Passwords don\'t match.'
				else
					return 'All ok!'
			},
			loginMsg: function() {
				if (this.username.length === 0 || this.password.length === 0)
					return 'Please enter credentials.'
				else
					return ''
			}
		},
	// props: ['count'],	
	methods: {
		login: function() {
			this.msg = ''
			if (!this.enableLogin)
				return

			this.$store.dispatch('login', {username: this.username, password: this.password})
			this.msg = "Invalid credentials."
		},
		register: function() {
			if (!this.enableRegister)
				return

			this.$store.dispatch('register', {username: this.username, password: this.password})
		},
		demo () {
			this.$store.dispatch('login', {username: 'dummy@example.com', password: 'Spendpass123'})
		},
		handleTabChange (val) {
			this.msg = ''
			this.activeTab = val
		}
	},
	beforeCreate: function() {
		this.$store.dispatch('initialize')
	}

}
</script>

