# Mangal

> Work in progress...

Fork of [metafates/mangal@v5](https://github.com/metafates/mangal/tree/v5).

Version 5 is now based on metafates' [mangalorg](https://github.com/mangalorg) suite, which I've forked to work on them myself:

- [libmangal](https://github.com/luevano/libmangal)
- [luaprovider](https://github.com/luevano/luaprovider)
- [gopher-luadoc](https://github.com/luevano/gopher-luadoc)
- [saturno](https://github.com/luevano/saturno)

## Development

For now, just notes for myself while I work on this.

### Requirements

The following needs to be installed before starting development:

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

Generate assets and code:

```sh
just generate
```

Then it can be used without compiling for testing purposes:

```sh
go run main.go ...
```

Or compiled/installed:

```sh
just build # or full or install
./mangal
```

## Usage

By default if just executing the `mangal` command, it will start the TUI. To list all options:

```sh
mangal -h
```

**Note:** If the previous mangal (v4) was used, then remove the config (usually located at `~/.config/mangal`) as otherwise these commands will fail.

### Config

Mangal contains sensitive defaults that can be edited with some command flags or by editing a config file. By default no config is written to disk, to write the default config to disk run:

```sh
mangal config write
```

It will be written to `$XDG_CONFIG_HOME/mangal/mangal.toml`.

Then, the config can be set either by modifying the file itself or by running:

```sh
mangal config set <key> <value>
```

For example:

```sh
mangal config set download.path $HOME/downloads_test
```

### Providers

These are Lua scripts that handle the site scrape logic (search mangas, list mangas/chapters/images, etc). As shown in [mangalorg/saturno](/mangalorg/saturno), each provider is a directory with the `main.lua` (scraper) and `mangal.toml` (metadata) files inside.

Providers need to be placed in `$XDG_CONFIG_HOME/mangal/providers`.

A new provider can be created by running:

```sh
mangal providers new
```

It is automatically placed in the required path.

### Templates

Special functions are available and can be shown by running:

```sh
mangal templates funcs
```

Run a template with custom functions available:

```sh
mangal templates exec <template> -v <values>
```

Where `<template>` is in the form of a go [text/template](https://pkg.go.dev/text/template) and `<values>` a JSON. For example:

```sh
mangal templates exec "{{sanitizeWindows .Test}}" -v "{\"Test\":\"asdfg<>:?\"}"
```

Which will output `asdfg`.

### Modes

#### Script

Similar to [mangalorg/mangalcli](https://github.com/mangalorg/mangalcli) where a `run.lua` and a "provider" is required:

```sh
mangal script -f run.lua -p mangapill -v "title=tengoku"
```

I'm using the sample `mangapill.lua` mentioned in [Providers](#providers) with the following modified `run.lua`:

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
