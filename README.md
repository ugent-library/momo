:warning: EXPERIMENTAL AND UNDER HEAVY DEVELOPMENT :warning:

## Requirements

* go
* postgresql
* elasticsearch 6.8
* npm

## Search index

To create the index:

```
go run cmd/momo-app/main.go index create
```

To delete the index:

```
go run cmd/momo-app/main.go index delete
```

See `fixes/README.md` on how to index data.

## Compile assets

```
cd assets
npm install
npx mix watch # live reload in development
npx mix --production # production 
```

Laravel Mix [documentation](https://laravel.com/docs/8.x).

## Start server

```
go run cmd/momo-app/main.go server
```

To run the server with live reload:

```
go get -u github.com/cosmtrek/air
air -c .air.toml
```
