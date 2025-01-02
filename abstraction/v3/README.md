# Обобщаем проектные абстракции

### Пример 1: Иерархия сущностей в системе управления складом

Исходный вариант:
- Entity (базовый класс).
- Item (предметы).
  - PerishableItem (скоропортящиеся товары).
  - NonPerishableItem (нескоропортящиеся товары).

Обоснование: 
Классы PerishableItem и NonPerishableItem используют разные операции для определения срока годности, но больше никак не связаны. Их объединение в Item лишь усложняет систему.

Решение:
1. Удалить Item как общий класс.
2. Создать интерфейс Expirable, реализуемый каждым классом.

~~~ go
type Expirable interface {
    ExpirationDate() time.Time
}

type PerishableItem struct { /* поля и методы */ }
func (p PerishableItem) ExpirationDate() time.Time { /* реализация */ }

type NonPerishableItem struct { /* поля и методы */ }
func (n NonPerishableItem) ExpirationDate() time.Time { /* реализация */ }
~~~

### Пример 2: Иерархия пользовательских уведомлений

Исходный вариант:
- Notification (базовый класс).
  - EmailNotification.
  - SMSNotification.
  - PushNotification.

Обоснование:
Классы уведомлений никак не связаны по поведению. Их объединение в Notification добавляет избыточность.

Решение:
Вместо наследования использовать интерфейс Notifier:

~~~go
type Notifier interface {
    Notify(message string) error
}

type EmailNotification struct { /* поля и методы */ }
func (e EmailNotification) Notify(message string) error { /* отправка email */ }

type SMSNotification struct { /* поля и методы */ }
func (s SMSNotification) Notify(message string) error { /* отправка SMS */ }

type PushNotification struct { /* поля и методы */ }
func (s PushNotification) Notify(message string) error { /* отправка Push */ }
~~~

### Расширение подхода через interface dispatch

Пример: Интерфейсы для операций с геометрическими фигурами

Исходная иерархия
- Shape (базовый класс).
  - Circle.
  - Rectangle.
  - Triangle.

Проблема:
Разные фигуры имеют уникальные операции, такие как расчёт площади или периметра. Общий класс затрудняет их поддержку.

Решение:
Создать отдельные интерфейсы:

~~~go
type AreaCalculator interface {
    Area() float64
}

type PerimeterCalculator interface {
    Perimeter() float64
}

type Circle struct { /* поля */ }
func (c Circle) Area() float64       { /* реализация */ }
func (c Circle) Perimeter() float64  { /* реализация */ }

type Rectangle struct { /* поля */ }
func (r Rectangle) Area() float64    { /* реализация */ }
func (r Rectangle) Perimeter() float64 { /* реализация */ }
~~~