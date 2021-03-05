## Dependencies

* perl >= 5.10
* cpanm Catmandu
* cpanm Catmandu::MARC

## Convert biblio data

```
cd fixes
cat biblio-export.json | ./biblio_to_momo.sh > momo_biblio.json
```

## Convert cageweb data

```
cd fixes
cat cageweb-export.json | ./cageweb_to_momo.sh > momo_cageweb.json
```

## Convert orpheus data

```
cd fixes
cat orpheus.marcxml | ./orpheus_to_momo.sh > momo_orpheus.json
```
