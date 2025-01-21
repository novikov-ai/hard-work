# Добиваемся совместимости между (унаследованными) типами

## Сохранение событий для ретроспективного анализа

В сервисе с высокой нагрузкой часто требуется логировать важные события для их последующего анализа или восстановления состояния. Например, это может быть система учёта заказов в маркетплейсе. Вместо хранения только текущего состояния заказов, можно логировать все изменения.

```go
type Event struct {
    Timestamp time.Time
    Action    string
    Metadata  map[string]string
}

type EventJournal struct {
    events []Event
}

func (j *EventJournal) RecordEvent(action string, metadata map[string]string) {
    j.events = append(j.events, Event{
        Timestamp: time.Now(),
        Action:    action,
        Metadata:  metadata,
    })
}

func (j *EventJournal) RollbackTo(index int) {
    if index >= 0 && index < len(j.events) {
        j.events = j.events[:index+1]
    }
}
```


## Миграция данных между форматами с сохранением обратной совместимости

Если сервис изменяет формат хранения данных (например, с JSON на Protobuf), может возникнуть необходимость поддерживать оба формата в течение переходного периода. Введение слоя адаптации позволяет сохранить данные в универсальном виде.

```go
type User struct {
    ID    string
    Name  string
    Email string
}

// Логирование действий с пользователем
type UserAction struct {
    User      User
    Timestamp time.Time
    Action    string
}

// Абстракция над форматами хранения
type DataAdapter interface {
    Save(user UserAction) error
    Load(userID string) ([]UserAction, error)
}

type JSONAdapter struct {}
type ProtobufAdapter struct {}

// Пример реализации адаптера для JSON
func (a *JSONAdapter) Save(userAction UserAction) error {
    data, err := json.Marshal(userAction)
    if err != nil {
        return err
    }
    fmt.Println("Saving JSON:", string(data))
    return nil
}

// Аналогично для Protobuf...
```

## Реализация ретроактивного пересчёта данных

В аналитических системах может возникнуть необходимость пересчитать старые данные с учётом новых правил (например, изменения в бизнес-логике). Хранение операций вместо их результатов помогает в таких случаях.

```go
type Operation struct {
    Name      string
    Arguments []int
}

type Calculator struct {
    operations []Operation
}

func (c *Calculator) Add(a, b int) int {
    c.operations = append(c.operations, Operation{"Add", []int{a, b}})
    return a + b
}

func (c *Calculator) Recalculate() int {
    result := 0
    for _, op := range c.operations {
        switch op.Name {
        case "Add":
            result += op.Arguments[0] + op.Arguments[1]
        }
    }
    return result
}
```