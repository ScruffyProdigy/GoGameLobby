package controller

import (
	"../rack"
	"net/http"
	"reflect"
	"strings"
	"../redirecter"
	"../templater"
)


type dispatchAction struct {
	m ModelMap
	name string
	action rack.Middleware
}

func (this dispatchAction) Run(r *http.Request,vars rack.Vars,next rack.NextFunc) (int,http.Header,[]byte) {
	actions := rack.NewRack()
	if r.Method == "POST" || r.Method == "PUT" {
		actions.Add(rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {

			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
			return next()
		}))
	}
	actions.Add(this.action)
	switch(r.Method) {
	case "GET":
		//if it was a get, the default action should be to render the template corresponding with the action
		actions.Add(templater.Templater{this.m.RouteName()+"/"+this.name})
	case "POST","PUT":
		//if it was a put or a post, we the default action should be to redirect to the affected item
		actions.Add(rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
			return redirecter.Go(r, vars, vars[this.m.VarName()].(Urler).Url())
		}))
	case "DELETE":
		//I'm not currently sure what the default action for deletion should be, perhaps redirecting to the parent route
	default:
		panic("Unknown method")
	}
	return actions.Run(r,vars,next)
}

func isControlFunc(m reflect.Method) bool{
	t := m.Type
	if t.Kind() != reflect.Func {	//it should be a function
		return false
	}
	if t.NumIn() != 1 {			//it should have one input parameter (the 'this' controller)
		return false
	}
	if t.NumOut() != 1 {		//it should have one output parameter
		return false
	}
	if t.Out(0).String() != "controller.Response" {		//the output should be a controller.Response
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
	for _,funcName := range([]string{"Index","Create","New","Show","Edit","Update","Destroy"}) {
		//check each function to make sure it's there
		method,methodExists := controllerType.MethodByName(funcName)
		if methodExists && isControlFunc(method) {
			caller := method.Func
			value := []reflect.Value{reflect.ValueOf(controller)}
			restfuncs[funcName] = &dispatchAction{m:mapper,name:strings.ToLower(funcName),action:rack.Func(func(r *http.Request,vars rack.Vars,next rack.NextFunc)(int,http.Header,[]byte){
				mapper.SetDefaultResponse(next)
				result,isResponse := caller.Call(value)[0].Interface().(Response)
				if !isResponse {
					panic("unexpected output")
				}
				return result.ToRack()
			})}
		}
	}
	
	return
}

type mapList struct {
	all,get,put,post,delete map[string]rack.Middleware
}

func GetGenericMapList(controller interface{}, functype string) (funcs mapList) {
	funcs.all = GetGenericMap(controller,functype)
	funcs.get = GetGenericMap(controller,"Get"+functype)
	funcs.put = GetGenericMap(controller,"Put"+functype)
	funcs.post = GetGenericMap(controller,"Post"+functype)
	funcs.delete = GetGenericMap(controller,"Delete"+functype)
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
	for i,c := 0,controllerType.NumMethod();i < c;i = i+1 {
		method := controllerType.Method(i)
		if len(method.Name) >= typelen && method.Name[:typelen] == functype && isControlFunc(method) {
			caller := method.Func
			value := []reflect.Value{reflect.ValueOf(controller)}
			name := method.Name[typelen:]
			funcs[name] = &dispatchAction{m:mapper,name:strings.ToLower(name),action:rack.Func(func(r *http.Request,vars rack.Vars,next rack.NextFunc)(int,http.Header,[]byte){
				mapper.SetDefaultResponse(next)
				result := caller.Call(value)
				return result[0].Interface().(int),result[1].Interface().(http.Header),result[2].Interface().([]byte)
			})}
		}
	}
	return
}