import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import '../css/app.scss'

Vue.use(BootstrapVue)

import Search from "../js/Search";

if (document.getElementById("search")) {
  new Vue({
    el: "#search",
    render: h => h(Search)
  })
}
