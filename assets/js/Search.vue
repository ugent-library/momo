// based on https://github.com/elastic/go-elasticsearch/tree/master/_examples/xkcdsearch

<template>
  <div>
    <div id="search-form">
        <form v-on:submit.prevent="">
        <input v-model.lazy="query" class="form-control" type="text" size="50" placeholder="Search...">
        </form>
        <p class="total"><span class="label">total: </span><span class="content">{{total}}</span></p>
    </div>

    <div v-if="total < 1 && !state.loading" class="no-results">
        <p>Sorry, no results for your query&hellip;</p>
    </div>

    <ul v-for="result in results" v-bind:id="result.id" :key="result.id" id="search-results" class="result">
        <li>
        <span v-html="result.title" class="title"></span>
        </li>
    </ul>

    <p v-show="state.fetching" id="loading-results">Loading results...</p>
    <p v-show="state.loading" id="loading-app">Loading the application...</p>
    <p v-if="state.error" id="app-error">[{{state.error.status}}] {{state.error.statusText}}</p>

    <b-pagination
        v-model="page"
        :total-rows="total"
        :per-page="perPage"
        aria-controls="search-results"
    ></b-pagination>
  </div>
</template>

<script>
export default {
    data() {
        return {
            state: {
            error: null,
            loading: true,
            fetching: false,
            replaceResults: true
            },
    
            query: '',
            results: [],
            total: 0,
            page: 1,
            perPage: 10
        };
    },
    methods: {
        loadResults: function () {
          var self = this
  
          if (self.state.fetching) { return }
          self.state.fetching = true
  
            var queryParams = `?page=${encodeURIComponent(self.page)}&q=${encodeURIComponent(self.query)}`
            window.history.pushState({}, '', queryParams)
  
            var p = window.location.pathname
            window.fetch(`${p}/search${queryParams}`)
            .then(function (response) {
              if (!response.ok) { return Promise.reject(response) }
              return response.json()
            })
            .then(function (response) {
              var results = []
  
              response.hits.forEach(function (r) {
                var result = {
                    title: r.title
                }
  
                results.push(result)
              })
  
              self.total = response.total
  
              if (self.state.replaceResults) {
                self.results = results
              }
            })
            .then(function () {
              self.state.loading = false
              self.state.fetching = false
            })
            .catch(function (response) {
              self.state.loading = false
              self.state.fetching = false
              self.state.error = response
            })
        },
    },
  
    watch: {
        query: function () {
          this.state.replaceResults = true
          this.loadResults()
        },
        page: function () {
          this.state.replaceResults = true
          this.loadResults()
        },
    },
  
    created: function () {
        var self = this
  
        var q = document.location.search.split('q=')[1]
        var p = document.location.search.split('page=')[1]
        if (q) { self.query = decodeURIComponent(q) }
        if (p) { self.page = decodeURIComponent(p) }
        self.loadResults()
    }
};
</script>