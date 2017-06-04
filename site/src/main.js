import Vue from 'vue'
// import Vuex from 'vuex'

import MuseUI from 'muse-ui'
import 'muse-ui/dist/muse-ui.css'
import 'muse-ui/dist/theme-carbon.css'
Vue.use(MuseUI)

import css from 'vue-material/dist/vue-material.css'
import VueMaterial from 'vue-material'
Vue.use(VueMaterial)

import {store} from './store.js'
import App from './App.vue'

new Vue({
  el: '#app',
  store,
  render: h => h(App)
})