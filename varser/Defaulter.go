package varser

type Default map[string]interface{}

func (this Default) Run(vars map[string]interface{}, next func()) {
	for k, v := range this {
		vars[k] = v
	}
	next()
}

type Override map[string]interface{}

func (this Override) Run(vars map[string]interface{}, next func()) {
	next()
	for k, v := range this {
		vars[k] = v
	}
}
