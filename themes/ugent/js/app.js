import { SearchApp } from '../../_common/js/Search'
import { IIIFViewerApp } from '../../_common/js/IIIFViewer'
import { TaggerApp } from '../../_common/js/Tagger'

window.addEventListener('DOMContentLoaded', () => {
  SearchApp.create('#search')
  IIIFViewerApp.create('#iiif-viewer', {})
  TaggerApp.create('#tagger')
})
