<template>
  <div id="iiif-leaflet-viewer">
    <l-map :zoom="zoom" :center="location">
      <l-iiif :url="manifestUrl" :options="opts" />
    </l-map>
  </div>
</template>

<script>
import L from 'leaflet'
import { LMap } from 'vue2-leaflet'
import 'leaflet-iiif'
import LIiif from './LeafletIIIF'

export default {
  components: {
    LMap,
    LIiif
  },
  props: [
    'initialZoom',
    'initialOptions',
    'initialInfoUrl',
    'manifestUrl'
  ],
  data () {
    return {
      location: L.latLng(0, 0),
      zoom: this.initialZoom,
      opts: this.initialOptions
    }
  }
}
</script>

<style>
#iiif-leaflet-viewer {
    height: 600px;
    width: 800px;
}
</style>
