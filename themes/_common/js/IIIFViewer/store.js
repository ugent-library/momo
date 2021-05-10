import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

const state = {
  manifest: null
}

const mutations = {
  setManifest (state, { manifest }) {
    state.manifest = manifest
  }
}

const actions = {
  loadManifest ({ commit }, manifestUrl) {
    axios
      .get(manifestUrl)
      .then(response => {
        commit('setManifest', { manifest: response.data })
      })
  }
}

const getters = {
  getManifest: (state) => { return state.manifest }
}

export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations
})
