# DIP с т.зр. ФП

## Пример 1

Для создания хэндлера используем DIP, чтобы была возможность не зависеть от конкретной реализации Composer:

~~~go
// handler.go
func New(composer Composer, metric internal.Metric, log internal.Log) *Handler {
	return &Handler{
		composer: composer,
		metric:   metric,
		log:      log,
	}
}
~~~

~~~go
// contract.go
type Composer interface {
	Compose(ids []string) ([]model.DtoV2, error)
}
~~~


## Пример 2

Для создания механизма, который рассылает уведомления, используем интерфейс контекста (Context):

~~~go
// sender.go
func (s sender) SubscriptionError(ctx context.Context, uid int64, errText, errDescription *string) {
	event := clickstream_events.NewErrorShowV2(
        ctx, // <-- передаем контекст, удовлетворяющий интерфейсу context.Context
		clickstream_events.ErrorShowV2{},
	)

	event.SetUid(&uid)

	if errDescription != nil {
		event.SetErrorDescription(errDescription)
		errText = &clientErrorText
	}

	event.SetErrorText(errText)

	// ...
}
~~~

## Пример 3

Для работы с внешним API используем DIP, чтобы не зависеть от конкретной реализации клиента:

```go
// service.go
func NewUserService(client APIClient, cache Cache) *UserService {
	return &UserService{
		client: client,
		cache:  cache,
	}
}
```

```go
// contract.go
type APIClient interface {
	GetUser(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}
```

## Пример 4

Для работы с очередями используем интерфейс Producer:

```go
// processor.go
func (p *Processor) HandleEvent(ctx context.Context, event Event) error {
	msg := Message{
		Type:    "user_event",
		Payload: event,
	}
	
	return p.producer.Send(ctx, "events", msg) // <-- producer реализует интерфейс
}
```

```go
// contract.go
type Producer interface {
	Send(ctx context.Context, topic string, msg interface{}) error
	Close() error
}
```

## Пример 5

Для работы с репозиторием используем DIP, чтобы можно было менять базу данных:

```go
// handler.go
func NewOrderHandler(repo OrderRepository, validator Validator) *OrderHandler {
	return &OrderHandler{
		repo:      repo,
		validator: validator,
	}
}
```

```go
// contract.go
type OrderRepository interface {
	FindByID(ctx context.Context, id int64) (*Order, error)
	Save(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id int64) error
}
```