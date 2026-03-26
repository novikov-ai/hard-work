# Дополнительная сложность -- мать всех запашков кода

## Пример 1 — Дублирование кода (`compose.go`)

`ComposePayload` и `ComposePayloadV2` содержат идентичный блок marshal → unmarshal.

Было:
```go
// ComposePayload
plEncoded, err := json.Marshal(pl)
if err != nil {
    return nil
}
var result map[string]interface{}
err = json.Unmarshal(plEncoded, &result)
if err != nil {
    return nil
}
return result

// ComposePayloadV2 — слово в слово то же самое
plEncoded, err := json.Marshal(pl)
if err != nil {
    return nil
}
var result map[string]interface{}
err = json.Unmarshal(plEncoded, &result)
if err != nil {
    return nil
}
return result
```

Стало:
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

---

## Пример 2 — Switch Statement (повторяется в двух местах)

`cmd/main.go` уже вынесен в `getWebhookURL()`, но `cmd/reminder/main.go` дублирует тот же switch вручную.

Было:
```go
// cmd/reminder/main.go — switch продублирован inline
webhookURL := ""
switch envFlag {
case "prod":
    webhookURL = os.Getenv("MATTERMOST_WEBHOOK_URL")
case "debug":
    webhookURL = os.Getenv("MATTERMOST_WEBHOOK_URL_DEBUG")
}
if webhookURL == "" {
    log.Fatal("Webhook url is empty")
}
```

Стало:
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

---

## Пример 3 — Primitive Obsession (`envFlag`)

`envFlag` — строка со скрытым enum-контрактом: допустимы только `"prod"` и `"debug"`. Опечатка (`"Prod"`, `"production"`) не даст ошибки компиляции.

Было:
```go
// cmd/main.go:19, cmd/reminder/main.go:16
var envFlag string

flag.StringVar(&envFlag, "e", "debug", "environment (prod/dev)")

switch envFlag {
case "prod":
    ...
case "debug":
    ...
}
```

Стало:
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

---

## Пример 4 — Comments as a Smell (`cmd/main.go`, `fetch.go`, `pgn.go`)

Комментарии пересказывают то, что и так читается из кода. Особенно заметно в `cmd/main.go`.

Было:
```go
// Предуловия: нет
// Постусловия: получена моделька пазла или ошибка, если были проблемы с формированием пазла
puzzle, err := fetching.DailyPuzzle()
if err != nil {
    log.Fatal("Can't fetch daily puzzle:", err)
}

// Предуловия: pgn — строка, представляющую собой формат файла для сохранения шахматных партий: Portable Game Notation
// Постусловия: валидный URL-адрес с корректным расположением доски и ошибка либо ее отсутствие
picURL, err := pgn.GetPictureURL(puzzle.Game.Pgn)
if err != nil {
    log.Fatal("Can't get picture from PGN:", err)
}
```

Стало:
```go
puzzle, err := fetching.DailyPuzzle()
if err != nil {
    log.Fatal("Can't fetch daily puzzle:", err)
}

picURL, err := pgn.GetPictureURL(puzzle.Game.Pgn)
if err != nil {
    log.Fatal("Can't get picture from PGN:", err)
}
```

---

## Пример 5 — Long Method (`cmd/main.go:21-82`, 61 строка)

`main()` делает всё подряд: загрузка конфига, логгер, webhook URL, fetch пазла, PGN → картинка, payload, JSON, HTTP POST.

Было:
```go
func main() {
    flag.StringVar(&envFlag, "e", "debug", "environment (prod/dev)")
    flag.Parse()

    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Can't find the config file")
    }
    tgLogger, err := telegram.NewLogger()
    if err != nil {
        log.Println("Can't create telegram logger:", err)
    }
    webhookURL := getWebhookURL(tgLogger)
    log.Println("Start fetching a new daily puzzle...")
    puzzle, err := fetching.DailyPuzzle()
    if err != nil {
        telegram.LogError(tgLogger, "Can't fetch daily puzzle", err)
        log.Fatal("Can't fetch daily puzzle:", err)
    }
    picURL, err := pgn.GetPictureURL(puzzle.Game.Pgn)
    if err != nil {
        telegram.LogError(tgLogger, "Can't get picture from PGN", err)
        log.Fatal("Can't get picture from PGN:", err)
    }
    payload := presentation.ComposePayload(puzzle.Puzzle.ID, picURL)
    if payload == nil {
        log.Fatal("Error composing payload:", err)
    }
    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        log.Fatal("Error creating JSON payload:", err)
    }
    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadJSON))
    if err != nil {
        log.Fatal("Error sending webhook request:", err)
    }
    defer resp.Body.Close()
    log.Println("Puzzle was sent successfully!")
}
```

Стало:
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
