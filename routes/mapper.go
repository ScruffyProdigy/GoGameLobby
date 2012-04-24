package routes

import (
	"../rack"
	"net/http"
)

type Newer interface {
	New(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type NewWrapper struct {
	Newer
}

func (this NewWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.New(r, vars, next)
}

type Creater interface {
	Create(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type CreateWrapper struct {
	Creater
}

func (this CreateWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.Create(r, vars, next)
}

type Indexer interface {
	Index(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type IndexWrapper struct {
	Indexer
}

func (this IndexWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.Index(r, vars, next)
}

type Shower interface {
	Show(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type ShowWrapper struct {
	Shower
}

func (this ShowWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.Show(r, vars, next)
}

type Editer interface {
	Edit(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type EditWrapper struct {
	Editer
}

func (this EditWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.Edit(r, vars, next)
}

type Updater interface {
	Update(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type UpdateWrapper struct {
	Updater
}

func (this UpdateWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.Update(r, vars, next)
}

type Deleter interface {
	Delete(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte)
}

type DeleteWrapper struct {
	Deleter
}

func (this DeleteWrapper) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return this.Delete(r, vars, next)
}

func GetRestMap(c interface{}) (restfuncs map[string]rack.Middleware) {
	restfuncs = make(map[string]rack.Middleware)
	mIndex, hasIndex := c.(Indexer)
	if hasIndex {
		restfuncs["index"] = IndexWrapper{mIndex}
	}
	mCreate, hasCreate := c.(Creater)
	if hasCreate {
		restfuncs["create"] = CreateWrapper{mCreate}
	}
	mNew, hasNew := c.(Newer)
	if hasNew {
		restfuncs["new"] = NewWrapper{mNew}
	}
	mShow, hasShow := c.(Shower)
	if hasShow {
		restfuncs["show"] = ShowWrapper{mShow}
	}
	mEdit, hasEdit := c.(Editer)
	if hasEdit {
		restfuncs["edit"] = EditWrapper{mEdit}
	}
	mUpdate, hasUpdate := c.(Updater)
	if hasUpdate {
		restfuncs["update"] = UpdateWrapper{mUpdate}
	}
	mDelete, hasDelete := c.(Deleter)
	if hasDelete {
		restfuncs["delete"] = DeleteWrapper{mDelete}
	}
	return
}
