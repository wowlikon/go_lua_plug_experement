package main

import (
	"fmt"
	"os"
	"path/filepath"

	lua "github.com/cathalgarvey/gopher-lua"
)

func Square(L *lua.LState) int {
	lv := L.ToInt(1)             /* get argument */
	L.Push(lua.LNumber(lv * lv)) /* push result */
	return 1                     /* number of results */
}

func loadLuaFile(state *lua.LState, path string) {
	if err := state.DoFile(path); err != nil {
		fmt.Println("Error loading library:", err)
	}
}

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
		state.Register("square", Square)

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
