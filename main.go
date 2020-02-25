package main

import (
	"github.com/astaxie/beego"
	"inaction/gitchat/core"
)

func main() {
	core.Handler03()

	beego.Run()
}
