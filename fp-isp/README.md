# ISP с т.зр. ФП

## Пример 1: Разделение интерфейса репозитория

**Было:**
```go
// Нарушает ISP - содержит методы для разных сущностей
type Repository interface {
    // User methods
    CreateUser(ctx context.Context, user *User) error
    GetUserByID(ctx context.Context, id string) (*User, error)
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id string) error
    
    // Order methods
    CreateOrder(ctx context.Context, order *Order) error
    GetOrderByID(ctx context.Context, id string) (*Order, error)
    ListOrdersByUser(ctx context.Context, userID string) ([]Order, error)
    
    // Product methods
    CreateProduct(ctx context.Context, product *Product) error
    GetProductByID(ctx context.Context, id string) (*Product, error)
}
```

**Стало (соблюдаем ISP и обобщенность):**
```go
// Базовый интерфейс CRUD операций - максимально обобщенный
type CRUDRepository[T any, ID comparable] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id ID) error
}

// Специализированные интерфейсы для конкретных сценариев
type UserRepository interface {
    CRUDRepository[User, string]
    GetByEmail(ctx context.Context, email string) (*User, error)
}

type OrderRepository interface {
    CRUDRepository[Order, string]
    ListByUser(ctx context.Context, userID string) ([]Order, error)
    ListByStatus(ctx context.Context, status string) ([]Order, error)
}

type ProductRepository interface {
    CRUDRepository[Product, string]
    Search(ctx context.Context, query string) ([]Product, error)
    ListByCategory(ctx context.Context, category string) ([]Product, error)
}
```

**Как сочетаем:** 
- `CRUDRepository` - обобщенный интерфейс с единственной ответственностью (CRUD операции)
- Специализированные интерфейсы добавляют только специфичные методы
- Клиенты зависят только от нужных им интерфейсов

## Пример 2: Разделение интерфейса обработчика HTTP

**Было:**
```go
// Монолитный интерфейс обработчика
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
    ValidateRequest(*http.Request) error
    LogRequest(*http.Request)
    RateLimit(*http.Request) error
    Authenticate(*http.Request) (*User, error)
    Authorize(*User, *http.Request) bool
}
```

**Стало:**
```go
// Базовый интерфейс обработчика HTTP - минимальный и обобщенный
type HTTPHandler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}

// Декораторы как отдельные интерфейсы
type RequestValidator interface {
    ValidateRequest(*http.Request) error
}

type Logger interface {
    LogRequest(*http.Request)
}

type RateLimiter interface {
    Allow(*http.Request) bool
}

type Authenticator interface {
    Authenticate(*http.Request) (*User, error)
}

type Authorizer interface {
    Authorize(*User, *http.Request) bool
}

// Композиция через декораторы
func WithValidation(h HTTPHandler, v RequestValidator) HTTPHandler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := v.ValidateRequest(r); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        h.ServeHTTP(w, r)
    })
}
```

**Как сочетаем:**
- Каждый интерфейс имеет одну ответственность
- `HTTPHandler` - максимально обобщенный (стандартный интерфейс из net/http)
- Декораторы можно комбинировать по необходимости

## Пример 3: Разделение интерфейса кэша

**Было:**
```go
type Cache interface {
    // Основные операции
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    
    // Статистика и метрики
    Hits() int64
    Misses() int64
    HitRate() float64
    
    // Управление памятью
    EvictOldest()
    Clear()
    Size() int64
    
    // Распределенный кэш
    GetFromReplica(key string) (interface{}, error)
    SyncWithPrimary() error
}
```

**Стало:**
```go
// Базовый обобщенный интерфейс кэша
type Cache[K comparable, V any] interface {
    Get(key K) (V, bool)
    Set(key K, value V, ttl time.Duration)
    Delete(key K)
}

// Интерфейс для метрик (отдельная ответственность)
type CacheMetrics interface {
    Hits() int64
    Misses() int64
    HitRate() float64
}

// Интерфейс для управления памятью
type MemoryManager interface {
    EvictOldest()
    Clear()
    Size() int64
}

// Интерфейс для распределенного кэша
type DistributedCache[K comparable, V any] interface {
    Cache[K, V]
    GetFromReplica(key K) (V, bool)
    SyncWithPrimary() error
}

// Реализация может реализовывать несколько интерфейсов
type LRUCache[K comparable, V any] struct {
    // ...
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) { /* ... */ }
func (c *LRUCache[K, V]) Set(key K, value V, ttl time.Duration) { /* ... */ }
func (c *LRUCache[K, V]) EvictOldest() { /* ... */ }
```

**Как сочетаем:**
- `Cache` - обобщенный с дженериками, минимальный интерфейс
- Каждая дополнительная функциональность в отдельном интерфейсе
- Клиенты могут требовать только нужные им интерфейсы

## Пример 4: Разделение интерфейса нотификаций

**Было:**
```go
type Notifier interface {
    // Разные типы нотификаций
    SendEmail(to, subject, body string) error
    SendSMS(phone, message string) error
    SendPush(userID, title, message string) error
    
    // Шаблонизация
    RenderTemplate(templateName string, data interface{}) (string, error)
    
    // Логирование
    LogNotification(notificationType, recipient, message string)
    
    // Ретри
    WithRetry(attempts int, delay time.Duration)
}
```

**Стало:**
```go
// Базовый обобщенный интерфейс отправки
type Sender[T any] interface {
    Send(recipient string, message T) error
}

// Конкретные типы сообщений
type Email struct {
    Subject string
    Body    string
}

type SMS struct {
    Text string
}

type PushNotification struct {
    Title   string
    Message string
    Badge   int
}

// Интерфейсы для конкретных каналов
type EmailSender interface {
    Sender[Email]
}

type SMSSender interface {
    Sender[SMS]
}

type PushSender interface {
    Sender[PushNotification]
}

// Дополнительные интерфейсы
type TemplateRenderer interface {
    Render(templateName string, data interface{}) (string, error)
}

type NotificationLogger interface {
    Log(notificationType, recipient string, metadata map[string]interface{})
}

type Retryable interface {
    WithRetry(attempts int, delay time.Duration)
}

// Композитный нотификатор
type CompositeNotifier struct {
    emailSender EmailSender
    smsSender   SMSSender
    // ...
}
```

**Как сочетаем:**
- `Sender[T]` - обобщенный интерфейс с дженериком
- Каждый канал нотификации имеет свой интерфейс
- Дополнительные возможности в отдельных интерфейсах