package main

import (
	"embed"
	"fmt"
	"os"
	"z0ne/cmd"
)

//go:embed all:assets
var assets embed.FS

func main() {
	var banner, err = os.ReadFile("assets/banners/z0ne.txt")
	if err != nil {
		fmt.Println("Error reading banner:", err)
	}
	fmt.Println(string(banner))
	cmd.Execute()
}
