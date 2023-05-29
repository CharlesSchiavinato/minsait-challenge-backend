package cache

type Cache interface {
	Check() error
	Close() error
}
