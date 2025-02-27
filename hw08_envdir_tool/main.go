package main

import (
	"fmt"

	"github.com/MaksimIschenko/hw_otus_golang/hw08_envdir_tool/envreader"
)

func main() {
	envMap, err := envreader.ReadDir("./testdata/env")
	if err != nil {
		panic(err)
	}
	for k, v := range envMap {
		fmt.Printf("%s: %s %v\n", k, v.Value, v.NeedRemove)
	}
}
