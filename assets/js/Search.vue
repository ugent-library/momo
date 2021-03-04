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

    <ul v-for="hit in hits" v-bind:id="hit.id" :key="hit.id" id="search-results" class="result">
        <li>
          <a :href="'/v/all/'+hit.id"><span v-html="hit.title" class="title"></span></a>
        </li>
    </ul>

    <p v-show="state.fetching" id="loading-results">Loading results...</p>
    <p v-show="state.loading" id="loading-app">Loading the application...</p>
    <p v-if="state.error" id="app-error">[{{state.error.status}}] {{state.error.statusText}}</p>

    <b-pagination
        v-model="page"
        :total-rows="total"
        :per-page="size"
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
                fetching: false
            },
    
            query: '',
            hits: [],
            total: 0,
            page: 1,
            size: 10
        };
    },
    methods: {
        loadResults: function () {
          var self = this
  
          if (self.state.fetching) { return }
            self.state.fetching = true
  
            var path = window.location.pathname
            var p = (new URL(window.location)).searchParams;
            p.set('q', self.query)
            p.set('skip', (self.page-1)*self.size)
            p.set('size', self.size)
            var pStr = '?'+p.toString()

            window.history.pushState({}, '', pStr)

            window.fetch(`${path}/search${pStr}`)
            .then(function (res) {
              if (!res.ok) { return Promise.reject(res) }
              return res.json()
            })
            .then(function (res) {
              self.total = res.total
              self.hits = res.hits
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
          this.loadResults()
        },
        page: function () {
          this.loadResults()
        },
    },
  
    created: function () {
        var self = this

        var p = (new URL(window.location)).searchParams;
        if (p.has('q')) {
            self.query = p.get('q')
        }
        if (p.has('size')) { 
            self.size = parseInt(p.get('size'), 10) 
        }
        if (p.has('skip')) {
            self.page = Math.ceil((parseInt(p.get('skip'), 10)+1)/self.size) 
        }

        self.loadResults()
    }
};
</script>