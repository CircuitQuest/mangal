# Mangal 5

> Work in progress...

Fork of [metafates/mangal@v5](/metafates/mangal/tree/v5) (originaly [metafates/mangal](/metafates/mangal)). Which is based on metafates' [mangalorg/libmangal](/mangalorg/libmangal), [mangalorg/luaprovider](/mangalorg/luaprovider) and [mangalorg/saturno](/mangalorg/saturno).

Basically all work originally from [metafates](/metafates), but he is kinda abandoning his work all the time so I'm making this fork as a personal thing in case I need to add/fix things.

## Development

For now, just notes for myself while I work on this.

### Requirements

- `go>=1.21`
- `npm`
  - [`oapi-codegen`](https://github.com/deepmap/oapi-codegen): `go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest`
  - [`enumer`](https://github.com/dmarkham/enumer): `go install github.com/dmarkham/enumer@latest`
- [`just`](https://github.com/casey/just)

### Install/use

Update modules:

```sh
just update
```

Install [openapi-typescript](https://www.npmjs.com/package/openapi-typescript):

```sh
cd web/ui
npm i -D openapi-typescript
npm audit fix
cd ../..
```

Generate assets:

```sh
just generate
```

Then it can be used directly while debugging:

```sh
go run main.go ...
```

Or by compiling it and then running it:

```sh
just build # or full or install
./mangal # or mangal if it was installed
```

**Note:** If the previous mangal was used, then remove the config (usually located at `~/.config/mangal`) as otherwise these commands will fail.
