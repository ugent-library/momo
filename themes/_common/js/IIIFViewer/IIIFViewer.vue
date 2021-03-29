<template>
  <div class="iiif-example">
    <l-map :zoom="zoom" :center="location">
      <l-iiif :url="infoUrl" :options="opts" />
    </l-map>
  </div>
</template>

<script>
import * as L from 'leaflet'
import { LMap } from 'vue2-leaflet'
import LIiif from './LeafletIIIF'

export default {
  components: {
    LMap,
    LIiif
  },
  props: [
    'initialZoom',
    'initialOptions',
    'initialInfoUrl'
  ],
  data () {
    return {
      location: L.latLng(0, 0),
      zoom: this.initialZoom,
      opts: this.initialOptions,
      infoUrl: this.initialInfoUrl
    }
  }
}
</script>

<style>
.iiif-example {
    height: 600px;
    width: 800px;
}
</style>
