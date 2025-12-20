# LSP с т.зр. ФП

### 1. **Работа с `io.Reader` и `io.Writer`**
**Контекст:** Любая функция, принимающая `io.Reader` или `io.Writer`.
```go
func ProcessData(r io.Reader) error {
    data, err := io.ReadAll(r)
    // ... обработка
}
```
**Почему LSP:** Функция `ProcessData` работает с абстракцией чтения (`io.Reader`). Мы можем подставить:
- `*os.File` (файл)
- `*bytes.Buffer` (буфер в памяти)
- `*strings.Reader` (строка)
- `net.Conn` (сетевое соединение)
- Кастомную структуру с методом `Read()`

Все эти типы взаимозаменяемы, потому что каждый соблюдает контракт `Read(p []byte) (n int, err error)`. Это чистая демонстрация LSP: код, ожидающий `io.Reader`, корректно работает с любым его подтипом.

### 2. **HTTP-обработчики (`http.Handler`)**
**Контекст:** Любой HTTP-сервер.
```go
http.Handle("/path", myHandler)
```
**Почему LSP:** Интерфейс `http.Handler` требует одного метода:
```go
ServeHTTP(ResponseWriter, *Request)
```
Можно подставлять:
- Стандартные хендлеры (`http.FileServer`, `http.RedirectHandler`)
- Структуры с middleware (цепочка обработки)
- Любые кастомные реализации

Даже `http.HandlerFunc` (адаптер для функций) — это яркий пример соблюдения LSP: тип-функция подставляется вместо интерфейса, сохраняя контракт.

### 3. **Сортировка через `sort.Interface`**
**Контекст:** Использование `sort.Sort()`.
```go
type Users []User
func (u Users) Len() int           { return len(u) }
func (u Users) Less(i, j int) bool { return u[i].Age < u[j].Age }
func (u Users) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

sort.Sort(Users(users)) // Работает с любой коллекцией
```
**Почему LSP:** Пакет `sort` требует только три метода. Неважно, сортируем мы пользователей, заказы или товары — алгоритм сортировки работает корректно с любой реализацией `sort.Interface`. Это гарантирует, что подстановка любого совместимого типа не сломает логику сортировки.

### 4. **Кастомные хранилища (Repository Pattern)**
**Контекст:** Слоистая архитектура приложения.
```go
type UserRepository interface {
    GetByID(id string) (*User, error)
    Save(user *User) error
}

// В production используем PostgreSQL
type PostgresUserRepo struct{ /* ... */ }

// В тестах используем mock
type MockUserRepo struct{ /* ... */ }

// В use-case
type UserService struct {
    repo UserRepository // Зависим от абстракции
}
```
**Почему LSP:** `UserService` зависит от абстракции `UserRepository`. Мы можем подставить:
- Реальную реализацию для БД
- In-memory реализацию для интеграционных тестов
- Mock для unit-тестов
- Stub для демо-режима

Каждая реализация соблюдает контракт, поэтому замена одной на другую не ломает `UserService`. Это прямое применение LSP для достижения слабой связанности.