# Модули важнее всего

## 1. **Источник данных: БД и кэш**
Реализация (`MemoryStore`) удовлетворяет интерфейсам `Database` и `Cache`:

```go
type Database interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte) error
}

type Cache interface {
    GetCached(key string) ([]byte, error)
}

// Одна реализация для двух интерфейсов
type MemoryStore struct {
    data map[string][]byte
}

func (m *MemoryStore) Get(key string) ([]byte, error) { 
    return m.data[key], nil 
}

func (m *MemoryStore) Set(key string, value []byte) error { 
    m.data[key] = value; return nil 
}

func (m *MemoryStore) GetCached(key string) ([]byte, error) {
    // Логика кэширования
    return m.Get(key)
}

// Использование
func SaveToDB(db Database) { /* ... */ }
func UseCache(cache Cache)  { /* ... */ }

func main() {
    store := &MemoryStore{data: make(map[string][]byte)}
    SaveToDB(store)  // Используем как Database
    UseCache(store)  // Используем как Cache
}
```

---

## 2. **Логгер: Консоль и файл**
Реализация (`MultiLogger`) удовлетворяет интерфейсам `ConsoleLogger` и `FileLogger`:

```go
type ConsoleLogger interface {
    LogToConsole(message string)
}

type FileLogger interface {
    LogToFile(message string) error
}

// Одна реализация для двух интерфейсов
type MultiLogger struct{}

func (m *MultiLogger) LogToConsole(message string) {
    fmt.Println("CONSOLE:", message)
}

func (m *MultiLogger) LogToFile(message string) error {
    // Запись в файл
    return nil
}

// Использование
func Debug(console ConsoleLogger) { console.LogToConsole("Debug") }
func Audit(file FileLogger)      { file.LogToFile("Audit") }

func main() {
    logger := &MultiLogger{}
    Debug(logger)  // Используем как ConsoleLogger
    Audit(logger)  // Используем как FileLogger
}
```

---

## 3. **Геометрическая фигура: Площадь и периметр**
Реализация (`Circle`) удовлетворяет интерфейсам `AreaCalculator` и `PerimeterCalculator`:

```go
type AreaCalculator interface {
    Area() float64
}

type PerimeterCalculator interface {
    Perimeter() float64
}

// Одна реализация для двух интерфейсов
type Circle struct{ Radius float64 }

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Использование
func PrintArea(a AreaCalculator) { fmt.Println("Area:", a.Area()) }
func PrintPerimeter(p PerimeterCalculator) { fmt.Println("Perimeter:", p.Perimeter()) }

func main() {
    circle := Circle{Radius: 5}
    PrintArea(circle)       // Используем как AreaCalculator
    PrintPerimeter(circle)  // Используем как PerimeterCalculator
}
```