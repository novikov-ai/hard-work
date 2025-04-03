# Избавляемся от условных инструкций

## 1. Управление состоянием подключения (State Pattern)

**Было (с использованием if для проверки состояния подключения):**

```go
package main

import "fmt"

type Connection struct {
	state string // "open", "closed" или "error"
}

func handleConnection(conn *Connection) {
	if conn.state == "open" {
		fmt.Println("Обрабатываем активное соединение")
	} else if conn.state == "closed" {
		fmt.Println("Соединение закрыто, закрываем ресурсы")
	} else if conn.state == "error" {
		fmt.Println("Ошибка соединения, выполняем восстановление")
	} else {
		fmt.Println("Неизвестное состояние")
	}
}

func main() {
	conn := &Connection{state: "open"}
	handleConnection(conn)
}
```

**Стало (использование паттерна «Состояние» через полиморфизм):**

```go
package main

import "fmt"

// Определяем интерфейс состояния
type ConnectionState interface {
	Handle(*Connection)
}

type Connection struct {
	state ConnectionState
}

// Конкретные реализации состояний:
type OpenState struct{}
func (s OpenState) Handle(conn *Connection) {
	fmt.Println("Обрабатываем активное соединение")
}

type ClosedState struct{}
func (s ClosedState) Handle(conn *Connection) {
	fmt.Println("Соединение закрыто, закрываем ресурсы")
}

type ErrorState struct{}
func (s ErrorState) Handle(conn *Connection) {
	fmt.Println("Ошибка соединения, выполняем восстановление")
}

func main() {
	// Инициализация соединения уже с корректным состоянием
	conn := &Connection{state: OpenState{}}
	conn.state.Handle(conn)
}
```

**Комментарий:**  
Здесь состояние подключения инкапсулировано в объекте, который сам знает, как обработать свою логику. Таким образом, основная логика не содержит условных операторов для ветвления.

---

## 2. Обработка платежей (Полиморфизм)

**Было (ветвление по типу платежа):**

```go
package main

import "fmt"

func ProcessPayment(paymentType string, amount float64) {
	if paymentType == "credit" {
		fmt.Printf("Обрабатываем кредитную карту: %.2f\n", amount)
	} else if paymentType == "paypal" {
		fmt.Printf("Обрабатываем PayPal: %.2f\n", amount)
	} else {
		fmt.Println("Неизвестный тип платежа")
	}
}

func main() {
	ProcessPayment("credit", 100.0)
}
```

**Стало (использование интерфейса PaymentMethod и фабрики):**

```go
package main

import "fmt"

// Интерфейс для метода оплаты
type PaymentMethod interface {
	Process(amount float64)
}

type CreditCard struct{}
func (c CreditCard) Process(amount float64) {
	fmt.Printf("Обрабатываем кредитную карту: %.2f\n", amount)
}

type PayPal struct{}
func (p PayPal) Process(amount float64) {
	fmt.Printf("Обрабатываем PayPal: %.2f\n", amount)
}

// Фабрика возвращает реализацию PaymentMethod
func NewPaymentMethod(method string) PaymentMethod {
	switch method {
	case "credit":
		return CreditCard{}
	case "paypal":
		return PayPal{}
	default:
		// Для неизвестного типа можно вернуть объект-заглушку
		return nil
	}
}

func main() {
	pm := NewPaymentMethod("credit")
	if pm == nil {
		fmt.Println("Неизвестный тип платежа")
		return
	}
	pm.Process(100.0)
}
```

**Комментарий:**  
Выбор способа оплаты делегируется фабрике, а каждая реализация самостоятельно знает, как обработать платеж. Таким образом, в основной логике не требуется ветвление по типу платежа.

---

## 3. Логирование в зависимости от окружения (Dependency Injection)

**Было (условная логика по окружению):**

```go
package main

import "fmt"

func Log(message, env string) {
	if env == "prod" {
		// Отправляем лог на удалённый сервер
		fmt.Println("[REMOTE]", message)
	} else {
		// Логируем в консоль
		fmt.Println("[LOCAL]", message)
	}
}

func main() {
	Log("Событие произошло", "dev")
}
```

**Стало (использование интерфейса Logger, внедряемого через DI):**

```go
package main

import "fmt"

// Интерфейс логгера
type Logger interface {
	Log(message string)
}

type RemoteLogger struct{}
func (r RemoteLogger) Log(message string) {
	// Реальная отправка лога на сервер
	fmt.Println("[REMOTE]", message)
}

type LocalLogger struct{}
func (l LocalLogger) Log(message string) {
	fmt.Println("[LOCAL]", message)
}

// Клиентский код использует Logger, не заботясь о реализации
func ProcessEvent(logger Logger, event string) {
	logger.Log("Событие: " + event)
}

func main() {
	// Выбор логгера происходит один раз при инициализации
	var logger Logger = LocalLogger{}
	// Для production можно было бы сделать: var logger Logger = RemoteLogger{}
	ProcessEvent(logger, "Событие произошло")
}
```

**Комментарий:**  
Здесь выбор способа логирования происходит на этапе конфигурации, а основной код просто вызывает метод Log интерфейса, что устраняет необходимость в условных проверках при каждом вызове.

---

## 4. Обработка цепочки трансформаций (Функциональная композиция)

**Было (последовательные if-проверки при обработке данных):**

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

func toUpper(s string) (string, error) {
	if s == "" {
		return "", errors.New("пустая строка")
	}
	return strings.ToUpper(s), nil
}

func addExclamation(s string) (string, error) {
	if len(s) > 20 {
		return "", errors.New("слишком длинная строка")
	}
	return s + "!", nil
}

func processData(data string) (string, error) {
	res, err := toUpper(data)
	if err != nil {
		return "", err
	}
	res, err = addExclamation(res)
	if err != nil {
		return "", err
	}
	return res, nil
}

func main() {
	result, err := processData("hello")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(result)
}
```

**Стало (композиция функций через общий пайплайн):**

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

// Тип функции-трансформера
type Transformer func(string) (string, error)

func toUpper(s string) (string, error) {
	if s == "" {
		return "", errors.New("пустая строка")
	}
	return strings.ToUpper(s), nil
}

func addExclamation(s string) (string, error) {
	if len(s) > 20 {
		return "", errors.New("слишком длинная строка")
	}
	return s + "!", nil
}

// Chain объединяет последовательность трансформеров в один
func Chain(transformers ...Transformer) Transformer {
	return func(input string) (string, error) {
		var err error
		for _, t := range transformers {
			input, err = t(input)
			if err != nil {
				return "", err
			}
		}
		return input, nil
	}
}

func main() {
	// Определяем цепочку преобразований
	processData := Chain(toUpper, addExclamation)
	result, err := processData("hello")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println(result)
}
```

**Комментарий:**  
Пайплайн (цепочка) трансформаций инкапсулирует последовательные шаги обработки данных. Основной код не содержит явных if для каждого шага — ошибки обрабатываются внутри композиции.

---

## 5. Реализация логики повторных попыток (Strategy Pattern)

**Было (цикл с проверками для повторной попытки операции):**

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

func fetchData() (string, error) {
	// Симуляция неудачного запроса
	return "", errors.New("сбой запроса")
}

func fetchDataWithRetry(maxAttempts int) (string, error) {
	attempt := 0
	for {
		data, err := fetchData()
		if err == nil {
			return data, nil
		}
		if attempt >= maxAttempts {
			return "", err
		}
		attempt++
		time.Sleep(1 * time.Second)
	}
}

func main() {
	data, err := fetchDataWithRetry(3)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Полученные данные:", data)
}
```

**Стало (инкапсуляция логики повторов в стратегию):**

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

// Интерфейс стратегии повторных попыток
type RetryStrategy interface {
	Execute(operation func() (string, error)) (string, error)
}

type ExponentialBackoff struct {
	maxAttempts int
	baseDelay   time.Duration
}

func (r ExponentialBackoff) Execute(operation func() (string, error)) (string, error) {
	var err error
	delay := r.baseDelay
	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		result, opErr := operation()
		if opErr == nil {
			return result, nil
		}
		err = opErr
		time.Sleep(delay)
		delay *= 2
	}
	return "", err
}

func fetchData() (string, error) {
	// Симуляция неудачного запроса
	return "", errors.New("сбой запроса")
}

func main() {
	strategy := ExponentialBackoff{
		maxAttempts: 3,
		baseDelay:   500 * time.Millisecond,
	}
	// Логика повторов инкапсулирована в стратегии, основному коду не нужно проверять условия.
	data, err := strategy.Execute(fetchData)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Полученные данные:", data)
}
```

**Комментарий:**  
Логика повторных попыток вынесена в отдельную стратегию, которую можно легко заменить или настроить. Основной код просто вызывает стратегию, избавляясь от явных циклических проверок.