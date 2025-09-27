# Что такое абстрация-2

## Пример 1: Абстракция для работы с денежными суммами

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
)

// Конкретное представление (вычислительное) - копейки
type Cents int64

// Абстрактное представление - денежная сумма с валютой
type Money struct {
    Amount   float64
    Currency string
}

// LIFT: Поднимаем конкретное значение в абстрактное представление
func LiftToMoney(cents Cents, currency string) Money {
    return Money{
        Amount:   float64(cents) / 100.0,
        Currency: currency,
    }
}

// LOWER: Опускаем абстрактное представление в конкретное
func (m Money) LowerToCents() Cents {
    return Cents(m.Amount * 100)
}

// Операции на абстрактном уровне (более понятные для человека)
func (m Money) Add(other Money) Money {
    if m.Currency != other.Currency {
        panic("Валюты не совпадают")
    }
    return Money{
        Amount:   m.Amount + other.Amount,
        Currency: m.Currency,
    }
}

func (m Money) Format() string {
    return fmt.Sprintf("%s %.2f", m.Currency, m.Amount)
}

func main() {
    // Конкретные значения (удобные для вычислений)
    price1 := Cents(1999) // 19.99 в копейках
    price2 := Cents(500)  // 5.00 в копейках
    
    // LIFT: Поднимаем в абстрактное представление
    money1 := LiftToMoney(price1, "USD")
    money2 := LiftToMoney(price2, "USD")
    
    fmt.Printf("Товар 1: %s\n", money1.Format()) // Товар 1: USD 19.99
    fmt.Printf("Товар 2: %s\n", money2.Format()) // Товар 2: USD 5.00
    
    // Работа на абстрактном уровне
    total := money1.Add(money2)
    fmt.Printf("Итого: %s\n", total.Format()) // Итого: USD 24.99
    
    // LOWER: Опускаем обратно для вычислений
    totalCents := total.LowerToCents()
    fmt.Printf("Всего копеек: %d\n", totalCents) // Всего копеек: 2499
}
```

## Пример 2: Абстракция для работы с температурой

```go
package main

import (
    "fmt"
    "math"
)

// Конкретное представление - Кельвины (удобно для вычислений)
type Kelvin float64

// Абстрактное представление - температура со шкалой
type Temperature struct {
    Value float64
    Scale string // "C", "F", "K"
}

// LIFT функции для разных шкал
func LiftCelsiusToTemp(c float64) Temperature {
    return Temperature{Value: c, Scale: "C"}
}

func LiftFahrenheitToTemp(f float64) Temperature {
    return Temperature{Value: f, Scale: "F"}
}

func LiftKelvinToTemp(k Kelvin) Temperature {
    return Temperature{Value: float64(k), Scale: "K"}
}

// LOWER: Опускаем в Кельвины (универсальная вычислительная форма)
func (t Temperature) LowerToKelvin() Kelvin {
    switch t.Scale {
    case "K":
        return Kelvin(t.Value)
    case "C":
        return Kelvin(t.Value + 273.15)
    case "F":
        return Kelvin((t.Value-32)*5/9 + 273.15)
    default:
        panic("Неизвестная шкала")
    }
}

// Операции на абстрактном уровне
func (t Temperature) ConvertTo(scale string) Temperature {
    kelvin := t.LowerToKelvin() // Сначала опускаем к универсальной форме
    
    switch scale {
    case "K":
        return Temperature{Value: float64(kelvin), Scale: "K"}
    case "C":
        return Temperature{Value: float64(kelvin) - 273.15, Scale: "C"}
    case "F":
        return Temperature{Value: (float64(kelvin)-273.15)*9/5 + 32, Scale: "F"}
    default:
        panic("Неизвестная шкала")
    }
}

func (t Temperature) String() string {
    return fmt.Sprintf("%.1f°%s", t.Value, t.Scale)
}

func main() {
    // Конкретные значения от разных источников
    roomTempC := 23.0
    weatherTempF := 68.0
    scienceTempK := Kelvin(300)
    
    // LIFT: Поднимаем к абстрактному представлению
    temp1 := LiftCelsiusToTemp(roomTempC)
    temp2 := LiftFahrenheitToTemp(weatherTempF)
    temp3 := LiftKelvinToTemp(scienceTempK)
    
    fmt.Printf("Комнатная температура: %s\n", temp1) // Комнатная температура: 23.0°C
    fmt.Printf("Погода на улице: %s\n", temp2)       // Погода на улице: 68.0°F
    fmt.Printf("Научный эксперимент: %s\n", temp3)   // Научный эксперимент: 300.0°K
    
    // Работаем на абстрактном уровне
    temp1InF := temp1.ConvertTo("F")
    fmt.Printf("Комнатная температура по Фаренгейту: %s\n", temp1InF)
    
    // LOWER: Для физических расчетов
    kelvin1 := temp1.LowerToKelvin()
    kelvin2 := temp2.LowerToKelvin()
    
    avgKelvin := (kelvin1 + kelvin2) / 2
    fmt.Printf("Средняя температура в Кельвинах: %.2f\n", avgKelvin)
}
```

## Пример 3: Абстракция для обработки пользовательского ввода

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

// Конкретное представление - сырые строки из формы
type RawFormData map[string]string

// Абстрактное представление - валидированные данные пользователя
type UserProfile struct {
    Name      string
    Email     string
    BirthDate time.Time
    Age       int
}

// LIFT: Поднимаем сырые данные в абстрактную модель
func LiftToUserProfile(raw RawFormData) (UserProfile, error) {
    var profile UserProfile
    var errors []string
    
    // Валидация и преобразование имени
    if name, ok := raw["name"]; ok && len(strings.TrimSpace(name)) > 0 {
        profile.Name = strings.TrimSpace(name)
    } else {
        errors = append(errors, "Имя обязательно")
    }
    
    // Валидация email
    if email, ok := raw["email"]; ok && strings.Contains(email, "@") {
        profile.Email = strings.TrimSpace(email)
    } else {
        errors = append(errors, "Некорректный email")
    }
    
    // Парсинг даты рождения
    if dateStr, ok := raw["birth_date"]; ok {
        if birthDate, err := time.Parse("2006-01-02", dateStr); err == nil {
            profile.BirthDate = birthDate
            profile.Age = calculateAge(birthDate)
        } else {
            errors = append(errors, "Некорректная дата рождения")
        }
    }
    
    if len(errors) > 0 {
        return profile, fmt.Errorf("Ошибки валидации: %s", strings.Join(errors, ", "))
    }
    
    return profile, nil
}

// LOWER: Опускаем абстрактную модель к форме, удобной для БД
func (p UserProfile) LowerToDBFormat() map[string]interface{} {
    return map[string]interface{}{
        "name":       p.Name,
        "email":      p.Email,
        "birth_date": p.BirthDate.Format("2006-01-02"),
        "age":        p.Age,
        "created_at": time.Now(),
    }
}

// Операции на абстрактном уровне
func (p UserProfile) IsAdult() bool {
    return p.Age >= 18
}

func (p UserProfile) WelcomeMessage() string {
    return fmt.Sprintf("Добро пожаловать, %s! Ваш email: %s", p.Name, p.Email)
}

func calculateAge(birthDate time.Time) int {
    today := time.Now()
    age := today.Year() - birthDate.Year()
    if today.YearDay() < birthDate.YearDay() {
        age--
    }
    return age
}

func main() {
    // Конкретные данные из HTML-формы
    rawData := RawFormData{
        "name":       "  Иван Петров  ",
        "email":      "ivan@example.com",
        "birth_date": "1990-05-15",
    }
    
    // LIFT: Поднимаем к абстрактному представлению
    profile, err := LiftToUserProfile(rawData)
    if err != nil {
        fmt.Println("Ошибка:", err)
        return
    }
    
    // Работаем на абстрактном уровне
    fmt.Println(profile.WelcomeMessage())
    fmt.Printf("Возраст: %d, Совершеннолетний: %t\n", profile.Age, profile.IsAdult())
    
    // LOWER: Опускаем к формату для базы данных
    dbData := profile.LowerToDBFormat()
    fmt.Printf("Данные для БД: %+v\n", dbData)
}
```

## Пример 4: Абстракция для математических выражений

```go
package main

import (
    "fmt"
    "strconv"
)

// Конкретное представление - токены парсера
type Token struct {
    Type  string // "number", "operator", "paren"
    Value string
}

// Абстрактное представление - AST (Abstract Syntax Tree)
type Expr interface {
    Eval() float64
}

type Number struct{ Value float64 }
type BinaryOp struct {
    Left  Expr
    Op    string
    Right Expr
}

func (n Number) Eval() float64 { return n.Value }
func (b BinaryOp) Eval() float64 {
    switch b.Op {
    case "+": return b.Left.Eval() + b.Right.Eval()
    case "-": return b.Left.Eval() - b.Right.Eval()
    case "*": return b.Left.Eval() * b.Right.Eval()
    case "/": return b.Left.Eval() / b.Right.Eval()
    default: panic("Неизвестный оператор")
    }
}

// LIFT: Поднимаем токены в AST
func LiftToAST(tokens []Token) Expr {
    return parseExpression(tokens)
}

// LOWER: Опускаем AST обратно к строковому представлению
func LowerToString(expr Expr) string {
    switch e := expr.(type) {
    case Number:
        return fmt.Sprintf("%g", e.Value)
    case BinaryOp:
        return fmt.Sprintf("(%s %s %s)", LowerToString(e.Left), e.Op, LowerToString(e.Right))
    default:
        return "???"
    }
}

// Парсер (упрощенный)
func parseExpression(tokens []Token) Expr {
    if len(tokens) == 1 && tokens[0].Type == "number" {
        val, _ := strconv.ParseFloat(tokens[0].Value, 64)
        return Number{Value: val}
    }
    
    // Ищем оператор с наименьшим приоритетом
    parenLevel := 0
    for i := len(tokens) - 1; i >= 0; i-- {
        token := tokens[i]
        
        if token.Type == "paren" {
            if token.Value == ")" { parenLevel++ }
            if token.Value == "(" { parenLevel-- }
            continue
        }
        
        if parenLevel == 0 && token.Type == "operator" && (token.Value == "+" || token.Value == "-") {
            return BinaryOp{
                Left:  parseExpression(tokens[:i]),
                Op:    token.Value,
                Right: parseExpression(tokens[i+1:]),
            }
        }
    }
    
    // Упрощение: обрабатываем только простые случаи
    if tokens[0].Value == "(" && tokens[len(tokens)-1].Value == ")" {
        return parseExpression(tokens[1 : len(tokens)-1])
    }
    
    panic("Не могу распарсить выражение")
}

func main() {
    // Конкретное представление - токены
    tokens := []Token{
        {Type: "paren", Value: "("},
        {Type: "number", Value: "10"},
        {Type: "operator", Value: "+"},
        {Type: "number", Value: "5"},
        {Type: "operator", Value: "*"},
        {Type: "number", Value: "2"},
        {Type: "paren", Value: ")"},
    }
    
    // LIFT: Поднимаем к абстрактному синтаксическому дереву
    ast := LiftToAST(tokens)
    
    // Работаем на абстрактном уровне
    result := ast.Eval()
    fmt.Printf("Результат вычисления: %g\n", result) // Результат вычисления: 20
    
    // LOWER: Опускаем обратно к представлению (упрощенному)
    expressionStr := LowerToString(ast)
    fmt.Printf("Строковое представление: %s\n", expressionStr)
}
```