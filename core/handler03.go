package core

import (
	"fmt"
	"inaction/gitchat/router"
	"net"
	"net/http"
	"reflect"
	"time"
)

type ControllerRegister struct {
}

func (c *ControllerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.RequestURI
	routerInfo, err := router.FindRouter(path)

	if err != nil {
		fmt.Fprintf(w, "获取Controller出错： %s", err.Error())
		return
	}

	execController := routerInfo.Initialize()

	vc := reflect.ValueOf(execController)
	method := vc.MethodByName("Index")
	method.Call(nil)

	execController.Init(w, r)

	execController.Prepare()
	execController.Render()
	execController.Finish()

}

// NewControllerRegister returns a new ControllerRegister.
func NewControllerRegister() *ControllerRegister {
	cr := &ControllerRegister{}
	return cr
}

func NewApp() *App {
	cr := NewControllerRegister()
	app := &App{Handlers: cr, Server: &http.Server{}}
	return app
}

type App struct {
	Handlers *ControllerRegister
	Server   *http.Server
}

func Handler03() {

	app := NewApp()
	app.Server.Handler = app.Handlers
	app.Server.Addr = ":8080"

	ln, err := net.Listen("tcp4", app.Server.Addr)
	if err != nil {
		time.Sleep(100 * time.Microsecond)
	}
	if err = app.Server.Serve(ln); err != nil {
		time.Sleep(100 * time.Microsecond)
	}

}
