// based on https://github.com/elastic/go-elasticsearch/tree/master/_examples/xkcdsearch

<template>
  <base-search>
    <div class="row">
        <div class="col-12">
          <basic-query></basic-query>
        </div>
    </div>

    <div class="row">
      <div class="col-12">
        <!-- <p>Start opnieuw</p>
        <p>Filters</p> -->
        <p>{{ total }} records found</p>
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
          ></basic-facet>

          <!-- <basic-facet
            v-if="facets['type']"
            v-bind:key="facets['type']"
            v-bind:facetId="facets['type'].id"
            v-bind:label="facets['type'].label"
            v-bind:buckets="facets['type'].buckets"
          ></basic-facet> -->
      </div>

      <div class="col-9">
        <div class="container search-results">

          <!-- <div class="row">
            <p>Sorteren</p>
            <p>Per pagina</p>
          </div> -->

          <div class="row">
            <div class="col-12">
              <div v-if="total < 1 && !loading" class="no-results">
                <p>Sorry, no results for your query&hellip;</p>
              </div>

              <hits v-bind:hits="hits"></hits>

              <p v-show="fetching" id="loading-results">Loading results...</p>
              <p v-show="loading" id="loading-app">Loading the application...</p>
              <p v-if="error" id="app-error">
                {{ error }}
              </p>
            </div>
          </div>

          <div class="row">
            <div class="col-12">
              <search-pagination
                v-if="initialized"
              ></search-pagination>
            </div>
          </div>
        </div>
      </div>

    </div>
  </base-search>
</template>

<script>
import BaseSearch from './BaseSearch.vue'
import BasicFacet from './BasicFacet.vue'
import BasicQuery from './BasicQuery.vue'
import SearchPagination from './SearchPagination.vue'
import Hits from './Hits.vue'

import { mapState } from 'vuex'

export default {
  components: {
    BaseSearch,
    BasicFacet,
    BasicQuery,
    Hits,
    SearchPagination
  },
  computed: {
    ...mapState([
      'hits',
      'facets',
      'total',
      'initialized',
      'loading',
      'fetching',
      'error'
    ])
  }
}
</script>
