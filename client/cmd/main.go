package main

import "github.com/ngalayko/url_shortner/client"

func main() {
	a := client.NewApplication()

	a.Serve()
}
