# Когда наследование лучше композиции?

## LSP в тестах

### Пример 1. Интерфейс Cache

```go
// cache.go
package cache

// Cache — контракт для key‑value хранилища
type Cache interface {
    Get(key string) (value []byte, ok bool)
    Set(key string, value []byte)
    Delete(key string)
}
```

```go
// cache_lsp_test.go
package cache_test

import (
    "testing"

    "github.com/yourorg/yourrepo/cache"
    "github.com/stretchr/testify/assert"
)

// runCacheTests — общий набор LSP‑тестов для любого cache.Cache
func runCacheTests(t *testing.T, c cache.Cache) {
    // 1) При установке ключа — Get возвращает точно тот же байт‑слайс
    c.Set("foo", []byte("bar"))
    v, ok := c.Get("foo")
    assert.True(t, ok)
    assert.Equal(t, []byte("bar"), v)

    // 2) После Delete ключ уже не должен присутствовать
    c.Delete("foo")
    _, ok = c.Get("foo")
    assert.False(t, ok)
}
```

```go
// memory_cache_test.go
package cache_test

import (
    "testing"

    "github.com/yourorg/yourrepo/cache"
)

// TestMemoryCache запускает LSP‑тесты для in‑memory реализации
func TestMemoryCache(t *testing.T) {
    mem := cache.NewMemoryCache()      // implements cache.Cache
    runCacheTests(t, mem)              // «наследуем» все тесты
    // можно добавить ещё специфичные тесты:
    t.Run("MemoryCache_Capacity", func(t *testing.T) {
        // проверяем, что capacity не превышается
    })
}
```

```go
// redis_cache_test.go
package cache_test

import (
    "testing"

    "github.com/yourorg/yourrepo/cache"
)

// TestRedisCache запускает тот же набор LSP‑тестов для Redis
func TestRedisCache(t *testing.T) {
    rds := cache.NewRedisCache("localhost:6379")
    runCacheTests(t, rds)
    // плюс тесты на отказоустойчивость и reconnect
}
```

Здесь `runCacheTests` — «тест суперкласса», автоматически применяемый к каждому конкретному кэшу.

---

### Пример 2. Интерфейс RateLimiter

```go
// ratelimit.go
package ratelimit

import "time"

// RateLimiter — контракт для ограничителя скорости
type RateLimiter interface {
    Allow() bool       // можно ли сейчас выполнить операцию
    Reset()            // сброс внутреннего состояния
}
```

```go
// ratelimit_lsp_test.go
package ratelimit_test

import (
    "testing"
    "time"

    "github.com/yourorg/yourrepo/ratelimit"
    "github.com/stretchr/testify/assert"
)

// runRateLimiterTests проверяет, что любая реализация RateLimiter
// не пропускает более N событий в секунду и корректно сбрасывается.
func runRateLimiterTests(t *testing.T, rl ratelimit.RateLimiter) {
    rl.Reset()
    // допустим, лимит 3 события в секунду
    allowed := 0
    for i := 0; i < 5; i++ {
        if rl.Allow() {
            allowed++
        }
    }
    assert.Equal(t, 3, allowed, "should allow exactly 3 events per second")

    // имитируем паузу в 1 секунду — после Reset должно снова пускать
    rl.Reset()
    time.Sleep(1 * time.Second)
    assert.True(t, rl.Allow(), "after reset и sleep limiter должен пускать запрос")
}
```

```go
// token_bucket_test.go
package ratelimit_test

import (
    "testing"
    "github.com/yourorg/yourrepo/ratelimit"
)

// TestTokenBucketLimiter применяет общий набор тестов к TokenBucket
func TestTokenBucketLimiter(t *testing.T) {
    // допустим, NewTokenBucket(3, 1s)
    tb := ratelimit.NewTokenBucketLimiter(3, time.Second)
    runRateLimiterTests(t, tb)
}
```

```go
// leaky_bucket_test.go
package ratelimit_test

import (
    "testing"
    "github.com/yourorg/yourrepo/ratelimit"
)

// TestLeakyBucketLimiter — тот же «суперкласс тестов» для LeakyBucket
func TestLeakyBucketLimiter(t *testing.T) {
    lb := ratelimit.NewLeakyBucketLimiter(3, time.Second)
    runRateLimiterTests(t, lb)
}
```

В этом примере любой новый лимитер автоматически проверяется на соблюдение LSP‑контракта (ограничение и сброс), без дублирования кода в каждом `*_test.go`.

## Удачные примеры подтипного полиморфизма

### 1. Система очередей задач

```go
// BaseTaskQueue.go
type Task interface {
    Execute() error
}

type TaskQueue interface {
    Enqueue(Task)
    Dequeue() Task
    Len() int
}

// DefaultQueue реализует базовый FIFO‑queue
type DefaultQueue struct {
    tasks []Task
}

func (q *DefaultQueue) Enqueue(t Task) { q.tasks = append(q.tasks, t) }
func (q *DefaultQueue) Dequeue() Task {
    if len(q.tasks) == 0 { return nil }
    t := q.tasks[0]
    q.tasks = q.tasks[1:]
    return t
}
func (q *DefaultQueue) Len() int { return len(q.tasks) }

// LoggingQueue «наследует» DefaultQueue и добавляет логгирование
type LoggingQueue struct {
    *DefaultQueue // embedding для повторного использования реализации
}

func NewLoggingQueue() *LoggingQueue {
    return &LoggingQueue{DefaultQueue: &DefaultQueue{}}
}

func (q *LoggingQueue) Enqueue(t Task) {
    log.Printf("Enqueue task: %T", t)
    q.DefaultQueue.Enqueue(t)
}

func (q *LoggingQueue) Dequeue() Task {
    t := q.DefaultQueue.Dequeue()
    log.Printf("Dequeue task: %T", t)
    return t
}

// В рантайме маршрутизатор видит TaskQueue и может подставлять любой подкласс
func ProcessAll(q TaskQueue) {
    for q.Len() > 0 {
        if err := q.Dequeue().Execute(); err != nil {
            log.Println("error:", err)
        }
    }
}
```

* `LoggingQueue` расширяет, но не меняет базовое поведение: LSP соблюдается, `ProcessAll` работает с обоими типами одинаково.

---

### 2. HTTP‑мидлвары с общим контрактом

```go
// middleware.go
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// BaseHandler предоставляет общую логику (напр., парсинг запроса)
type BaseHandler struct{}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // common pre‑processing...
    next := h.Next()
    next.ServeHTTP(w, r)
}

// Метод-«заглушка», чтобы embedding‑наследники предоставили реальный Next
func (h *BaseHandler) Next() Handler {
    panic("Next() must be overridden")
}

// AuthHandler наследует BaseHandler и добавляет проверку токена
type AuthHandler struct {
    *BaseHandler
    next Handler
}

func (a *AuthHandler) Next() Handler { return a.next }

func NewAuthHandler(next Handler) *AuthHandler {
    return &AuthHandler{
        BaseHandler: &BaseHandler{},
        next:        next,
    }
}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")
    if !validateToken(token) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    a.BaseHandler.ServeHTTP(w, r)
}

// В цепочке мидлваров каждый элемент — Handler, полиморфно вызывается
mux.Handle("/secure", NewAuthHandler(myBusinessLogicHandler))
```

* Все `Handler` взаимозаменяемы: маршрутизатор не замечает разницы между «чистым» и «авторизующим» хэндлером.

---

### 3. Коннекторы к разным БД

```go
// db.go
type DB interface {
    Query(query string, args ...interface{}) (*Rows, error)
    Exec(query string, args ...interface{}) (Result, error)
}

// BaseDB реализует общие вспомогательные методы
type BaseDB struct {
    driverName string
    dsn        string
    conn       *sql.DB
}

func (b *BaseDB) Connect() error {
    db, err := sql.Open(b.driverName, b.dsn)
    if err != nil { return err }
    b.conn = db
    return nil
}

// MySQLDB наследует BaseDB, задаёт параметры и при необходимости расширяет
type MySQLDB struct{ *BaseDB }

func NewMySQLDB(dsn string) *MySQLDB {
    return &MySQLDB{BaseDB: &BaseDB{driverName: "mysql", dsn: dsn}}
}

// При необходимости переопределяем Exec, чтобы логировать медленные запросы
func (m *MySQLDB) Exec(query string, args ...interface{}) (Result, error) {
    start := time.Now()
    res, err := m.conn.Exec(query, args...)
    if time.Since(start) > 100*time.Millisecond {
        log.Printf("SLOW QUERY: %s", query)
    }
    return res, err
}

// Везде, где ожидается DB, можно подставить MySQLDB, PostgreSQLDB и т.д.
func RunMigration(db DB) error {
    _, err := db.Exec("CREATE TABLE ...")
    return err
}
```

* Расширяем, но не меняем контракт `DB`: LSP соблюдается.
