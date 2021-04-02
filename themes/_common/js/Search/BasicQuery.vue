<template>
    <form v-on:submit.prevent="submit">
        <div class="form-group">
            <input
            v-model="query"
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
      query: this.$store.state.query
    }
  },
  methods: {
    submit () {
      this.$store.dispatch('submitQuery', this.query)
    }
  },
  mounted: function () {
    const p = new URL(window.location).searchParams

    if (p.has('q')) {
      this.$store.dispatch('submitQuery', p.get('q'))
      this.query = p.get('q')
    }
  }
}
</script>
