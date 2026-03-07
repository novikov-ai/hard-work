# Логика Хоара для программистов-2 (2)

##  Ошибка модульного рассуждения — Тип 1 (завязка на передачу аргументов более широких, чем предполагает спецификация)

### Пример 1

`json.Unmarshal` проверяет json-тэги, но несмотря на несовпадение регистра полей все равно их заполняет, что ведет к угрозам безопасности и нежелательному поведению.

Свежий пример:
~~~go
// modelcontextprotocol/go-sdk — структура запроса JSON-RPC

type JSONRPCRequest struct {
    JSONRPC string `json:"jsonrpc"`
    ID      int    `json:"id"`
    Method  string `json:"method"` // "method" (lowercase)
}

// Пример атаки:
malicious := `{
    "jsonrpc": "2.0",
    "id": 1, 
    "Method": "tools/call"
}`
// "Method" (uppercase)

var req JSONRPCRequest
json.Unmarshal([]byte(malicious), &req)

fmt.Println(req.Method)

// WAF/прокси проверяет поле "method" (lowercase) и не видит угрозы
// Go-сервис видит Method как method и выполняет вызов
~~~

### Пример 2

В Go пакет net/http автоматически определяет Content-Type, если хендлер его не выставил.

~~~go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Write([]byte("<html><body>Hello</body></html>"))
}
~~~

И до Go 1.11 это работало, но с Go 1.11 перестало, так как перестал проставляться Content-Type.


##  Ошибка модульного рассуждения — Тип 2 (завязка на результат, который не гарантирован спецификацией)

### Пример 1

~~~go
type Task struct {
    Priority int
    Name     string
}

tasks := []Task{
    {1, "deploy"},
    {1, "test"},
    {1, "lint"},
    {2, "build"},
}

sort.Slice(tasks, func(i, j int) bool {
    return tasks[i].Priority < tasks[j].Priority
})
~~~

### Пример 2

~~~go
import (
    "github.com/prometheus/prometheus/tsdb"
    "github.com/prometheus/prometheus/promql"
    "github.com/prometheus/prometheus/model/labels"
)

func queryMetrics() {
    eng := promql.NewEngine(promql.EngineOpts{
        MaxSamples: 50000000,
        Timeout:    2 * time.Minute,
    })
}
~~~