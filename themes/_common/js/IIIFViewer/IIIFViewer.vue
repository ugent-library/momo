<template>
  <div id="iiif-openseadragon-viewer" class="mt-2 mb-2">
  </div>
</template>

<script>

import OpenSeadragon from 'openseadragon'
import { mapState } from 'vuex'

export default {
  props: [
    'manifestUrl'
  ],
  computed: {
    ...mapState([
      'manifest',
    ])
  },
  watch: {
    manifest: function (manifest) {
        const layers = []
        manifest.sequences[0].canvases.forEach(val => {
          layers.push(val.images[0].resource.service['@id'] + '/info.json')
        })
        OpenSeadragon({
            id:                 "iiif-openseadragon-viewer",
            prefixUrl:          "https://openseadragon.github.io/openseadragon/images/",
            preserveViewport:   true,
            sequenceMode:  true,
            showReferenceStrip: true,
            tileSources: layers
        });
    }
  },
  mounted () {
    this.$store.dispatch('loadManifest', this.manifestUrl)
  }
}
</script>

<style>
#iiif-openseadragon-viewer {
    height: 800px;
}
</style>
