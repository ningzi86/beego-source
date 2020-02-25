package controller

import (
	"fmt"
	"os"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.TplName = "index.html"

	dir, _ := os.Getwd()
	fs := []string{
		dir + "/templates/layout.html",
		dir + "/templates/layout/body.html",
		dir + "/templates/layout/css.html",
		dir + "/templates/layout/scripts.html",
		dir + "/templates/index.html",
	}

	data := make(map[string]interface{})
	data["fs"] = fs
	data["Title"] = "hello"

	c.Data = data

	fmt.Println("IndexController", "Index")
}

func (c *IndexController) Prepare() {
	fmt.Println("IndexController", "Prepare", c.R, c.W)
}
