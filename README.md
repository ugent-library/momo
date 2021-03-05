:warning: EXPERIMENTAL AND UNDER HEAVY DEVELOPMENT :warning:

## Requirements

* go
* postgresql
* elasticsearch 6.x
* npm

## Search index

To create the index:

```
go run cmd/momo/main.go rec create-index
```

To delete the index:

```
go run cmd/momo/main.go rec delete-index
```

## Import records

```
go run cmd/momo/main.go rec add myrecs1.json myrecs2.json
```

See `fixes/README.md` on how to convert data.

## Configuration

Configuration can be passed as an argument:

```
go run cmd/momo/main.go server --port 4000
```

Or as an env variable:

```
MOMO_PORT=4000 go run cmd/momo/main.go server
```

## Themes

Install node dependencies:

```bash
npm install
```

Momo contains a default theme called Opale. This theme will be compiled and installed during installation.

Compile a theme manually. Replace THEME with the name of your theme:

```
npx mix --mix-config themes/THEME/webpack.mix.js
```

Watching:

```
npx mix --mix-config themes/opale/webpack.mix.js watch # live reload in development
npx mix --production # production
```

Laravel Mix [documentation](https://laravel.com/docs/8.x).

## Start server

```
go run cmd/momo/main.go server
```

To run the server with live reload:

```
go get -u github.com/cosmtrek/air
air -c .air.toml
```
