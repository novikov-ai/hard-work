# Как правильно готовить юнит-тесты

## На примере LRU-кэша напишем юнит-тесты, выделяя и тестируя свойства, а не сценарии

### Тестируемые методы
~~~go
// ...
func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()

	item, exists := cache.items[key]

	if !exists {
		cache.queue.PushFront(cacheItem{key: key, value: value})
		cache.items[key] = cache.queue.Front()

		deleteUnusedIfOverflow(cache)
	} else {
		updatingCacheItem := item.Value.(cacheItem)
		updatingCacheItem.value = value

		item.Value = updatingCacheItem

		cache.queue.MoveToFront(item)
	}

	cache.mutex.Unlock()

	return exists
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item, exists := cache.items[key]
	var value interface{}

	if exists {
		cache.queue.MoveToFront(item)
		value = item.Value.(cacheItem).value
	}

	return value, exists
}
// ...
~~~

### Множество свойств, которые задают корректность LRU-кэша

1. Если элемент устанавливают в кэш и его раньше там не было, то он помещается в кэш на первое место
2. Если элемент устанавливают в кэш и он был там раньше, то он помещается в кэш на первое место
3. Если элемент устанавливают в кэш и происходит переполнение кэша, то самый старый элемент удаляется
4. Если элемент пытаются получить из кэша и он там есть, то он возвращается и помещается на первое место
5. Если элемент пытаются получить из кэша и его там нет, то он не возвращается

Каждое свойство лучше всего протестировать модульным тестом (юнитом) 



~~~go
func TestCache(t *testing.T) {
	t.Run("set new value to cache", func(t *testing.T) {
		cache := NewCache(10) // cache with capacity = 10

        newValue := "aaa"

		_, ok := cache.Get(newValue)
		require.False(t, ok)

		cache.Set(newValue)

        _, ok := cache.Get(newValue)
		require.True(t, ok)

        require.Equal(t, 10, cache.Front())
	})

	t.Run("set existing value to cache", func(t *testing.T) {
		cache := NewCache(10) // cache with capacity = 10

        existingValue := "aaa"

		_, ok = cache.Get(existingValue)
        require.True(t, ok)

        cache.Set(existingValue)

        require.Equal(t, 10, cache.Front())
	})

    t.Run("remove old element", func(t *testing.T) {
		cache := NewCache(10) // cache with capacity = 10

        newValue := "aaa"

        cache.Back() // out: "old"

        cache.Set(newValue)

        require.Equal(t, "replaced-old", cache.Back())
	})

    t.Run("get existing element", func(t *testing.T) {
		cache := NewCache(10) // cache with capacity = 10

        existingValue := "aaa"

		value, ok = cache.Get(existingValue)

        require.True(t, ok)
        require.Equal(t, "aaa", value)
        require.Equal(t, "aaa", cache.Front())
	})

     t.Run("get non existent element", func(t *testing.T) {
		cache := NewCache(10) // cache with capacity = 10

        nonExistentValue := "dddd"

		value, ok = cache.Get(nonExistentValue)
        
        require.False(t, ok)
	})
}
~~~