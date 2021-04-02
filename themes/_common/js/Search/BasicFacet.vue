<template>
  <div class="facet">
    <h2 class="title" v-html="label"></h2>

    <b-form-group>
        <b-form-checkbox-group
            v-bind:id="`facet-checkbox-group-${facetId}`"
            v-bind:name="facetId"
            v-model="selected"
            stacked
        >
            <b-form-checkbox
            v-for="bucket in buckets"
            v-bind:id="bucket.key"
            :key="bucket.key"
            class="facet-checkbox-bucket"
            v-bind:value="bucket.key"
            >{{ bucket.key }} ({{ bucket.doc_count }})</b-form-checkbox>
        </b-form-checkbox-group>
    </b-form-group>
  </div>
</template>

<script>
export default {
  props: [
    'label', 'facetId', 'buckets'
  ],
  data () {
    return {
      // Keep local state separate from store state, the latter is used for Elastic.
      // The former is used for tracking the active state of the checkbox.
      selected: []
    }
  },
  watch: {
    selected: function (newSelected, oldSelected) {
      if (newSelected.length !== oldSelected.length) {
        this.$store.dispatch('toggleFacet', { facetId: this.facetId, buckets: this.selected })
      }
    }
  },
  mounted: function () {
    const p = new URL(window.location).searchParams
    const pKey = `f[${this.facetId}]`

    if (p.has(pKey)) {
      this.selected = p.getAll(pKey)
      this.$store.dispatch('toggleFacet', { facetId: this.facetId, buckets: this.selected })
    }
  }
}
</script>
