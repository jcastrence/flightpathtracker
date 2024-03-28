package main

import (
	"fmt"

	"github.com/jcastrence/flightpathtracker/src/router"
)

func main() {
	fmt.Println("Starting Flight Path Tracker microservice...")

	e := router.New()

	e.Logger.Fatal(e.Start(":8080"))
}
