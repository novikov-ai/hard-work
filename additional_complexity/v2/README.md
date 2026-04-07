# Дополнительная сложность -- мать всех запашков кода (2)

## Пример 1

```go
func toMap(v any) map[string]interface{} {
    encoded, err := json.Marshal(v)
    if err != nil {
        return nil
    }
    var result map[string]interface{}
    if err = json.Unmarshal(encoded, &result); err != nil {
        return nil
    }
    return result
}

func ComposePayload(gameID, gamePicURL string) map[string]interface{} {
    // ...
    return toMap(pl)
}

func ComposePayloadV2() map[string]interface{} {
    // ...
    return toMap(pl)
}
```

Я могу чётко сформулировать, никуда больше не перескакивая, что код приводит к мапе любой аргумент, если удалось его анмаршалить, или nil.

Функция действительно должна занимать 8 строк кода и 2 условия. 

## Пример 2

```go
// internal/usecases/config/env.go — единая точка
func WebhookURL(env string) string {
    switch env {
    case "prod":
        return os.Getenv("MATTERMOST_WEBHOOK_URL")
    case "debug":
        return os.Getenv("MATTERMOST_WEBHOOK_URL_DEBUG")
    }
    return ""
}

// cmd/main.go и cmd/reminder/main.go — оба используют одну функцию
webhookURL := config.WebhookURL(envFlag)
if webhookURL == "" {
    log.Fatal("Webhook url is empty")
}
```

Я могу чётко сформулировать, никуда больше не перескакивая, что код формирует webhook url в зависимости от переданных значений env.

Функция действительно должна занимать 6 строк кода и 3 условия. 

## Пример 3

```go
type Env string

const (
    EnvProd  Env = "prod"
    EnvDebug Env = "debug"
)

func (e Env) IsValid() bool {
    return e == EnvProd || e == EnvDebug
}

var envFlag string
flag.StringVar(&envFlag, "e", "debug", "environment (prod/dev)")

env := Env(envFlag)
if !env.IsValid() {
    log.Fatalf("unknown environment: %q", envFlag)
}
```

Я могу чётко сформулировать, никуда больше не перескакивая, что код формирует проверяет валидность переменных окружения.

Функция действительно должна состоять из двух условий.

## Пример 4

```go
func fetchPuzzle() (models.Puzzle, error) {
	resp, err := http.Get(apiPuzzleDaily)
	if err != nil {
		return models.Puzzle{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Puzzle{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Puzzle{}, err
	}

	var puzzle models.Puzzle
	err = json.Unmarshal(body, &puzzle)
	if err != nil {
		return models.Puzzle{}, err
	}

	return puzzle, nil
}
```

Я могу чётко сформулировать, никуда больше не перескакивая, что код делает сетевой запрос и обрабатывает результат.

Функция действительно должна занимать 22 строки кода и 4 условия. 

## Пример 5

```go
func main() {
    flag.StringVar(&envFlag, "e", "debug", "environment (prod/dev)")
    flag.Parse()

    cfg := mustLoadConfig()
    if err := run(cfg); err != nil {
        cfg.logger.Error(err.Error())
        log.Fatal(err)
    }
}

func run(cfg appConfig) error {
    puzzle, err := fetchPuzzle(cfg.logger)
    if err != nil {
        return err
    }
    payload, err := buildPayload(puzzle, cfg.logger)
    if err != nil {
        return err
    }
    return sendWebhook(cfg.webhookURL, payload, cfg.logger)
}
```

Я могу чётко сформулировать, никуда больше не перескакивая, что код запускает приложение с заданным конфигом.

Функция действительно должна занимать 8 строк кода и 2 условия. 


