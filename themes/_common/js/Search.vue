// based on https://github.com/elastic/go-elasticsearch/tree/master/_examples/xkcdsearch

<template>
  <div class="container momo-search">
    <div class="row">
        <div class="col-12">
          <form v-on:submit.prevent="">
            <div class="form-group">
              <input
                v-model.lazy="query"
                class="form-control"
                type="text"
                placeholder="Search..."
              />
            </div>
          </form>
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
        <div
          v-for="facet in facets"
          v-bind:id="`facet-${facet.id}`"
          :key="facet.id"
          class="facet"
        >
          <h2 class="title" v-html="facet.label"></h2>

          <b-form-group>
            <b-form-checkbox-group
              v-bind:id="`facet-checkbox-group-${facet.id}`"
              v-model="state[facet.id]"
              v-bind:name="facet.id"
              stacked
            >
              <b-form-checkbox
                v-for="bucket in facet.buckets"
                v-bind:id="bucket.key"
                :key="bucket.key"
                class="facet-checkbox-bucket"
                v-bind:value="bucket.key"
              >{{ bucket.key }} ({{ bucket.doc_count }})</b-form-checkbox>
            </b-form-checkbox-group>
          </b-form-group>

        </div>

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
              <b-pagination
                v-if="state.initialized"
                v-model="page"
                :total-rows="total"
                :per-page="size"
                aria-controls="search-results"
                class="list-unstyled"
              ></b-pagination>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      state: {
        initialized: false,
        loading: true,
        fetching: false,
        error: null,
        type: [],
        collection: []
      },

      query: "",
      type: [],
      collection: [],
      hits: [],
      facets: [],
      total: 0,
      page: 1,
      size: 10,
    };
  },
  methods: {
    loadResults: function () {
      var self = this;

      if (self.state.fetching) {
        return;
      }
      self.state.fetching = true;

      var path = window.location.pathname;
      var p = new URL(window.location).searchParams;
      p.set("q", self.query);
      p.set("skip", (self.page - 1) * self.size);
      p.set("size", self.size);
      // TODO remove existing
      for (let f of ["collection", "type"]) {
        var k = `f[${f}]`
        p.delete(k)
        if (!self[f].length) {
          continue
        }
        for (let term of self[f]) {
          p.append(k, term)
        }
      }

      var pStr = "?" + p.toString();

      window.history.pushState({}, "", pStr);

      window
        .fetch(`${path}/search${pStr}`)
        .then(function (res) {
          if (!res.ok) {
            return Promise.reject(res);
          }
          return res.json();
        })
        .then(function (res) {
          var hits = [];
          res.hits.forEach(function (h) {
            var hit = {
              id: h.id,
              type: h.type
            };

            if (h.highlight && h.highlight["metadata.title.ngram"]) {
              hit.title = h.highlight["metadata.title.ngram"][0];
            } else {
              hit.title = h.metadata.title;
            }
            hits.push(hit);
          });

          var facets = ["collection", "type"].map(function(key) {
              return {
                "id" : key,
                "label": key.replace (/^(.)/, (_, c) => c.toUpperCase()), // uppertitle
                "buckets": res.aggregation.facets[key].facet.buckets
              };
          });

          self.facets = facets;
          self.total = res.total;
          self.hits = hits;
        })
        .then(function () {
          self.state.initialized = true;
          self.state.loading = false;
          self.state.fetching = false;
        })
        .catch(function (response) {
          self.state.loading = false;
          self.state.fetching = false;
          self.state.error = response;
        });
    },
    hitUrl: function (hit) {
      return window.location.pathname + "/" + hit.id;
    },
    toggleBucket: function(facet, bucket) {
      var self = this;
      self[facet]=self.state[facet];
    },
  },
  watch: {
    'state.type': function(buckets) {
      this.toggleBucket('type', buckets);
    },
    'state.collection': function(buckets) {
      this.toggleBucket('collection', buckets);
    },
    query: function () {
      this.loadResults();
    },
    page: function () {
      this.loadResults();
    },
    type: function() {
      this.loadResults();
    },
    collection: function() {
      this.loadResults();
    }
  },

  mounted: function () {
    var self = this;

    var p = new URL(window.location).searchParams;

    if (p.has("q")) {
      self.query = p.get("q");
    }
    for (let f of ["collection", "type"]) {
      var terms = p.getAll(`f[${f}]`)
      if (terms.size) {
        self[f] = terms
        self.state[f] = terms
      }
    }
    if (p.has("size")) {
      self.size = parseInt(p.get("size"), 10);
    }
    if (p.has("skip")) {
      self.page = Math.ceil((parseInt(p.get("skip"), 10) + 1) / self.size);
    }

    self.loadResults();
  },
};
</script>