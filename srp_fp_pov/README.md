# SRP с т.зр. ФП

### 1. Валидация и обработка входящих данных (HTTP API)
**Проблема**: Обработчик эндпоинта делал слишком много: парсинг JSON, валидацию, преобразование данных и бизнес-логику.

**SRP+ФП решение**:
```go
// Чистая функция для валидации email
func validateEmail(email string) error {
    if !regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+`).MatchString(email) {
        return errors.New("invalid email format")
    }
    return nil
}

// Чистая функция для нормализации email (приведение к нижнему регистру)
func normalizeEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

// Чистая функция для создания пользователя (бизнес-логика)
func createUser(email string, name string) (*User, error) {
    return &User{
        ID:    generateID(),
        Email: email,
        Name:  name,
    }, nil
}

// Композиция в обработчике
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var request struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Композиция чистых функций
    if err := validateEmail(request.Email); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    normalizedEmail := normalizeEmail(request.Email)
    user, err := createUser(normalizedEmail, request.Name)
    if err != nil {
        http.Error(w, "Creation failed", http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}
```

**Преимущества**:
- Каждая функция имеет одну ответственность
- Легко тестировать изолированно
- `normalizeEmail` и `validateEmail` - чистые функции

---

### 2. Трансформация данных в конвейере обработки
**Проблема**: Монолитная функция обработки данных, которую сложно модифицировать.

**SRP+ФП решение**:
```go
// Тип для чистых функций трансформации
type StringTransformer func(string) string

// Маленькие функции с одной ответственностью
func trimSpaces(s string) string {
    return strings.TrimSpace(s)
}

func removeSpecialChars(s string) string {
    return regexp.MustCompile(`[^a-zA-Z0-9\s]`).ReplaceAllString(s, "")
}

func toLowerCase(s string) string {
    return strings.ToLower(s)
}

// Композиция трансформаций
func processString(input string, transformers ...StringTransformer) string {
    result := input
    for _, transform := range transformers {
        result = transform(result)
    }
    return result
}

// Использование
func main() {
    input := "  Hello, World!  "
    
    result := processString(input,
        trimSpaces,
        removeSpecialChars,
        toLowerCase,
    )
    
    fmt.Println(result) // "hello world"
}
```

---

### 3. Фильтрация и агрегация данных
**Проблема**: Сложный цикл с множеством условий и побочных эффектов.

**SRP+ФП решение**:
```go
// Чистые функции-предикаты
type ProductFilter func(Product) bool

func byCategory(category string) ProductFilter {
    return func(p Product) bool {
        return p.Category == category
    }
}

func inPriceRange(min, max float64) ProductFilter {
    return func(p Product) bool {
        return p.Price >= min && p.Price <= max
    }
}

func isAvailable() ProductFilter {
    return func(p Product) bool {
        return p.StockCount > 0
    }
}

// Композиция фильтров
func filterProducts(products []Product, filters ...ProductFilter) []Product {
    var result []Product
    for _, product := range products {
        valid := true
        for _, filter := range filters {
            if !filter(product) {
                valid = false
                break
            }
        }
        if valid {
            result = append(result, product)
        }
    }
    return result
}

// Чистая функция для агрегации
func calculateTotalPrice(products []Product) float64 {
    total := 0.0
    for _, p := range products {
        total += p.Price
    }
    return total
}

// Использование
availableElectronics := filterProducts(products,
    byCategory("electronics"),
    inPriceRange(100, 1000),
    isAvailable(),
)

total := calculateTotalPrice(availableElectronics)
```

---

### 4. Обработка ошибок с функциональными обёртками
**Проблема**: Повторяющийся код обработки ошибок.

**SRP+ФП решение**:
```go
// Функция высшего порядка для обработки ошибок
func withErrorHandling(fn func() error, onError func(error)) {
    if err := fn(); err != nil {
        onError(err)
    }
}

// Специализированные обработчики с одной ответственностью
func logError(err error) {
    log.Printf("Error occurred: %v", err)
}

func metricsError(err error) {
    metrics.Increment("error_count")
}

func notifyTeam(err error) {
    if isCritical(err) {
        slack.SendAlert(err.Error())
    }
}

// Композиция обработчиков
func handleDatabaseOperation(op func() error) {
    withErrorHandling(op, func(err error) {
        logError(err)
        metricsError(err)
        notifyTeam(err)
    })
}

// Использование
func updateUserProfile(userID string, profile Profile) {
    handleDatabaseOperation(func() error {
        return db.Users.Update(userID, profile)
    })
}
```

---

### 5. Функциональный пайплайн для ETL
**Проблема**: Сложный ETL-процесс с перемешанной логикой.

**SRP+ФП решение**:
```go
// Типы для функционального пайплайна
type Extractor func() ([]Data, error)
type Transformer func(Data) Data
type Loader func(Data) error

// Чистые функции-трансформеры
func sanitizeData(d Data) Data {
    return Data{
        ID:    strings.TrimSpace(d.ID),
        Value: math.Max(0, d.Value), // Убираем отрицательные значения
    }
}

func enrichWithTimestamp(d Data) Data {
    d.ProcessedAt = time.Now()
    return d
}

func calculateDerivedFields(d Data) Data {
    d.DerivedValue = d.Value * 1.1 // Добавляем 10%
    return d
}

// Композиция пайплайна
func runETLPipeline(extract Extractor, transform []Transformer, load Loader) error {
    data, err := extract()
    if err != nil {
        return err
    }
    
    for i := range data {
        // Применяем все трансформеры последовательно
        for _, t := range transform {
            data[i] = t(data[i])
        }
        
        if err := load(data[i]); err != nil {
            return err
        }
    }
    
    return nil
}

// Использование
func main() {
    transformers := []Transformer{
        sanitizeData,
        enrichWithTimestamp,
        calculateDerivedFields,
    }
    
    err := runETLPipeline(
        extractFromCSV,    // Функция-экстрактор
        transformers,      // Цепочка трансформеров
        saveToDatabase,    // Функция-загрузчик
    )
}
```