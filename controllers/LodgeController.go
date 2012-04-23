package controller

import (
	"../login"
	"../models/lodge"
	"../models"
	"../rack"
	"../redirecter"
	"../routes"
	"../session"
	"net/http"
)


type LodgeController struct {
	L *lodge.LodgeCollection
}

func (LodgeController) RouteName() string {
	return "lodges"
}

func (LodgeController) VarName() string {
	return "Lodge"
}

func (this LodgeController) Indexer(s string) (interface{},bool) {
	result := this.L.LodgeFromName(s)
	return result,result!=nil
}


func (this LodgeController) Index(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	var lodges []lodge.Lodge
	err := this.L.AllLodges(&lodges)
	if err != nil {
		panic(err)
	}

	vars["Lodges"] = lodges
	vars["Title"] = "Mason Lodges"

	return next()
}

func (this LodgeController) Show(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	l := vars["Lodge"].(*lodge.Lodge)
	
	vars["Title"] = l.Name

	return next()
}

func (this LodgeController) New(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {

	vars["Title"] = "Create a Mason Lodge"

	return next()
}

func (this LodgeController) Create(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	var l lodge.Lodge
	defer func() {
		rec := recover()
		if rec != nil {
			status,header,message = redirecter.Go(r,vars,"/lodges/new",
				session.AddFlash("You fucked something up, please try again"))
		}
	}()

	l.Name = r.FormValue("Lodge[Name]")
	l.AddMason(login.CurrentUser(vars))
	
	model.Save(&l)

	return redirecter.Go(r,vars,l.Url())
}


func init() {
	routes.Resource(LodgeController{lodge.L}).AddTo(routes.Root)
}
