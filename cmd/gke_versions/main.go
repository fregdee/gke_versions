package main

import (
	"fmt"
	"os"

	"github.com/fregdee/gke_versions"
)

func main() {
	err := gke_versions.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
