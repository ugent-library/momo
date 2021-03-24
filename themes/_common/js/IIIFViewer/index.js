import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import 'leaflet/dist/leaflet.css'
import IIIFViewer from './IIIFViewer'

export class IIIFViewerApp {
  constructor (element) {
    const apps = [];

    (function () {
      Vue.use(BootstrapVue)

      const matches = document.body.querySelectorAll(element)
      for (const match of matches) {
        if (document.body.contains(match)) {
          apps.push(new Vue({
            el: match,
            render: h => h(IIIFViewer)
          }))
        }
      }
    })(apps)
  }
}
