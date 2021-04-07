import axios from 'axios'

export default {

  async query (query, page, size, selectedBuckets, commit) {
    const path = window.location.pathname
    const p = new URLSearchParams()

    p.append('q', query)
    p.append('skip', (page - 1) * size)
    p.append('size', size)

    Object
      .keys(selectedBuckets)
      .forEach((key) => {
        const pKey = `f[${key}]`
        p.delete(pKey)
        selectedBuckets[key].forEach(el => { p.append(pKey, el) })
      })

    const pStr = '?' + p.toString()

    const res = await axios.get(`${path}/search${pStr}`)

    if (res.status !== 200 && res.statusText !== 'OK') {
      throw new Error(`[HTTP: ${res.status}] ${res.statusText}`)
    }

    commit(
      res.data.hits,
      res.data.aggregation.facets,
      res.data.total
    )
  }
}
