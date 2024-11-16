# Повышаем полиморфность кода

### Пример 1

Паттерн "наблюдатель" работал только с определенными функциями о чем также свидетельствовал нейминг. Сделал поведение и нейминг полиморфным, чтобы была возможность работать не только с типом `string`.

Было:
~~~go
type StringObserver func(string)

type StringNotifier struct {
	observers []StringObserver
}

func (n *StringNotifier) Register(observer StringObserver) {
	n.observers = append(n.observers, observer)
}

func (n *StringNotifier) Notify(data string) {
	for _, observer := range n.observers {
		observer(data)
	}
}
~~~

Стало:
~~~go
type Observer[T any] func(T)

type Notifier[T any] struct {
	observers []Observer[T]
}

func (n *Notifier[T]) Register(observer Observer[T]) {
	n.observers = append(n.observers, observer)
}

func (n *Notifier[T]) Notify(data T) {
	for _, observer := range n.observers {
		observer(data)
	}
}
~~~

### Пример 2

Был определен кэш, работающий только с `int`. Используя дженерики добавил полиморфизма и переименовал функцию.

Было:
~~~go
type IntCache struct {
	data map[string]int
}

func (c *IntCache) Set(key string, value int) {
	c.data[key] = value
}

func (c *IntCache) Get(key string) (int, bool) {
	value, exists := c.data[key]
	return value, exists
}
~~~

Стало:
~~~go
type Cache[T any] struct {
	data map[string]T
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make(map[string]T)}
}

func (c *Cache[T]) Set(key string, value T) {
	c.data[key] = value
}

func (c *Cache[T]) Get(key string) (T, bool) {
	value, exists := c.data[key]
	return value, exists
}
~~~

### Пример 3

`ExecuteUserCommand` была завязана на работу со структурой `UserCommand`. Сделал поведение более общим и позволил работать с любым типом команд. 

Было:
~~~go
type UserCommand struct {
	UserID int
	Action string
}

func ExecuteUserCommand(cmd UserCommand) {
	fmt.Printf("Executing %s for user ID %d\n", cmd.Action, cmd.UserID)
}
~~~

Стало:
~~~go
type Command[T any] struct {
	Target T
	Action string
}

func ExecuteCommand[T any](cmd Command[T], executor func(T, string)) {
	executor(cmd.Target, cmd.Action)
}
~~~

### Пример 4

`HandleUserEvent` был завязан на работу с событиями пользователя. Сделал общий, полиморфный обработчик для любых событий.

Было:
~~~go
type User struct {
	ID   int
	Name string
}

type UserEvent struct {
	User  User
	Event string
}

func HandleUserEvent(event UserEvent) {
	fmt.Printf("Handling event %s for user %s\n", event.Event, event.User.Name)
}
~~~

Стало:
~~~go
type Event[T any] struct {
	Entity T
	Event  string
}

func HandleEvent[T any](event Event[T], handler func(T, string)) {
	handler(event.Entity, event.Event)
}
~~~

### Пример 5

Хранилище предполагало работу с определенным типом данных. Добавил интерфейс над хранилищем и поменял реализацию для методов, чтобы была возможность работать с любыми типами.

Было:
~~~go
type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{data: make(map[string]string)}
}

func (s *InMemoryStorage) Save(key, value string) {
	s.data[key] = value
}

func (s *InMemoryStorage) Load(key string) (string, bool) {
	value, exists := s.data[key]
	return value, exists
}

func (s *InMemoryStorage) Delete(key string) {
	delete(s.data, key)
}
~~~

Стало:
~~~go
type StorageBackend[T any] interface {
	Save(key string, value T)
	Load(key string) (T, bool)
	Delete(key string)
}

type InMemoryBackend[T any] struct {
	data map[string]T
}

func NewInMemoryBackend[T any]() *InMemoryBackend[T] {
	return &InMemoryBackend[T]{data: make(map[string]T)}
}

func (b *InMemoryBackend[T]) Save(key string, value T) { b.data[key] = value }
func (b *InMemoryBackend[T]) Load(key string) (T, bool) {
	value, exists := b.data[key]
	return value, exists
}
func (b *InMemoryBackend[T]) Delete(key string) { delete(b.data, key) }

type Storage[T any] struct {
	backend StorageBackend[T]
}

func NewStorage[T any](backend StorageBackend[T]) *Storage[T] {
	return &Storage[T]{backend: backend}
}

func (s *Storage[T]) Save(key string, value T) { s.backend.Save(key, value) }
func (s *Storage[T]) Load(key string) (T, bool) { return s.backend.Load(key) }
func (s *Storage[T]) Delete(key string) { s.backend.Delete(key) }
~~~