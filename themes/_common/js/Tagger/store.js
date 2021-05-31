import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
  form: {
    tags: []
  }
}

const mutations = {
  addToTags (state, { tag }) {
    const found = state.form.tags.find(el => { return tag === el })
    if (found === undefined) {
      state.form.tags.push(tag)
    }
  },
  removeFromTags (state, { tag }) {
    state.form.tags = state.form.tags.filter(el => {
      return el !== tag
    })
  }
}

const actions = {
  addTag ({ commit }, tag) {
    if (tag !== '') {
      commit('addToTags', tag)
    }
  },
  removeTag ({ commit }, tag) {
    commit('removeFromTags', tag)
  }
}

const getters = {
  getAllTags: (state) => { return state.form.tags }
}

export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations
})
