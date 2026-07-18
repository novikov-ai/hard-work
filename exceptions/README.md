# Разбираемся с исключениями

В Go как таковых исключений нет, есть ошибки и паники.

## 1. `panic` как guard clause в конструкторе

**Слабо.** Проверка есть, но она ничего не гарантирует статически, и падение происходит в рантайме далеко от места ошибки.

```go
func NewClient(baseURL string, timeout time.Duration) *Client {
	if baseURL == "" {
		panic("baseURL is required")
	}
	if timeout <= 0 {
		panic("timeout must be positive")
	}
	return &Client{baseURL: baseURL, timeout: timeout}
}
```

**Правильно.** Обязательное выносим в параметры, необязательное — в опции со значением по умолчанию. Тогда забыть `baseURL` нельзя, а невалидный `timeout` становится ошибкой значения, а не паникой.

```go
type Option func(*Client)

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { if d > 0 { c.timeout = d } }
}

// baseURL обязателен — это видно из сигнатуры, IDE подскажет
func NewClient(baseURL URL, opts ...Option) *Client {
	c := &Client{baseURL: baseURL, timeout: 30 * time.Second}
	for _, o := range opts { o(c) }
	return c
}

type URL struct{ raw string }

func ParseURL(s string) (URL, error) {
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return URL{}, fmt.Errorf("invalid base URL %q", s)
	}
	return URL{raw: u.String()}, nil
}

func (u URL) String() string { return u.raw }
```

---

## 2. `panic`/`recover` как управление потоком

**Слабо.** Классический трюк из рекурсивных парсеров: паника вместо проброса ошибки через десять уровней.

```go
func (p *parser) expect(tok token) {
	if p.next() != tok {
		panic(parseError{pos: p.pos, want: tok})
	}
}

func Parse(src []byte) (ast *Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(parseError)
		}
	}()
	return newParser(src).parseFile(), nil
}
```

Проблема в `r.(parseError)` — непроверенное приведение проглотит чужую панику, например `nil pointer dereference`, и превратит её в невнятную ошибку парсинга.

**Правильно.** Если оставляем `recover` — только с проверкой типа и с обязательным пробросом чужой паники дальше:

```go
defer func() {
	r := recover()
	if r == nil {
		return
	}
	pe, ok := r.(parseError)
	if !ok {
		panic(r) // чужая паника — не наше дело, пробрасываем
	}
	err = pe
}()
```

---

## 3. Ошибка, которая не может произойти

**Слабо.** Сигнатура возвращает `error`, но вызвать его невозможно — и все вызывающие пишут либо `_`, либо бессмысленную проверку.

```go
func (c Currency) Symbol() (string, error) {
	switch c {
	case USD: return "$", nil
	case EUR: return "€", nil
	}
	return "", fmt.Errorf("unknown currency %v", c)
}

Здесь `error` — это налог, который платят все вызывающие из-за того, что `Currency` объявлена как `type Currency int` и технически может принять `Currency(999)`.

**Правильно.** Сделать невалидное состояние непредставимым, а разбор — единственной точкой проверки.

```go
type Currency struct{ code string; symbol string }

var (
	USD = Currency{"USD", "$"}
	EUR = Currency{"EUR", "€"}
)

// Единственный вход извне
func ParseCurrency(code string) (Currency, error) {
	switch code {
	case "USD": return USD, nil
	case "EUR": return EUR, nil
	}
	return Currency{}, fmt.Errorf("unknown currency %q", code)
}

// Ошибки больше нет — она невозможна
func (c Currency) Symbol() string { return c.symbol }
```

---

## 4. Ошибка как строка

**Слабо.** Сравнение по тексту — хрупкая связь между пакетами: изменение формулировки в логах ломает логику ретраев.

```go
if err != nil {
	if strings.Contains(err.Error(), "not found") {
		return defaultConfig, nil
	}
	if err.Error() == "connection refused" {
		return nil, retry(ctx)
	}
	return nil, err
}
```

**Правильно.** Классифицировать ошибки типами — это и есть «явная классификация потенциальных ошибок» из статьи. Причём для «есть/нет» достаточно sentinel-значения, а для ошибок с данными нужен тип.

```go
var ErrNotFound = errors.New("config: not found")

type TransientError struct {
	Op  string
	Err error
}

func (e *TransientError) Error() string { return e.Op + ": " + e.Err.Error() }
func (e *TransientError) Unwrap() error { return e.Err }

// на вызывающей стороне:
switch {
case errors.Is(err, ErrNotFound):
	return defaultConfig, nil
case new(TransientError) != nil && errors.As(err, &te):
	return nil, retry(ctx, te)
default:
	return nil, err
}
```

---

## 5. Проброс без контекста

**Слабо.** Самый частый паттерн в Go-коде вообще. Формально ошибка обработана, фактически — «переложили вину на других».

```go
func (s *Store) UserOrders(id UserID) ([]Order, error) {
	u, err := s.loadUser(id)
	if err != nil {
		return nil, err
	}
	o, err := s.loadOrders(u.AccountID)
	if err != nil {
		return nil, err
	}
	return o, nil
}
```

На верхнем уровне вы увидите `sql: no rows in result set` — и не узнаете ни какой запрос упал, ни для какого пользователя.

**Правильно.** Оборачивать с `%w`: цепочка сохраняется для `errors.Is`, а сообщение накапливает контекст.

```go
func (s *Store) UserOrders(id UserID) ([]Order, error) {
	u, err := s.loadUser(id)
	if err != nil {
		return nil, fmt.Errorf("load user %d: %w", id, err)
	}
	o, err := s.loadOrders(u.AccountID)
	if err != nil {
		return nil, fmt.Errorf("load orders for account %d: %w", u.AccountID, err)
	}
	return o, nil
}
```