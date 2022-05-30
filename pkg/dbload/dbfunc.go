package dbload

type Dface interface {
	GetKey
}

type GetKey interface {
	GetKeyWithPrefix(k string) ([]string, error)
}
