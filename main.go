package main

import (
	"MyGram/config"
	"MyGram/router"
	"log"
)

func main() {
    config.Connect()

    app := router.Uri()

    log.Fatal(app.Listen(":3000"))
}