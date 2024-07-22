# Интерфейс компактнее реализации?

### Призрачное состояние

1. Локальная переменная `cache`:

~~~go
type Calculator struct {
    cache map[int]int
}

func (c *Calculator) Compute(value int) int {
    if result, found := c.cache[value]; found {
        return result
    }
    result := heavyComputation(value)
    c.cache[value] = result
    return result
}

func heavyComputation(value int) int {
    // сложные вычисления
    return value * value
}
~~~

2. Флаг состояния isUpdated:

~~~go
type DataProcessor struct {
    data     []int
    isUpdated bool
}

func (p *DataProcessor) Process() {
    if !p.isUpdated {
        p.data = preprocess(p.data)
        p.isUpdated = true
    }
    p.data = performProcessing(p.data)
}
~~~

### Погрешности/неточности

1. Функция жестко задает фиксированную длину строки в 10 символов, чем ограничивает её гибкость:

~~~go
// package stringutil

// PadString pads the input string with spaces up to a fixed length of 10 characters.
func PadString(input string) string {
    if len(input) >= 10 {
        return input[:10]
    }
    return input + strings.Repeat(" ", 10-len(input))
}
~~~

~~~go
// PadString pads the input string with spaces up to the specified length.
// Spec: The length should be specified as an argument to the function, allowing more flexibility.
func PadString(input string, length int) string {
    if len(input) >= length {
        return input[:length]
    }
    return input + strings.Repeat(" ", length-len(input))
}
~~~

2. Функция использует жестко заданный тайм-аут, что может быть неподходящим для различных сценариев:

~~~go
// package network

// ConnectToServer attempts to connect to a server with a fixed timeout of 5 seconds.
func ConnectToServer(address string) (Connection, error) {
    timeout := 5 * time.Second
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return nil, err
    }
    return conn, nil
}
~~~

~~~go
// ConnectToServer attempts to connect to a server with a specified timeout.
// Spec: The timeout should be specified as an argument to the function, allowing more flexibility.
func ConnectToServer(address string, timeout time.Duration) (Connection, error) {
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return nil, err
    }
    return conn, nil
}
~~~

3. Функция использует жестко заданный порог для определения значительных изменений:

~~~go
// package analytics

// IsSignificantChange checks if the change in value is significant, using a fixed threshold of 0.1.
func IsSignificantChange(oldValue, newValue float64) bool {
    return math.Abs(newValue-oldValue) > 0.1
}
~~~

~~~go
// IsSignificantChange checks if the change in value is significant, using a specified threshold.
// Spec: The threshold should be specified as an argument to the function, allowing more flexibility.
func IsSignificantChange(oldValue, newValue, threshold float64) bool {
    return math.Abs(newValue-oldValue) > threshold
}
~~~

### Интерфейс явно не должен быть проще реализации

1. В GORM, популярной ORM-библиотеке для Go, интерфейс для управления транзакциями кажется простым, но его реализация может быть очень сложной:

~~~go
// Transactional function
func (db *DB) Transaction(fc func(tx *DB) error) (err error) {
    // Implementation details...
}
~~~

Использование подходящих типов данных сделает интерфейс очевиднее.

~~~go
type TransactionFunc func(tx *DB) error

type TransactionResult struct {
    Success bool
    Error   error
}

// Transactional function with explicit result type
func (db *DB) Transaction(fc TransactionFunc) TransactionResult {
    // Implementation details...
}
~~~

Теперь возвращаемый тип TransactionResult явно указывает на успех или неудачу транзакции, что улучшает читаемость и понимание кода.

2. В стандартной библиотеке Go пакет net/http предоставляет интерфейс для обработки HTTP-запросов.

~~~go
type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}
~~~

Реализация обработки HTTP-запросов включает множество аспектов, таких как: маршрутизация, обработка заголовков и тела запроса, управление состоянием сессий и т.д.

Использование более конкретных типов данных и структур может помочь сделать интерфейс более наглядным.

~~~go
type HTTPMethod string

const (
    GET    HTTPMethod = "GET"
    POST   HTTPMethod = "POST"
    PUT    HTTPMethod = "PUT"
    DELETE HTTPMethod = "DELETE"
)

type HTTPRequest struct {
    Method  HTTPMethod
    Headers map[string]string
    Body    []byte
}

type HTTPResponse struct {
    StatusCode int
    Headers    map[string]string
    Body       []byte
}

type Handler interface {
    ServeHTTP(req HTTPRequest) HTTPResponse
}
~~~

3. Logrus — популярная библиотека логирования для Go, но интерфейс выглядит сложно читаемым.

~~~go
type Logger interface {
    Info(args ...interface{})
    Error(args ...interface{})
    // Другие методы
}
~~~

Подходящие типы данных помогают "разгрузить" интерфейс и сделать его более очевидным для использования.

~~~go
type LogLevel string

const (
    InfoLevel  LogLevel = "INFO"
    ErrorLevel LogLevel = "ERROR"
    // Другие уровни
)

type LogMessage struct {
    Level   LogLevel
    Message string
    Fields  map[string]interface{}
}

type Logger interface {
    Log(message LogMessage)
}
~~~