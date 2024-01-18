package main

import (
	"backend/api"
	"backend/db"
)

// entry point
func main() {
	db.Init()
	server := api.NewServer(":3000")
	server.Start()
}
