<template>
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
            <a :href="hitUrl(hit)">
                <h3 class="title" v-html="hit.title"></h3>
            </a>
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
