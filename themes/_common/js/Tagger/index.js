import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import Tagger from './Tagger'
import store from './store'

export class TaggerApp {
  static create (element, props) {
    Vue.use(BootstrapVue)

    const matches = document.body.querySelectorAll(element)
    for (const match of matches) {
      if (document.body.contains(match)) {
        const App = Vue.extend(Tagger)
        new App({ // eslint-disable-line no-new
          store,
          el: match,
          props: props,
          propsData: { ...match.dataset }
        })
      }
    }
  }
}
