package controller

import (
	"fmt"
	"net/http"
)

// ControllerInterface is an interface to uniform all controller handler.
type ControllerInterface interface {
	Init(http.ResponseWriter, *http.Request)
	Prepare()
	Finish()
	Render() error
}

type BaseController struct {
	R       *http.Request
	W       http.ResponseWriter
	TplName string
	Data    interface{}
}

func (c *BaseController) Init(w http.ResponseWriter, r *http.Request) {
	c.R = r
	c.W = w
}

func (c *BaseController) Prepare() {
	fmt.Println(c.R, c.W)
	fmt.Println("BaseController", "Prepare")
}

func (c *BaseController) Finish() {
	fmt.Println("BaseController", "Finish")
}

func (c *BaseController) Render() error {
	ExecuteViewPathTemplate(c.W, c.TplName, c.TplName, c.Data)

	return nil
}
