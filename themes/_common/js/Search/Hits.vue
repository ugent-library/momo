<template>
    <ul id="search-results" class="list-group-flush mb-4">
        <li
        v-for="hit in hits"
        v-bind:id="hit.id"
        :key="hit.id"
        class="result search-result-item list-group-item"
        >
            <ul class="list-inline">
                <li>
                    <span class="badge badge-primary" v-html="hit.type"></span>
                </li>
            </ul>
            <a class="title" :href="hitUrl(hit)" v-html="hit.title">
            </a>
            <ul class="list-unstyled">
              <li v-for="c in hit.contributors">
                <span v-html="c"></span>
              </li>
            </ul>
        </li>
    </ul>
</template>

<script>
export default {
  props: [
    'hits'
  ],
  methods: {
    hitUrl: function (hit) {
      return window.location.pathname + '/' + hit.id
    }
  }
}
</script>
