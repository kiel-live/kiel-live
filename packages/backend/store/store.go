package store

type Store interface {
	Load() error
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
	Unload() error
}
