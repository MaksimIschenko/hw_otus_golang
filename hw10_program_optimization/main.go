package main

import (
	"fmt"
	"os"

	"github.com/MaksimIschenko/hw_otus_golang/hw10_program_optimization/stats"
)

func main() {
	f, err := os.OpenFile("./testdata/users.dat", os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	statsDomain, err := stats.GetDomainStat(f, "info")
	if err != nil {
		fmt.Println("stats.GetDomainStat error:", err)
	}

	for k, v := range statsDomain {
		fmt.Println(k, v)
	}
}
