package main

import (
	"embed"
	"fmt"
	"time"

	"github.com/zer0ne-hub/z0ne/cmd"
)

//go:embed all:assets
var Assets embed.FS

func main() {
	//profiler
	start := time.Now()
	var banner, err = Assets.ReadFile("assets/banners/z0ne.txt")
	if err != nil {
		fmt.Println("Error reading banner:", err)
	}
	fmt.Println(string(banner))
	cmd.Execute()
	fmt.Println("Executed in:", time.Since(start))
}
