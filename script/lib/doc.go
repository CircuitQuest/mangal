package lib

import lua "github.com/yuin/gopher-lua"

// Mangal script Lua doc.
func LuaDoc() string {
	return Lib(lua.NewState(), nil).LuaDoc()
}
