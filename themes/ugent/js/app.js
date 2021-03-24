import { SearchApp } from '../../_common/js/Search'

window.addEventListener('DOMContentLoaded', () => {
  const apps = [];

  (function () {
    apps.push(new SearchApp('#search'))
  })(apps)
})
