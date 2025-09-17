# Функциональные интерфейсы

## Пример 1: Сервис аутентификации

### До (объектно-императивный стиль)
```go
type AuthService struct {
    userRepo  *UserRepository
    validator *Validator
    tokenTTL  time.Duration
}

func NewAuthService(userRepo *UserRepository, validator *Validator) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        validator: validator,
        tokenTTL:  24 * time.Hour,
    }
}

func (s *AuthService) Login(email, password string) (string, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", err
    }
    
    if !s.validator.ValidatePassword(user.PasswordHash, password) {
        return "", errors.New("invalid password")
    }
    
    token, err := s.generateToken(user.ID)
    if err != nil {
        return "", err
    }
    
    return token, nil
}

func (s *AuthService) generateToken(userID string) (string, error) {
    // Императивная генерация токена
    // ...
}
```

### После (функциональный интерфейс с IoC)
```go
// Определяем интерфейсы для зависимостей
type UserProvider interface {
    FindByEmail(email string) (*User, error)
}

type PasswordValidator interface {
    ValidatePassword(hash, password string) bool
}

type TokenGenerator interface {
    GenerateToken(userID string, ttl time.Duration) (string, error)
}

// Чисто функциональный интерфейс
type AuthDependencies struct {
    UserProvider      UserProvider
    PasswordValidator PasswordValidator
    TokenGenerator    TokenGenerator
    TokenTTL          time.Duration
}

func Login(deps AuthDependencies, email, password string) (string, error) {
    user, err := deps.UserProvider.FindByEmail(email)
    if err != nil {
        return "", err
    }
    
    if !deps.PasswordValidator.ValidatePassword(user.PasswordHash, password) {
        return "", errors.New("invalid password")
    }
    
    token, err := deps.TokenGenerator.GenerateToken(user.ID, deps.TokenTTL)
    if err != nil {
        return "", err
    }
    
    return token, nil
}

// Каррированная версия для удобства использования
type LoginFunc func(email, password string) (string, error)

func CreateLoginFunc(deps AuthDependencies) LoginFunc {
    return func(email, password string) (string, error) {
        return Login(deps, email, password)
    }
}
```

## Пример 2: Обработчик данных

### До (объектно-императивный стиль)
```go
type DataProcessor struct {
    transformer *DataTransformer
    validator   *DataValidator
    repository  *DataRepository
}

func (p *DataProcessor) Process(data []byte) error {
    if !p.validator.Validate(data) {
        return errors.New("invalid data")
    }
    
    transformed, err := p.transformer.Transform(data)
    if err != nil {
        return err
    }
    
    return p.repository.Save(transformed)
}
```

### После (функциональный интерфейс с IoC)
```go
// Определяем интерфейсы для зависимостей
type DataValidator interface {
    Validate(data []byte) bool
}

type DataTransformer interface {
    Transform(data []byte) ([]byte, error)
}

type DataRepository interface {
    Save(data []byte) error
}

// Чисто функциональный интерфейс
type ProcessDataDeps struct {
    Validator    DataValidator
    Transformer  DataTransformer
    Repository   DataRepository
}

func ProcessData(deps ProcessDataDeps, data []byte) error {
    if !deps.Validator.Validate(data) {
        return errors.New("invalid data")
    }
    
    transformed, err := deps.Transformer.Transform(data)
    if err != nil {
        return err
    }
    
    return deps.Repository.Save(transformed)
}

// Функциональная композиция
type DataProcessor func([]byte) error

func CreateDataProcessor(deps ProcessDataDeps) DataProcessor {
    return func(data []byte) error {
        return ProcessData(deps, data)
    }
}
```

## Пример 3: Кэширующий сервис

### До (объектно-императивный стиль)
```go
type CacheService struct {
    store map[string]interface{}
    mu    sync.RWMutex
    ttl   time.Duration
}

func (c *CacheService) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    value, exists := c.store[key]
    return value, exists
}

func (c *CacheService) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.store[key] = value
}
```

### После (функциональный интерфейс с IoC)
```go
// Определяем интерфейсы для зависимостей
type CacheStore interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{})
}

// Чисто функциональный интерфейс
type CacheDependencies struct {
    Store CacheStore
}

func GetFromCache(deps CacheDependencies, key string) (interface{}, bool) {
    return deps.Store.Get(key)
}

func SetToCache(deps CacheDependencies, key string, value interface{}) {
    deps.Store.Set(key, value)
}

// Функциональные обертки с каррированием
type CacheGetter func(key string) (interface{}, bool)
type CacheSetter func(key string, value interface{})

func CreateCacheGetter(deps CacheDependencies) CacheGetter {
    return func(key string) (interface{}, bool) {
        return GetFromCache(deps, key)
    }
}

func CreateCacheSetter(deps CacheDependencies) CacheSetter {
    return func(key string, value interface{}) {
        SetToCache(deps, key, value)
    }
}
```

## Пример 4: HTTP-обработчик

### До (объектно-императивный стиль)
```go
type UserHandler struct {
    userService *UserService
    renderer    *TemplateRenderer
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    user, err := h.userService.GetUser(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if err := h.renderer.Render(w, "user_template", user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

### После (функциональный интерфейс с IoC)
```go
// Определяем интерфейсы для зависимостей
type UserService interface {
    GetUser(id string) (*User, error)
}

type TemplateRenderer interface {
    Render(w io.Writer, template string, data interface{}) error
}

// Чисто функциональный интерфейс
type UserHandlerDeps struct {
    UserService UserService
    Renderer    TemplateRenderer
}

func HandleGetUser(deps UserHandlerDeps, w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    user, err := deps.UserService.GetUser(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    if err := deps.Renderer.Render(w, "user_template", user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Фабрика для создания обработчиков
func CreateUserHandler(deps UserHandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        HandleGetUser(deps, w, r)
    }
}
```