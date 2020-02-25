package router

import "inaction/gitchat/controller"

func init() {
	Router("/", &controller.IndexController{})
}
