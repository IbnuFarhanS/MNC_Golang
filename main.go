package main

import "github.com/IbnuFarhanS/Golang_MNC/api"

func main() {
	app := api.NewApp()
	app.Initialize()
	app.Run("8080")
}
