package mutex

type Mutex interface {
	Try(action func()) bool
	Force(action func())
}
