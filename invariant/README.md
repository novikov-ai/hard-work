# Инварианты и качественный код (1)

## Заготовка для ассертов

~~~go
package assert
 
import "fmt"
 
func True(cond bool, msg string, args ...any) {
	if !cond {
		panic("assert: " + fmt.Sprintf(msg, args...))
	}
}
 
func NotNil(v any, name string) {
	if v == nil {
		panic("assert: " + name + " must not be nil")
	}
}
~~~

## Пример 1

Без инварианта (плохо):
~~~go
type User struct {
    ID    string
    Email string
    Name  string
}

func CreateUser(email, name string) (*User, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }
    // ... и так в каждом месте, где создаём User
    return &User{ID: uuid.New().String(), Email: email, Name: name}, nil
}

func SendNotification(u *User) error {
    if u.Email == "" {
        // Защитная проверка "на всякий случай"
        return errors.New("user has no email")
    }
    // ...
}
~~~

Вводим инвариант NonEmptyString:
~~~go
package types
 
import (
	"errors"
	"strings"
 
	"internal/assert"
)
 
var ErrEmptyString = errors.New("string is empty")
 
type NonEmptyString struct {
	v string
}
 
func NewNonEmptyString(s string) (NonEmptyString, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return NonEmptyString{}, ErrEmptyString
	}
	return NonEmptyString{v: s}, nil
}
 
// MustNonEmptyString — для констант и тестов, где пустая строка означает баг.
func MustNonEmptyString(s string) NonEmptyString {
	n, err := NewNonEmptyString(s)
	if err != nil {
		panic(err)
	}
	return n
}
 
func (n NonEmptyString) String() string {
	// Ловит zero-value: NonEmptyString{} в обход конструктора.
	assert.True(n.v != "", "NonEmptyString is empty (zero-value used?)")
	return n.v
}
 
func (n NonEmptyString) Equal(other NonEmptyString) bool {
	return n.v == other.v
}
~~~

~~~go
type User struct {
    ID    string
    Email types.NonEmptyString  // компилятор гарантирует непустоту
    Name  types.NonEmptyString
}

func SendNotification(u *User) error {
    // Email гарантированно не пуст без проверок
    return mailer.Send(u.Email.String(), ...)
}
~~~

## Пример 2

Без инварианта (плохо):
~~~go
type Order struct {
    Amount   float64  // float для денег — преступление
    Currency string   // просто строка, можно передать "USDD"
}

func Transfer(from, to *Account, amount float64, currency string) error {
    if amount < 0 {
        return errors.New("amount must be positive")
    }
    if from.Currency != currency {
        return errors.New("currency mismatch")
    }
    // ... и так далее
}
~~~

Добавляем инвариант, гарантируя корректность операций:
~~~go
package money
 
import (
	"errors"
	"fmt"
	"math"
 
	"internal/assert"
)
 
type Currency string
 
const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	RUB Currency = "RUB"
)
 
func (c Currency) valid() bool {
	switch c {
	case USD, EUR, RUB:
		return true
	}
	return false
}
 
var (
	ErrNegativeAmount  = errors.New("amount must be non-negative")
	ErrUnknownCurrency = errors.New("unknown currency")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)
 
// Money хранит сумму в минимальных единицах (копейки, центы).
// Никаких float — намеренно.
type Money struct {
	minor    int64
	currency Currency
}
 
func New(minor int64, c Currency) (Money, error) {
	if minor < 0 {
		return Money{}, ErrNegativeAmount
	}
	if !c.valid() {
		return Money{}, fmt.Errorf("%w: %q", ErrUnknownCurrency, c)
	}
	return Money{minor: minor, currency: c}, nil
}
 
func Zero(c Currency) Money {
	m, err := New(0, c)
	if err != nil {
		panic(err) // невалидная валюта в коде — это баг
	}
	return m
}
 
func (m Money) Minor() int64       { return m.minor }
func (m Money) Currency() Currency { return m.currency }
 
func (m Money) Add(other Money) (Money, error) {
	assert.True(m.minor >= 0, "receiver invariant broken: minor=%d", m.minor)
	assert.True(other.minor >= 0, "operand invariant broken: minor=%d", other.minor)
 
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("%w: %s vs %s", ErrCurrencyMismatch, m.currency, other.currency)
	}
	if m.minor > math.MaxInt64-other.minor {
		return Money{}, fmt.Errorf("overflow: %d + %d", m.minor, other.minor)
	}
 
	return Money{minor: m.minor + other.minor, currency: m.currency}, nil
}
 
func (m Money) Sub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("%w: %s vs %s", ErrCurrencyMismatch, m.currency, other.currency)
	}
	if other.minor > m.minor {
		return Money{}, ErrNegativeAmount
	}
	return Money{minor: m.minor - other.minor, currency: m.currency}, nil
}
 
func (m Money) String() string {
	return fmt.Sprintf("%d.%02d %s", m.minor/100, m.minor%100, m.currency)
}
~~~

## Пример 3

Без инварианта (плохо):
~~~go
type HTTPClient struct {
    BaseURL    string
    Timeout    time.Duration
    APIKey     string
    httpClient *http.Client
}

func (c *HTTPClient) Get(path string) (*Response, error) {
    if c.BaseURL == "" {
        return nil, errors.New("BaseURL is required")
    }
    if c.Timeout == 0 {
        return nil, errors.New("Timeout is required")
    }
    if c.APIKey == "" {
        return nil, errors.New("APIKey is required")
    }
    if c.httpClient == nil {
        return nil, errors.New("client not initialized — call Init() first")
    }
    // ... запрос
}
~~~

Добавляем инвариант:
~~~go

type Config struct {
	BaseURL string
	Timeout time.Duration
	APIKey  string
}
 
type Client struct {
	baseURL *url.URL
	apiKey  string
	hc      *http.Client
}
 
func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("BaseURL is required")
	}
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse BaseURL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("BaseURL scheme: want http/https, got %q", u.Scheme)
	}
	if u.Host == "" {
		return nil, errors.New("BaseURL host is empty")
	}
	if cfg.Timeout <= 0 {
		return nil, errors.New("Timeout must be positive")
	}
	if cfg.APIKey == "" {
		return nil, errors.New("APIKey is required")
	}
 
	return &Client{
		baseURL: u,
		apiKey:  cfg.APIKey,
		hc:      &http.Client{Timeout: cfg.Timeout},
	}, nil
}
 
func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	assert.NotNil(c, "client")
	assert.NotNil(c.baseURL, "baseURL")
	assert.NotNil(c.hc, "http client")
	assert.True(c.apiKey != "", "apiKey is empty")
 
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL.JoinPath(path).String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
 
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
 
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
~~~