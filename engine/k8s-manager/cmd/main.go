package main

import (
	"log"

	"github.com/konstellation-io/kai/engine/k8s-manager/cmd/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
