package cache

import (
	"somefun/pkg/log"
	"sync"
)

type Cache interface {
	Load(chatId int) (interface{}, bool)
	Store(chatId int, value interface{}) error
	View()
}

var cacheOnce sync.Once

var cache *memoryCache

func GetCache() Cache {
	cacheOnce.Do(func() {
		cache = &memoryCache{
			repo: make(map[int]interface{}),
		}
		go cleanCache(cache)
	})
	return cache
}

func cleanCache(cache *memoryCache) {
	// todo 清理长时间未使用的chatId
	for {
	}
}

type memoryCache struct {
	sync.Mutex
	repo map[int]interface{}
}

func (mc *memoryCache) Load(chatId int) (interface{}, bool) {
	mc.Lock()
	defer func() {
		mc.Unlock()
	}()
	v, ok := mc.repo[chatId]
	if !ok {
		return nil, false
	}
	return v, true
}

func (mc *memoryCache) Store(chatId int, value interface{}) error {
	mc.Lock()
	defer func() {
		mc.Unlock()
	}()
	mc.repo[chatId] = value
	return nil
}

func (mc *memoryCache) View() {
	log.Info("memoryCache is ", mc.repo)
}
