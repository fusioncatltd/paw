package main

import (
	"context"
	"github.com/fusioncatalyst/paw/router"
	"log"
	"os"
)

func main() {
	cli := router.GetCLIRouter()
	if err := cli.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
