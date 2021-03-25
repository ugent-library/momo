import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import Search from './Search'

export class SearchApp {
  static create (element) {
    Vue.use(BootstrapVue)

    const matches = document.body.querySelectorAll(element)
    for (const match of matches) {
      if (document.body.contains(match)) {
        const vm = new Vue({
          render: h => h(Search)
        })
        vm.$mount(match)
      }
    }
  }
}
