package controllers

import (
	"github.com/HairyMezican/Middleware/renderer"
	"github.com/HairyMezican/Middleware/router"
)

var Root = router.BasicRoute("", renderer.Renderer{"test"})
