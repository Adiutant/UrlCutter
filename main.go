package main

import "shortUrl/http_server"

func main() {
	server, err := http_server.MakeUrlServer()
	if err != nil {
		return
	}
	server.SetRoutes()
	err = server.Run()
	if err != nil {
		return
	}
}
