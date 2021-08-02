import { createApp } from 'vue'
import ElementPlus from 'element-plus';
import 'element-plus/lib/theme-chalk/index.css';
import {tcliPlugin} from './js/tcli'
import App from './App.vue'


import { createStore } from 'vuex'


const token = window.localStorage.getItem('usertoken')
const store = createStore({
  state () {
    return {
      userid: window.localStorage.getItem('userid'),
      username: window.localStorage.getItem('username'),
      usertoken: token,
    }
  },
  mutations: {
    login (state, user) {
      state.userid = user.userid
      state.username = user.username
      state.usertoken = user.usertoken
      window.localStorage.setItem('userid', user.userid)
      window.localStorage.setItem('username', user.username)
      window.localStorage.setItem('usertoken', user.usertoken)
    },
    logout (state) {
      state.userid = null
      state.username = null
      state.usertoken = null
      window.localStorage.removeItem('userid')
      window.localStorage.removeItem('username')
      window.localStorage.removeItem('usertoken')
    }
  }
})

const app = createApp(App)
app.use(ElementPlus)
app.use(store)
app.use(tcliPlugin, token, process.env.VUE_APP_TESTING_ADMIN_BASEURL)
window.vm = app.mount('#app')
