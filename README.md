:warning: EXPERIMENTAL AND UNDER HEAVY DEVELOPMENT :warning:

## Requirements

* go
* postgresql
* elasticsearch 6.x
* npm

## Search index

To create the index:

```
go run cmd/momo/main.go rec index create
```

To delete the index:

```
go run cmd/momo/main.go rec index delete
```

## Import records

```
go run cmd/momo/main.go rec add myrecs1.json myrecs2.json
```

## Destructive reset

```
go run cmd/momo/main.go reset --confirm
```

## Configuration

Configuration can be passed as an argument:

```
go run cmd/momo/main.go server start --port 4000
```

Or as an env variable:

```
MOMO_PORT=4000 go run cmd/momo/main.go server start
```

## Themes

Install node dependencies:

```bash
npm install
```

Momo contains a default theme called Ugent. This theme will be compiled and installed during installation.

Compile a theme manually. Replace THEME with the name of your theme:

```
npx mix --mix-config themes/THEME/webpack.mix.js
```

Watching:

```
npx mix --mix-config themes/ugent/webpack.mix.js watch # live reload in development
npx mix --production # production
```

Building:

```
npx mix --production # production
```

Laravel Mix [documentation](https://laravel.com/docs/8.x).

## Start server

To start the development server with live reload:

```
npx run dev
```

Run the server directly:

```
go run cmd/momo/main.go server start
```

## Build

```
go build -o momo cmd/momo/main.go
```
