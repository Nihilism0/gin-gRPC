package main

import (
	"gin_demo/api"
	"gin_demo/dao"
)

func main() {
	dao.DataLoad()
	api.InitRouter()
}
