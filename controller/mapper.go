package controller

import (
	"github.com/HairyMezican/Middleware/redirecter"
	"github.com/HairyMezican/Middleware/renderer"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"reflect"
	"strings"
)

type dispatchAction struct {
	m      ModelMap
	name   string
	action rack.Middleware
}

type Urler interface {
	Url() string
}

func (this dispatchAction) Run(vars map[string]interface{}, next func()) {
	actions := rack.New()
	actions.Add(this.action)
	switch (httper.V)(vars).GetRequest().Method {
	case "GET":
		//if it was a get, the default action should be to render the template corresponding with the action
		actions.Add(renderer.Renderer{this.m.RouteName() + "/" + this.name})
	case "POST", "PUT":
		//if it was a put or a post, we the default action should be to redirect to the affected item
		actions.Add(rack.Func(func(vars map[string]interface{}, next func()) {
			urler, isUrler := vars[this.m.VarName()].(Urler)
			if !isUrler {
				panic("Object doesn't have an URL to direct to")
			}
			(redirecter.V)(vars).Redirect(urler.Url())
		}))
	case "DELETE":
		//I'm not currently sure what the default action for deletion should be, perhaps redirecting to the parent route
	default:
		panic("Unknown method")
	}
	actions.Run(vars, next)
}

func isControlFunc(m reflect.Method) bool {
	t := m.Type
	if t.Kind() != reflect.Func { //it should be a function
		return false
	}
	if t.NumIn() != 1 { //it should have one input parameter (the 'this' controller)
		return false
	}
	if t.NumOut() != 0 { //it should have no output parameters
		return false
	}
	return true
}

func GetRestMap(controller interface{}) (restfuncs map[string]rack.Middleware) {
	restfuncs = make(map[string]rack.Middleware)

	mapper, canMap := controller.(ModelMap)
	if !canMap {
		panic("Can't set rack variables!")
	}

	controllerType := reflect.TypeOf(controller)
	for _, funcName := range []string{"Index", "Create", "New", "Show", "Edit", "Update", "Destroy"} {
		//check each function to make sure it's there
		method, methodExists := controllerType.MethodByName(funcName)
		if methodExists && isControlFunc(method) {
			caller := method.Func
			value := []reflect.Value{reflect.ValueOf(controller)}
			restfuncs[funcName] = &dispatchAction{
				m:    mapper,
				name: strings.ToLower(funcName),
				action: rack.Func(func(vars map[string]interface{}, next func()) {
					mapper.SetFinish(next)
					caller.Call(value)
				}),
			}
		}
	}

	return
}

type mapList struct {
	all, get, put, post, delete map[string]rack.Middleware
}

func GetGenericMapList(controller interface{}, functype string) (funcs mapList) {
	funcs.all = GetGenericMap(controller, functype)
	funcs.get = GetGenericMap(controller, "Get"+functype)
	funcs.put = GetGenericMap(controller, "Put"+functype)
	funcs.post = GetGenericMap(controller, "Post"+functype)
	funcs.delete = GetGenericMap(controller, "Delete"+functype)
	return
}

func GetGenericMap(controller interface{}, functype string) (funcs map[string]rack.Middleware) {
	funcs = make(map[string]rack.Middleware)
	typelen := len(functype)

	mapper, canMap := controller.(ModelMap)
	if !canMap {
		panic("Can't set rack variables!")
	}

	controllerType := reflect.TypeOf(controller)
	for i, c := 0, controllerType.NumMethod(); i < c; i = i + 1 {
		method := controllerType.Method(i)
		if len(method.Name) >= typelen {
			if method.Name[:typelen] == functype && isControlFunc(method) {
				caller := method.Func
				value := []reflect.Value{reflect.ValueOf(controller)}
				name := method.Name[typelen:]
				funcs[name] = &dispatchAction{
					m:    mapper,
					name: strings.ToLower(name),
					action: rack.Func(func(vars map[string]interface{}, next func()) {
						mapper.SetFinish(next)
						caller.Call(value)
					}),
				}
			}
		}
	}
	return
}
