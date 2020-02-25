package router

import (
	"errors"
	"inaction/gitchat/controller"
	"net/http"
	"reflect"
)

type ControllerInfo struct {
	pattern        string
	controllerType reflect.Type
	methods        map[string]string
	handler        http.Handler
	routerType     int
	Initialize     func() controller.ControllerInterface
}

var controllers map[string]*ControllerInfo

func Router(pattern string, c controller.ControllerInterface) {

	if (controllers == nil) {
		controllers = make(map[string]*ControllerInfo)
	}


	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	methods := make(map[string]string)

	route := &ControllerInfo{}
	route.pattern = pattern
	route.methods = methods
	route.controllerType = t
	route.Initialize = func() controller.ControllerInterface {
		vc := reflect.New(route.controllerType)
		execController, ok := vc.Interface().(controller.ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface")
		}
		return execController
	}

	controllers[pattern] = route

}

func FindRouter(path string) (*ControllerInfo, error) {
	if v, ok := controllers[path]; ok {
		return v, nil
	}
	return nil, errors.New("没有找到Controller")
}
