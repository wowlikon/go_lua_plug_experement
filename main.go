package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	lua "github.com/cathalgarvey/gopher-lua"
)

func main() {
	pluginsDir := "./plugins"
	files, err := os.ReadDir(pluginsDir)
	if err != nil {
		fmt.Println("Error reading plugins directory:", err)
		return
	}

	for _, dir := range files {
		if !dir.IsDir() {
			continue
		}

		pluginDir := filepath.Join(pluginsDir, dir.Name())

		state := lua.NewState()
		state.Register("get", Get)
		state.Register("exit", Exit)
		state.Register("square", Square)
		state.Register("json_decode", JsonDecode)
		state.Register("json_encode", JsonEncode)

		files, err := os.ReadDir(pluginDir)
		if err != nil {
			fmt.Println("Error reading plugin directory:", err)
			continue
		}

		for _, file := range files {
			if file.IsDir() || file.Name() == "main.lua" {
				continue
			}

			loadLuaFile(state, filepath.Join(pluginDir, file.Name()))
		}

		if err := state.DoFile(filepath.Join(pluginDir, "main.lua")); err != nil {
			fmt.Printf("Error running %s: %s\n", filepath.Join(pluginDir, "main.lua"), err)
		}
	}
}

func loadLuaFile(state *lua.LState, path string) {
	if err := state.DoFile(path); err != nil {
		fmt.Println("Error loading library:", err)
	}
}

func Square(L *lua.LState) int {
	lv := L.ToInt(1)
	L.Push(lua.LNumber(lv * lv))
	return 1
}

func Exit(L *lua.LState) int {
	lv := L.ToInt(1)
	os.Exit(lv)
	return 1
}

func Get(L *lua.LState) int {
	url := L.ToString(1)
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0
	}
	L.Push(lua.LString(string(body)))
	return 1
}

func JsonEncode(L *lua.LState) int {
	table := L.ToTable(1)
	json, err := json.Marshal(table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(json)))
	return 1
}

func JsonDecode(L *lua.LState) int {
	jsonStr := L.ToString(1)
	var table map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	luaTable := L.NewTable()
	for k, v := range table {
		switch v := v.(type) {
		case string:
			luaTable.RawSetString(k, lua.LString(v))
		case int:
			luaTable.RawSetString(k, lua.LNumber(int64(v)))
		case float64:
			luaTable.RawSetString(k, lua.LNumber(v))
		case bool:
			luaTable.RawSetString(k, lua.LBool(v))
		case nil:
			luaTable.RawSetString(k, lua.LNil)
		default:
			L.Push(lua.LNil)
			L.Push(lua.LString("unsupported type"))
			return 2
		}
	}
	L.Push(luaTable)
	return 1
}

func HtmlBySelector(html string) string {
	return ""
}
