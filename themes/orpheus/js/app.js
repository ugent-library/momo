import { IIIFViewerApp } from '../../_common/js/IIIFViewer'
import { SearchApp } from '../../_common/js/Search'

window.addEventListener('DOMContentLoaded', () => {
  const apps = [];
  (function () {
    apps.push(new SearchApp('#search'))
  })(apps)
})

window.IIIFViewer = {
  init: (element, config) => {
    const apps = [];
    (function () {
      apps.push(new IIIFViewerApp(element, config))
    })(apps)
  }
}
