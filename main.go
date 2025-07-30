package main

import (
	"yanblog/routers"
	"yanblog/model"
)

func main() {
	model.InitDB()
	routers.InitRouter()

}
