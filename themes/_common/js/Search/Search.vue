// based on https://github.com/elastic/go-elasticsearch/tree/master/_examples/xkcdsearch

<template>
  <div class="container momo-search">
    <div class="row">
        <div class="col-12">
          <basic-query
            v-on:query-submitted="loadResults"
          ></basic-query>
        </div>
    </div>

    <div class="row">
      <div class="col-12">
        <!-- <p>Start opnieuw</p>
        <p>Filters</p> -->
        <p>{{total }} gevonden</p>
      </div>
    </div>

    <div class="row">

      <div class="col-3">
          <basic-facet
            v-for="facet in facets"
            v-bind:key="facet.id"
            v-bind:facetId="facet.id"
            v-bind:label="facet.label"
            v-bind:buckets="facet.buckets"
            v-on:bucket-selected="loadResults"
          ></basic-facet>
      </div>

      <div class="col-9">
        <div class="container search-results">

          <!-- <div class="row">
            <p>Sorteren</p>
            <p>Per pagina</p>
          </div> -->

          <div class="row">
            <div class="col-12">
              <div v-if="total < 1 && !state.loading" class="no-results">
                <p>Sorry, no results for your query&hellip;</p>
              </div>

              <ul id="search-results" class="list-unstyled">
                <li
                v-for="hit in hits"
                v-bind:id="hit.id"
                :key="hit.id"
                class="result search-result-item"
                >
                  <ul class="list-inline mbottom-small">
                    <li>
                      <span class="btn btn-primary btn-tag" v-html="hit.type"></span>
                    </li>
                  </ul>
                  <a :href="hitUrl(hit)"
                    ><h2 class="title" v-html="hit.title"></h2>
                  </a>
                </li>
              </ul>

              <p v-show="state.fetching" id="loading-results">Loading results...</p>
              <p v-show="state.loading" id="loading-app">Loading the application...</p>
              <p v-if="state.error" id="app-error">
                [{{ state.error.status }}] {{ state.error.statusText }}
              </p>
            </div>
          </div>

          <div class="row">
            <div class="col-12">
              <search-pagination
                v-if="state.initialized"
                v-bind:total="total"
                v-on:page-selected="loadResults"
              ></search-pagination>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
import BasicFacet from './BasicFacet.vue'
import BasicQuery from './BasicQuery.vue'
import SearchPagination from './SearchPagination.vue'

export default {
  components: {
    BasicFacet,
    BasicQuery,
    SearchPagination
  },
  data () {
    return {
      state: {
        initialized: false,
        loading: true,
        fetching: false,
        error: null
      },
      hits: [],
      facets: {},
      total: 0
    }
  },
  methods: {
    loadResults: function () {
      const self = this

      if (self.state.fetching) {
        return
      }
      self.state.fetching = true

      const path = window.location.pathname
      const p = new URL(window.location).searchParams
      const pStr = '?' + p.toString()

      window.history.pushState({}, '', pStr)

      window
        .fetch(`${path}/search${pStr}`)
        .then(function (res) {
          if (!res.ok) {
            return Promise.reject(res)
          }
          return res.json()
        })
        .then(function (res) {
          const hits = []
          res.hits.forEach(function (h) {
            const hit = {
              id: h.id,
              type: h.type
            }

            if (h.highlight && h.highlight['metadata.title.ngram']) {
              hit.title = h.highlight['metadata.title.ngram'][0]
            } else {
              hit.title = h.metadata.title
            }
            hits.push(hit)
          })

          const facets = res.aggregation.facets

          Object
            .keys(facets)
            .filter( (key) => { return key !== 'doc_count' })
            .map( (key) => {
              self.facets[key] = {
                id: key,
                label: key.replace(/^(.)/, (_, c) => c.toUpperCase()), // uppertitle
                buckets: facets[key].facet.buckets
              }
            })

          self.total = res.total
          self.hits = hits
        })
        .then(function () {
          self.state.initialized = true
          self.state.loading = false
          self.state.fetching = false
        })
        .catch(function (response) {
          self.state.loading = false
          self.state.fetching = false
          self.state.error = response
        })
    },
    hitUrl: function (hit) {
      return window.location.pathname + '/' + hit.id
    }
  },

  mounted: function () {
    this.loadResults()
  }
}
</script>
