package cache

type LruCache interface {
	NewCache(size int)
	Get(key string) (string, bool)
	Insert(key string, value string, expire int64)
	Delete(key string)
	Empty()
}
