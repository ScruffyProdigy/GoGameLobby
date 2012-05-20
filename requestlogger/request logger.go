package requestlogger

import (
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/TheRack/httper"
)

type Middleware struct{}

func (Middleware) Run(vars map[string]interface{}, next func()) {
	r := (httper.V)(vars).GetRequest()
	(logger.V)(vars).Get().Println(r.Method, r.URL.String())
	//	(logger.V)(vars).Get().Println(r)
	next()
}

var M Middleware
