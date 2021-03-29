<template>
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
</template>

<script>
export default {
  data () {
    return {
      query: ""
    }
  },
  watch: {
    query: function () {
      const p = new URL(window.location).searchParams
      p.set('q', this.query)
      const pStr = '?' + p.toString()

      window.history.pushState({}, '', pStr)

      this.$emit('query-submitted', this.query)
    }
  },
  mounted: function () {
    const p = new URL(window.location).searchParams

    if (p.has('q')) {
      this.query = p.get('q')
    }
  }
}
</script>