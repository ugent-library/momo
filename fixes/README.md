These instructions assume Elasticsearch 6.x running on `localhost:9200`.

## Index recs

Dependencies:
* perl >= 5.10
* cpanm Catmandu
* cpanm Catmandu::Store::ElasticSearch
* cpanm Catmandu::MARC

```
cd fixes
cat recs.json | ./index_recs.sh
```

## Convert biblio data

```
cd fixes
cat biblio-export.json | ./biblio_to_momo.sh > momo_biblio.json
```

## Convert vlerick data

```
cd fixes
cat vlerick-export.json | ./vlerick_to_momo.sh > momo_biblio.json
```

## Convert orpheus data

```
cd fixes
cat orpheus.marcxml | ./orpheus_to_momo.sh > momo_orpheus.json
```
