import { IIIFViewerApp } from '../../_common/js/IIIFViewer'
import { SearchApp } from '../../_common/js/Search'

window.addEventListener('DOMContentLoaded', () => {
  SearchApp.create('#search')
})

window.IIIFViewer = {
  init: (element, config) => {
    IIIFViewerApp.create(element, config)
  }
}
