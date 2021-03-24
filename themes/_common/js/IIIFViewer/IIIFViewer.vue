<template>
    <div class="iiif-example">
        <l-map :zoom="initialZoom" :center="initialLocation">
            <l-iiif :url="infoUrl" :options="opts" />
        </l-map>
    </div>
</template>

<script>
import * as L from 'leaflet'
import { LMap } from 'vue2-leaflet'
import LIiif from './LeafletIIIF'
const iiifOptions = {
  tileFormat: 'jpg',
  tileSize: 512
}
export default {
  components: {
    LMap,
    LIiif
  },
  data () {
    return {
      initialLocation: L.latLng(0, 0),
      initialZoom: 2,
      opts: iiifOptions,
      infoUrl: 'http://iiif.hualab.nl:8080/iiif/2/saftleven-panorama-1684.jpg/info.json'
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
