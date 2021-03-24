import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import Search from './Search'

export class SearchApp {
  constructor (element) {
    const apps = [];

    (function () {
      Vue.use(BootstrapVue)

      const matches = document.body.querySelectorAll(element)
      for (const match of matches) {
        if (document.body.contains(match)) {
          apps.push(new Vue({
            el: match,
            render: h => h(Search)
          }))
        }
      }
    })(apps)
  }
}
