import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import 'leaflet/dist/leaflet.css'
import IIIFViewer from './IIIFViewer'

export class IIIFViewerApp {
  static create (element, props) {
    Vue.use(BootstrapVue)

    const matches = document.body.querySelectorAll(element)
    for (const match of matches) {
      if (document.body.contains(match)) {
        const App = Vue.extend(IIIFViewer)
        const vm = new App({
          el: match,
          props: props,
          propsData: { ...match.dataset }
        })
        vm.$mount(match)
      }
    }
  }
}
