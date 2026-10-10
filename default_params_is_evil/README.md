# Почему дефлотные параметры -- зло

## Go

> **Важно:** в Go нет механики дефолтных параметров. Нельзя написать `func f(x int = 1)`.
> Все примеры ниже — это эмуляция дефолтов средствами языка.

### Плохо: нулевое значение как «не задано»

```go
type Config struct {
    Port    int
    Retries int
}

func NewServer(cfg Config) *Server {
    if cfg.Port == 0 {
        cfg.Port = 8080
    }
    if cfg.Retries == 0 {
        cfg.Retries = 3 // дефолт
    }
    return &Server{cfg}
}

// Хочу отключить ретраи, передаю 0 — а получаю 3.
srv := NewServer(Config{Retries: 0})
```

Проблема: нельзя отличить «не задано» от «явно задан ноль». Дефолт молча
перезаписывает осознанное решение вызывающего, а поведение неочевидно из сигнатуры.

### Хорошо: functional options

```go
type Server struct {
    port    int
    retries int
}

type Option func(*Server)

func WithPort(p int) Option    { return func(s *Server) { s.port = p } }
func WithRetries(n int) Option { return func(s *Server) { s.retries = n } }

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080, retries: 3} // дефолты в одном месте
    for _, o := range opts {
        o(s)
    }
    return s
}

srv := NewServer(WithRetries(0)) // ретраи отключены явно, порт — 8080
```

Почему лучше: дефолты собраны в одном месте, каждое переопределение явное и
читается на месте вызова, `0` — полноценное значение, новые опции добавляются без
поломки API.
