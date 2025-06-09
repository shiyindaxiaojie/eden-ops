interface AppState {
  sidebar: {
    opened: boolean
    withoutAnimation: boolean
  }
  device: string
}

const state: AppState = {
  sidebar: {
    opened: localStorage.getItem('sidebarStatus') ? !!+localStorage.getItem('sidebarStatus') : true,
    withoutAnimation: false
  },
  device: 'desktop'
}

const mutations = {
  TOGGLE_SIDEBAR: (state: AppState) => {
    state.sidebar.opened = !state.sidebar.opened
    state.sidebar.withoutAnimation = false
    if (state.sidebar.opened) {
      localStorage.setItem('sidebarStatus', '1')
    } else {
      localStorage.setItem('sidebarStatus', '0')
    }
  },
  CLOSE_SIDEBAR: (state: AppState, withoutAnimation: boolean) => {
    localStorage.setItem('sidebarStatus', '0')
    state.sidebar.opened = false
    state.sidebar.withoutAnimation = withoutAnimation
  },
  TOGGLE_DEVICE: (state: AppState, device: string) => {
    state.device = device
  }
}

const actions = {
  toggleSideBar({ commit }) {
    commit('TOGGLE_SIDEBAR')
  },
  closeSideBar({ commit }, { withoutAnimation }) {
    commit('CLOSE_SIDEBAR', withoutAnimation)
  },
  toggleDevice({ commit }, device) {
    commit('TOGGLE_DEVICE', device)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
} 