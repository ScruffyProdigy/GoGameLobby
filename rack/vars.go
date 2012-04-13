package rack

type Vars map[string]interface{}
type VarFunc func(Vars) interface{}

func (this Vars) Apply(f VarFunc) interface{} {
	return f(this)
}
