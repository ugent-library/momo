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
      selected: []
    }
  },
  watch: {
    selected: function () {
      const p = new URL(window.location).searchParams
      const k = `f[${this.facetId}]`
      p.delete(k)

      for (const bucket of this.selected) {
        p.append(k, bucket)
      }

      const pStr = '?' + p.toString()

      window.history.pushState({}, '', pStr)

      this.$emit('bucket-selected', { facetId: this.facetId, buckets: this.selected })
    }
  }
}
</script>