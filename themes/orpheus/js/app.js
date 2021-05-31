import { IIIFViewerApp } from '../../_common/js/IIIFViewer'
import { SearchApp } from '../../_common/js/Search'
import { TaggerApp } from '../../_common/js/Tagger'

window.addEventListener('DOMContentLoaded', () => {
  SearchApp.create('#search')
  IIIFViewerApp.create('#iiif-viewer', {})
  TaggerApp.create('#tagger')
})
