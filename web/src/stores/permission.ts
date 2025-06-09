import { RouteRecordRaw } from 'vue-router'
import { constantRoutes } from '@/router'

interface PermissionState {
  routes: RouteRecordRaw[]
  addRoutes: RouteRecordRaw[]
}

const state: PermissionState = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state: PermissionState, routes: RouteRecordRaw[]) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes({ commit }, roles) {
    return new Promise(resolve => {
      let accessedRoutes
      if (roles.includes('admin')) {
        accessedRoutes = constantRoutes || []
      } else {
        accessedRoutes = filterAsyncRoutes(constantRoutes, roles)
      }
      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

function filterAsyncRoutes(routes: RouteRecordRaw[], roles: string[]) {
  const res: RouteRecordRaw[] = []

  routes.forEach(route => {
    const tmp = { ...route }
    if (hasPermission(roles, tmp)) {
      if (tmp.children) {
        tmp.children = filterAsyncRoutes(tmp.children, roles)
      }
      res.push(tmp)
    }
  })

  return res
}

function hasPermission(roles: string[], route: RouteRecordRaw) {
  if (route.meta && route.meta.roles) {
    return roles.some(role => route.meta.roles.includes(role))
  } else {
    return true
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
} 