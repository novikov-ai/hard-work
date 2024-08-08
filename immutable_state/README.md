# Имутабельные состояния лучше, чем передача сущностей по ссылке

### Пример 1: Работа с данными пользователя в микросервисах

#### До: Модификация пользовательского объекта по ссылке

В типичном микросервисном приложении может возникнуть ситуация, когда данные пользователя передаются по ссылке между несколькими сервисами, и каждый сервис может модифицировать эти данные.

```go
package main

import (
	"fmt"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	UpdatedAt time.Time
}

func updateUserEmail(user *User, newEmail string) {
	user.Email = newEmail
	user.UpdatedAt = time.Now()
}

func updateUserName(user *User, newName string) {
	user.Name = newName
	user.UpdatedAt = time.Now()
}

func main() {
	user := &User{ID: "123", Name: "Alice", Email: "alice@example.com", UpdatedAt: time.Now()}
	updateUserEmail(user, "alice.new@example.com")
	updateUserName(user, "Alice New")
	fmt.Println(user) // User{Name: "Alice New", Email: "alice.new@example.com", UpdatedAt: ...}
}
```

**Проблема:** Модификация данных пользователя по ссылке приводит к тому, что разные сервисы могут вносить изменения в один и тот же объект, что затрудняет отслеживание изменений и может привести к состояниям гонки в многопоточных приложениях.

#### После: Использование иммутабельного объекта

Вместо изменения оригинального объекта, каждый сервис будет создавать новый объект с обновленными данными.

```go
package main

import (
	"fmt"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	UpdatedAt time.Time
}

func (u User) withUpdatedEmail(newEmail string) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     newEmail,
		UpdatedAt: time.Now(),
	}
}

func (u User) withUpdatedName(newName string) User {
	return User{
		ID:        u.ID,
		Name:      newName,
		Email:     u.Email,
		UpdatedAt: time.Now(),
	}
}

func main() {
	user := User{ID: "123", Name: "Alice", Email: "alice@example.com", UpdatedAt: time.Now()}
	updatedUser := user.withUpdatedEmail("alice.new@example.com")
	updatedUser = updatedUser.withUpdatedName("Alice New")
	fmt.Println(updatedUser) // User{Name: "Alice New", Email: "alice.new@example.com", UpdatedAt: ...}
	fmt.Println(user)        // User{Name: "Alice", Email: "alice@example.com", UpdatedAt: ...}
}
```

**Комментарий:** Создание новых объектов при внесении изменений позволяет избежать непреднамеренных побочных эффектов и улучшает предсказуемость системы, особенно в распределённых и многопоточных приложениях.

### Пример 2: Управление конфигурациями в микросервисной архитектуре

#### До: Изменение общей конфигурации по ссылке

В системе, где несколько микросервисов используют общий объект конфигурации, изменение конфигурации по ссылке может привести к непредсказуемому поведению.

```go
package main

import (
	"fmt"
	"time"
)

type Config struct {
	DatabaseURL string
	CacheTTL    time.Duration
}

func updateDatabaseURL(config *Config, newURL string) {
	config.DatabaseURL = newURL
}

func updateCacheTTL(config *Config, newTTL time.Duration) {
	config.CacheTTL = newTTL
}

func main() {
	config := &Config{DatabaseURL: "postgres://localhost:5432/mydb", CacheTTL: 10 * time.Minute}
	updateDatabaseURL(config, "postgres://localhost:5432/newdb")
	updateCacheTTL(config, 20*time.Minute)
	fmt.Println(config) // Config{DatabaseURL: "postgres://localhost:5432/newdb", CacheTTL: 20m}
}
```

**Проблема:** Изменение конфигурации может неожиданно повлиять на другие части системы, которые используют этот объект, особенно если эти изменения происходят асинхронно.

#### После: Иммутабельная конфигурация

Вместо изменения общей конфигурации создаём новый объект конфигурации с необходимыми изменениями.

```go
package main

import (
	"fmt"
	"time"
)

type Config struct {
	DatabaseURL string
	CacheTTL    time.Duration
}

func (c Config) withUpdatedDatabaseURL(newURL string) Config {
	return Config{
		DatabaseURL: newURL,
		CacheTTL:    c.CacheTTL,
	}
}

func (c Config) withUpdatedCacheTTL(newTTL time.Duration) Config {
	return Config{
		DatabaseURL: c.DatabaseURL,
		CacheTTL:    newTTL,
	}
}

func main() {
	config := Config{DatabaseURL: "postgres://localhost:5432/mydb", CacheTTL: 10 * time.Minute}
	newConfig := config.withUpdatedDatabaseURL("postgres://localhost:5432/newdb")
	newConfig = newConfig.withUpdatedCacheTTL(20 * time.Minute)
	fmt.Println(config)    // Config{DatabaseURL: "postgres://localhost:5432/mydb", CacheTTL: 10m}
	fmt.Println(newConfig) // Config{DatabaseURL: "postgres://localhost:5432/newdb", CacheTTL: 20m}
}
```

**Комментарий:** Иммутабельные объекты конфигурации уменьшают вероятность возникновения ошибок, связанных с изменением конфигурации во время выполнения программы, и упрощают управление конфигурацией в распределённых системах.

### Пример 3: Обработка транзакций в банковской системе

#### До: Модификация состояния транзакции по ссылке

В банковской системе транзакции могут обрабатываться последовательно, и каждая транзакция может изменять состояние предыдущей.

```go
package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        string
	Amount    float64
	Status    string
	Timestamp time.Time
}

func processTransaction(transaction *Transaction, status string) {
	transaction.Status = status
	transaction.Timestamp = time.Now()
}

func main() {
	transaction := &Transaction{ID: "tx123", Amount: 100.0, Status: "pending"}
	processTransaction(transaction, "completed")
	fmt.Println(transaction) // Transaction{ID: "tx123", Amount: 100.0, Status: "completed", Timestamp: ...}
}
```

**Проблема:** Изменение статуса транзакции напрямую может привести к тому, что несколько операций будут работать с одним и тем же объектом, что затрудняет отслеживание изменений и может привести к ошибкам.

#### После: Использование иммутабельного состояния транзакции

Вместо изменения состояния транзакции создаётся новая транзакция с обновлённым состоянием.

```go
package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        string
	Amount    float64
	Status    string
	Timestamp time.Time
}

func (t Transaction) withUpdatedStatus(status string) Transaction {
	return Transaction{
		ID:        t.ID,
		Amount:    t.Amount,
		Status:    status,
		Timestamp: time.Now(),
	}
}

func main() {
	transaction := Transaction{ID: "tx123", Amount: 100.0, Status: "pending"}
	newTransaction := transaction.withUpdatedStatus("completed")
	fmt.Println(transaction)      // Transaction{ID: "tx123", Amount: 100.0, Status: "pending", Timestamp: ...}
	fmt.Println(newTransaction)   // Transaction{ID: "tx123", Amount: 100.0, Status: "completed", Timestamp: ...}
}
```

**Комментарий:** Использование иммутабельных состояний транзакций помогает избежать ошибок, связанных с параллельной обработкой данных, и облегчает отслеживание и аудит изменений.