import { createStore } from 'vuex'
import app from './app'
import user from './modules/user'
import permission from './modules/permission'
import settings from './settings'
import getters from './getters'

const store = createStore({
  modules: {
    app,
    user,
    permission,
    settings
  },
  getters
})

export default store 