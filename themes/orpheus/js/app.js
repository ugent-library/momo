import { SearchApp } from '../../_common/js/Search'
import { IIIFViewerApp } from '../../_common/js/IIIFViewer'

window.addEventListener('DOMContentLoaded', () => {
  SearchApp.create('#search')
  IIIFViewerApp.create('#iiif-viewer', {})
})
