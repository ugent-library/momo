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
export default {
  props: [
    'total'
  ],
  data () {
    return {
      page: 1,
      size: 10
    }
  },
  watch: {
    page: function () {
      const p = new URL(window.location).searchParams
      p.set('skip', (this.page - 1) * this.size)
      p.set('size', this.size)
      const pStr = '?' + p.toString()

      window.history.pushState({}, '', pStr)

      this.$emit('page-selected', this.page)
    }
  },
  mounted: function () {
    const p = new URL(window.location).searchParams

    if (p.has('size')) {
      this.size = parseInt(p.get('size'), 10)
    }
    if (p.has('skip')) {
      this.page = Math.ceil((parseInt(p.get('skip'), 10) + 1) / this.size)
    }
  }
}
</script>