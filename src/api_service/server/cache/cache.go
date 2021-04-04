package cache

import (
	"container/list"
	"sync"
	"time"
)

// Store contains an LRU Cache
type Store struct {
	mutex *sync.Mutex
	store map[string]*list.Element
	listStore    *list.List
	buffer   int // Zero for unlimited
}

// Node maps a value to a key
type Node struct {
	key     string
	value   string
	expire  int64  // Unix time
}

var s *Store


func NewCache(size int) *Store {
	s = New(size)
	return s
}

// Create new cache
func New(buffer int) *Store {
	s := &Store{
		mutex: &sync.Mutex{},
		store: make(map[string]*list.Element),
		listStore:    list.New(),
		buffer:   buffer,
	}
	return s
}

// Get a key from cache
func (s *Store) Get(key string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	current, exist := s.store[key]
	if exist {
		expire := int64(current.Value.(*Node).expire)
		if expire == 0 || expire > time.Now().Unix() {
			s.listStore.MoveToFront(current)
			return current.Value.(*Node).value, true
		}
	}
	return "", false
}

// Insert key to cache
func (s *Store) Insert(key string, value string, expire int64) {
	current, exist := s.store[key]
	if exist != true {
		s.store[key] = s.listStore.PushFront(&Node{
			key:    key,
			value:  value,
			expire: expire,
		})
		if s.buffer != 0 && s.listStore.Len() > s.buffer {
			s.Delete(s.listStore.Remove(s.listStore.Back()).(*Node).key)
		}
		return
	}
	current.Value.(*Node).value = value
	current.Value.(*Node).expire = expire
	s.listStore.MoveToFront(current)
}

// delete key from cache
func (s *Store) Delete(key string) {
	current, exist := s.store[key]
	if exist != true {
		return
	}
	s.listStore.Remove(current)
	delete(s.store, key)
}

// Flush all keys
func (s *Store) Empty() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store = make(map[string]*list.Element)
	s.listStore = list.New()
}
