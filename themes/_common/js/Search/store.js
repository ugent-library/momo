import Vue from 'vue'
import Vuex from 'vuex'
import Elastic from './es'

Vue.use(Vuex)

const state = {
  query: '',
  total: 0,
  page: 1,
  size: 10,
  hits: [],
  facets: {},
  selectedBuckets: {},
  initialized: false,
  loading: true,
  error: false,
  fetching: false
}

const mutations = {
  setLoading (state, { loading }) {
    state.loading = loading
  },
  setError (state, { error }) {
    state.error = error
  },
  setFetching (state, { fetching }) {
    state.fetching = fetching
  },
  setInitialized (state, { initialized }) {
    state.initialized = initialized
  },
  submitQuery (state, { query }) {
    state.query = query
  },
  changePage (state, { page }) {
    state.page = page
  },
  changeSize (state, { size }) {
    state.size = size
  },
  setTotal (state, { total }) {
    state.total = total
  },
  setHits (state, { hits }) {
    const result = []
    hits.forEach(function (h) {
      const hit = {
        id: h.id,
        type: h.type
      }
      if (h.highlight && h.highlight['metadata.title.ngram']) {
        hit.title = h.highlight['metadata.title.ngram'][0]
      } else {
        hit.title = h.metadata.title
      }
      if (h.highlight && h.highlight['metadata.author.name.ngram']) {
        hit.contributors = h.highlight['metadata.author.name.ngram']
      }
      result.push(hit)
    })

    state.hits = result
  },
  setFacets (state, { facets }) {
    const result = {}

    Object
      .keys(facets)
      .filter((key) => { return key !== 'doc_count' })
      .forEach((key) => {
        result[key] = {
          id: key,
          label: key.replace(/^(.)/, (_, c) => c.toUpperCase()), // uppertitle
          buckets: facets[key].facet.buckets
        }
      })

    state.facets = result
  },
  setSelectedBuckets (state, { facet }) {
    state.selectedBuckets[facet.facetId] = facet.buckets
  }
}

const actions = {
  setTotal ({ commit }, total) {
    commit('setTotal', { total })
  },
  submitQuery ({ dispatch, commit }, query) {
    commit('submitQuery', { query })

    const p = new URL(window.location).searchParams
    p.set('q', this.state.query)
    const pStr = '?' + p.toString()

    window.history.pushState({}, '', pStr)

    dispatch('loadResults')
  },
  toggleFacet ({ dispatch, commit }, facet) {
    commit('setSelectedBuckets', { facet })

    const p = new URL(window.location).searchParams
    const pKey = `f[${facet.facetId}]`
    p.delete(pKey)

    for (const bucket of facet.buckets) {
      p.append(pKey, bucket)
    }

    const pStr = '?' + p.toString()

    window.history.pushState({}, '', pStr)

    dispatch('loadResults')
  },
  changePage ({ dispatch, commit }, page) {
    commit('changePage', { page })
    dispatch('setURLParams')
    dispatch('loadResults')
  },
  changeSize ({ dispatch, commit }, size) {
    commit('changeSize', { size })
    dispatch('setURLParams')
    dispatch('loadResults')
  },
  setURLParams () {
    const p = new URL(window.location).searchParams
    p.set('page', this.state.page)
    p.set('size', this.state.size)
    const pStr = '?' + p.toString()

    window.history.pushState({}, '', pStr)
  },
  loadResults ({ commit }) {
    commit('setFetching', { fetching: true })

    Elastic
      .query(
        this.state.query,
        this.state.page,
        this.state.size,
        this.state.selectedBuckets,
        (hits, facets, total) => {
          commit('setHits', { hits })
          commit('setFacets', { facets })
          commit('setTotal', { total })
        }
      )
      .then(() => {
        // Response was OK, switching status flags
        commit('setLoading', { loading: false })
        commit('setFetching', { fetching: false })
        commit('setInitialized', { initialized: true })
      })
      .catch(e => {
        // Error occurred
        commit('setLoading', { loading: false })
        commit('setFetching', { fetching: false })
        commit('setError', { error: e.message })
      })
  }
}

const getters = {
  getPage: (state) => { return state.page },
  getSize: (state) => { return state.size }
}

export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations
})
