# OCP с т.зр. ФП

## 1. Middleware в веб-фреймворке (Композиция функций)

Роутер/обработчик закрыт для модификации, но открыт для расширения за счёт цепочки middleware.

**Было (без OCP):**
```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Логика обработки
}

// Где-то в роутинге
http.HandleFunc("/path", MyHandler) // Невозможно добавить логирование или аутентификацию без изменения кода
```

**Стало (OCP через композицию):**
```go
// Тип-алиас для функционального типа — основа "чёрного ящика"
type Middleware func(http.Handler) http.Handler

// Middleware для логирования (расширение без изменения существующих обработчиков)
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Middleware для аутентификации (ещё одно расширение)
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Композиция middleware (создание нового поведения из существующих функций)
router.Handle("/api/secure",
    Logging(Auth(http.HandlerFunc(MyHandler)))) // <- Не изменяя MyHandler, мы добавили функциональность
```

## 2. Конфигурируемая валидация данных (Функции высшего порядка как стратегии)
Вместо жёсткой иерархии валидаторов используем функции, которые можно комбинировать:

```go
// Тип-валидатор: функция, принимающая данные и возвращающая ошибку (или nil)
type Validator[T any] func(T) error

// Базовые валидаторы (закрыты для модификации, готовы к использованию)
func NotEmpty(s string) error {
    if s == "" {
        return errors.New("must not be empty")
    }
    return nil
}

func MinLength(n int) Validator[string] {
    return func(s string) error {
        if len(s) < n {
            return fmt.Errorf("length must be at least %d", n)
        }
        return nil
    }
}

func GreaterThan(threshold int) Validator[int] {
    return func(v int) error {
        if v <= threshold {
            return fmt.Errorf("must be greater than %d", threshold)
        }
        return nil
    }
}

// Композитор валидаторов (открыт для расширения новыми стратегиями)
func All[T any](validators ...Validator[T]) Validator[T] {
    return func(v T) error {
        for _, validator := range validators {
            if err := validator(v); err != nil {
                return err
            }
        }
        return nil
    }
}

// Использование (расширяем систему валидации, не меняя существующий код)
validateUsername := All(
    NotEmpty,
    MinLength(3),
    // Легко добавить новый валидатор без изменения All, NotEmpty и т.д.
    func(s string) error { // Анонимная функция как валидатор
        if strings.Contains(s, "admin") {
            return errors.New("username cannot contain 'admin'")
        }
        return nil
    },
)

err := validateUsername("johndoe") // Валидация через композицию
```

## 3. Обработка элементов коллекции (Map/Reduce/Filter как "чёрные ящики")

Функции высшего порядка для работы с коллекциями — проявление OCP в ФП:

```go
// SliceMap - абстрактная операция преобразования (закрыта для модификации)
func SliceMap[T, U any](in []T, mapper func(T) U) []U {
    out := make([]U, len(in))
    for i, v := range in {
        out[i] = mapper(v) // Поведение определяется переданной функцией
    }
    return out
}

// SliceFilter - абстрактная операция фильтрации (закрыта для модификации)
func SliceFilter[T any](in []T, predicate func(T) bool) []T {
    var out []T
    for _, v := range in {
        if predicate(v) {
            out = append(out, v)
        }
    }
    return out
}

// Использование в бизнес-логике (расширяем, передавая разные функции)
users := []User{...}

// Получить emails активных пользователей
activeEmails := SliceMap(
    SliceFilter(users, func(u User) bool { return u.IsActive }),
    func(u User) string { return u.Email },
)
// Завтра можем добавить другую логику фильтрации или маппинга,
// не изменяя SliceMap и SliceFilter.
```

## 4. Парсинг конфигурации с использованием алгебраических типов данных (имитация через интерфейсы)

В Go нет ADT, но их можно имитировать, что соответствует OCP:

```go
// Абстрактный тип конфига (закрыт для модификации кода, который его потребляет)
type Config interface {
    // Маркер-метод, чтобы сделать интерфейс закрытым (необязательно, но для безопасности)
    isConfig()
}

// Варианты конфига (открыты для добавления новых типов)
type FileConfig struct{
    Path string
}

func (FileConfig) isConfig() {}

type EnvConfig struct{
    Prefix string
}

func (EnvConfig) isConfig() {}

type RemoteConfig struct{
    URL string
}

func (RemoteConfig) isConfig() {}

// Парсер (закрыт для модификации, но открыт для расширения через новые типы Config)
func ParseConfig(cfg Config) (*Settings, error) {
    switch v := cfg.(type) { // Тип определяет поведение
    case FileConfig:
        return parseFromFile(v.Path) // Эти функции можно развивать независимо
    case EnvConfig:
        return parseFromEnv(v.Prefix)
    case RemoteConfig:
        return parseFromRemote(v.URL)
    default:
        // Компилятор НЕ предупредит о новом типе, но можно добавить панику
        panic("unhandled config type") // или вернуть ошибку
    }
}

// Завтра добавляем DatabaseConfig, не меняя сигнатуру ParseConfig
// (хотя и придётся добавить case - это слабое место)
```
