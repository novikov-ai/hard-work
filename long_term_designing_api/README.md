# Долгосрочное проектирование API

## **Кейс 1: Сериализация данных (JSON/Protobuf)**  
**Проблема обратной совместимости**  
Клиенты неявно полагаются на порядок полей, формат значений или отсутствие дополнительных ключей в сериализованных данных, хотя контракт гарантирует только семантику полей.

**Применение приёма**  
```go
// testutil/chaos_serializer.go
type ChaosSerializer struct{ base serializer.Serializer }

func (s *ChaosSerializer) Marshal(v any) ([]byte, error) {
    if chaos.Enabled() {
        // 1. Динамически меняем порядок полей
        // 2. Добавляем случайные ключи
        // 3. Меняем типы значений (числа → строки)
        data := map[string]any{
            "__chaos__":      true,
            "injected_field": rand.String(10),
            "value":         fmt.Sprintf("chaos:%v", reflect.ValueOf(v).Field(0)),
        }
        return json.Marshal(data)
    }
    return s.base.Marshal(v)
}
```

**Что обнаруживает**:
- Клиенты, использующие regex-парсинг вместо стандартных кодеков
- Жёсткую зависимость от порядка полей
- Ожидания определённых типов значений
- Уязвимости к битым UTF-8 и инъекциям

**Обобщение для API**:  
> "Формат сериализации — изменяемая деталь реализации. Клиенты должны использовать только официальные десериализаторы, обрабатывающие любую валидную перестановку полей."


## **Кейс 2: Идентификаторы ресурсов**  
**Проблема обратной совместимости**  
Клиенты предполагают фиксированную структуру ID (число/UUID), хотя контракт определяет их как непрозрачные строки.

**Применение приёма**  
```go
// testutil/chaos_id_generator.go
type ChaosIDGenerator struct{}

func (g *ChaosIDGenerator) NewID() string {
    switch rand.Intn(4) {
    case 0:  // Числовой формат
        return strconv.Itoa(rand.Intn(100000))
    case 1:  // UUID с префиксом
        return "id_" + uuid.NewString()
    case 2:  // Юникод-строки
        return "标识_" + rand.UnicodeString(8)
    default: // Base64
        return base64.RawURLEncoding.EncodeToString(rand.Bytes(16))
    }
}
```

**Что обнаруживает**:
- Конвертацию ID в int без проверок
- Парсинг через strings.Split/Regexp
- Ограничение длины буферов
- Неподдержку Unicode

**Обобщение для API**:  
> "Идентификаторы — абстрактные токены. Любая попытка их интерпретации нарушает контракт и приведёт к несовместимости при смене формата генерации."

## **Кейс 3: Конкурентное выполнение операций**  
**Проблема обратной совместимости**  
Клиенты полагаются на недокументированные гарантии: порядок обработки, время выполнения, отсутствие ошибок.

**Применение приёма**  
```go
// testutil/chaos_executor.go
type ChaosExecutor struct {
    base     TaskExecutor
    latency  time.Duration
    errRate  float64
}

func (e *ChaosExecutor) Execute(ctx context.Context, task Task) error {
    if chaos.Enabled() {
        // 1. Случайные задержки
        delay := time.Duration(rand.Float64() * float64(e.latency))
        time.Sleep(delay)
        
        // 2. Инъекция ошибок
        if rand.Float64() < e.errRate {
            return ChaosError("injected failure")
        }
    }
    return e.base.Execute(ctx, task)
}
```

**Что обнаруживает**:
- Предположение FIFO-обработки
- Отсутствие таймаутов и retry-логики
- Игнорирование ошибок выполнения
- Неподготовленность к сетевой нестабильности

**Обобщение для API**:  
> "Порядок и время выполнения операций не гарантируются. Клиенты должны проектироваться для работы в условиях частичных отказов и произвольных задержек."