<template>
    <b-pagination
      v-model="page"
      :total-rows="total"
      :per-page="size"
      aria-controls="search-results"
      class="list-unstyled"
    ></b-pagination>
</template>

<script>
import { mapState } from 'vuex'

export default {
  data () {
    return {
      page: this.$store.getters.getPage,
      size: this.$store.getters.getSize
    }
  },
  computed: {
    ...mapState([
      'total',
    ])
  },
  watch: {
    page: function (newPage, oldPage) {
      if (newPage !== oldPage) {
        this.$store.dispatch('changePage', newPage)
      }
    }
  },
  mounted: function () {
    const p = new URL(window.location).searchParams

    if (p.has('size')) {
      this.$store.dispatch('changeSize', p.get('size'))
    }

    if (p.has('page')) {
      this.$store.dispatch('changePage', p.get('page'))
    }
  }
}
</script>
