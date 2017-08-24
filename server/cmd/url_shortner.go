package main

import (
	"github.com/ngalayko/url_shortner/server"
)

func main() {
	app := server.NewApplication()

	app.Healthcheck()
}
