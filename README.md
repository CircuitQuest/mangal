# Mangal 5

> Work in progress...

Fork of [metafates/mangal@v5](/metafates/mangal/tree/v5) (originaly [metafates/mangal](/metafates/mangal)). Which is based on metafates' [mangalorg/libmangal](/mangalorg/libmangal), [mangalorg/luaprovider](/mangalorg/luaprovider) and [mangalorg/saturno](/mangalorg/saturno).

Basically all work originally from [metafates](/metafates), but he is kinda abandoning his work all the time so I'm making this fork as a personal thing in case I need to add/fix things.

## Development

For now, just notes for myself while I work on this.

### Requirements

- `go>=1.21`
  - [`oapi-codegen`](https://github.com/deepmap/oapi-codegen): `go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest`
  - [`enumer`](https://github.com/dmarkham/enumer): `go install github.com/dmarkham/enumer@latest`
- [`just`](https://github.com/casey/just)
- `npm`

### Setup

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

### Usage

Can be used without compiling for testing purposes:

```sh
go run main.go ...
```

Or compiled/installed:

```sh
just build # or full or install
./mangal # or mangal if it was installed
```

**Note:** If the previous mangal was used, then remove the config (usually located at `~/.config/mangal`) as otherwise these commands will fail.

#### Providers

A new provider can be created by running:

```sh
mangal providers new
```

And a new test provider will be placed in `$XDG_CONFIG/mangal/providers/test`. A sample provider (mangapill) can be found in [mangalorg/saturno](/mangalorg/saturno).

#### Script mode

There is no documentation for any of the code, but I figured out how to use the script mode, it is a bit similar to [mangalorg/mangalcli](/mangalorg/mangalcli): a `run.lua` and a provider is required.

I'm using the sample `mangapill.lua` mentioned above with the following modified `run.lua`:

```lua
local mangal = require('mangal')
local client = mangal.client
local fzf = mangal.prompt.fzf
local json = mangal.json

local mangas = client.search_mangas(Vars.title) -- search with the given title
local volumes = client.manga_volumes(fzf.find_one(mangas, function(manga)
  return manga:info().title
end)) -- select the first manga

local chapters = {}

-- get all chapters of the manga
for _, volume in ipairs(volumes) do
  for _, chapter in ipairs(client.volume_chapters(volume)) do
    table.insert(chapters, chapter)
  end
end

-- chapters encoded in json format for later use, e.g. pipe to jq
json.print(chapters)
```
