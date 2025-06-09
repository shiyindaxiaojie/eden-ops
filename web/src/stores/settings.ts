interface SettingsState {
  theme: string
  sidebarLogo: boolean
}

const state: SettingsState = {
  theme: localStorage.getItem('theme') || 'default',
  sidebarLogo: true
}

const mutations = {
  CHANGE_SETTING: (state: SettingsState, { key, value }) => {
    if (Object.prototype.hasOwnProperty.call(state, key)) {
      state[key] = value
      if (key === 'theme') {
        localStorage.setItem('theme', value)
      }
    }
  }
}

const actions = {
  changeSetting({ commit }, data) {
    commit('CHANGE_SETTING', data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
} 