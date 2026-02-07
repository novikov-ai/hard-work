# Логика Хоара для программистов

## Пример 1

~~~go
// Предусловия: в качестве массива строк передаются корректные пути на файлы, которые существуют в системе
// Постусловия: если в файлах были найдены переменные окружения, то они были установлены или ошибка при наличии проблем
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return 
		}
	}
	return
}
~~~

Место вызова функции:
~~~go
func main() {
	flag.StringVar(&envFlag, "e", "debug", "environment (prod/dev)")
	flag.Parse()

	// Предусловия: в качестве массива строк передаются корректные пути на файлы, которые существуют в системе
	// Постусловия: если в файлах были найдены переменные окружения, то они были установлены или ошибка при наличии проблем
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't find the config file")
	}

	tgLogger, err := telegram.NewLogger()
	if err != nil {
		log.Println("Can't create telegram logger:", err)
	}

	// Предусловия: установлены переменные окружения для вебхуков
	// Постусловия: получен URL веб-хука или завершение работы программы
	webhookURL := getWebhookURL(tgLogger)

    // ... 
}
~~~

Происходит стыковка спецификаций: постусловия для `Load` являются предусловиями для `getWebhookURL`.

## Пример 2

~~~go
// Предусловия: установлены переменные окружения для вебхуков
// Постусловия: получен URL веб-хука или завершение работы программы
func getWebhookURL(tgLogger *slog.Logger) string {
	webhookURL := ""

	switch envFlag {
	case "prod":
		webhookURL = os.Getenv("MATTERMOST_WEBHOOK_URL")
	case "debug":
		webhookURL = os.Getenv("MATTERMOST_WEBHOOK_URL_DEBUG")
	}

	if webhookURL == "" {
		telegram.LogError(tgLogger, "Webhook url is empty", nil)
		log.Fatal("Webhook url is empty")
	}

	return webhookURL
}
~~~

Происходит стыковка спецификаций: постусловия для `Load` являются предусловиями для `getWebhookURL` (см. прошлый пример).

## Пример 3

~~~go
// Предусловия: возможность установить http-соединение
// Постусловия: получена моделька пазла или ошибка, если были проблемы с формированием пазла
func DailyPuzzle() (models.Puzzle, error) {
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
~~~

Место вызова функции:
~~~go
func main(){
    // Предуловия: возможность установить http-соединение
    // Постусловия: получена валидная моделька пазла или ошибка, если были проблемы с формированием пазла
	puzzle, err := fetching.DailyPuzzle()
	if err != nil {
		telegram.LogError(tgLogger, "Can't fetch daily puzzle", err)
		log.Fatal("Can't fetch daily puzzle:", err)
	}

	// Предусловия: pgn — строка, представляющую собой формат файла для сохранения шахматных партий: Portable Game Notation
	// Постусловия: валидный URL-адрес с корректным расположением доски и ошибка либо ее отсутствие
	picURL, err := pgn.GetPictureURL(puzzle.Game.Pgn)
	if err != nil {
		telegram.LogError(tgLogger, "Can't get picture from PGN", err)
		log.Fatal("Can't get picture from PGN:", err)
	}
    // ...
}
~~~

Стыковка спецификаций происходит между `DailyPuzzle` и `GetPictureURL`: постусловия первой функции становятся предусловиями второй.

## Пример 4

~~~go
// Предуловия: pgn — строка, представляющую собой формат файла для сохранения шахматных партий: Portable Game Notation
// Постусловия: валидный URL-адрес с корректным расположением доски и ошибка либо ее отсутствие 
func GetPictureURL(pgn string) (string, error) {
	respHTML, err := pgnImportRetrieveHTML(pgn)
	if err != nil {
		return "", err
	}

	picURL := getPositionURL(respHTML)
	if picURL == "" {
		return "", errors.New("picture not found")
	}

	flipBoard, err := blackMove(pgn)
	if err != nil {
		return "", err
	}

	if flipBoard {
		picURL = picUrlWithColorBlack(picURL)
	}

	return picURL, nil
}
~~~

Место вызова функции:
~~~go
func main(){
    // Предуловия: возможность установить http-соединение
    // Постусловия: получена валидная моделька пазла или ошибка, если были проблемы с формированием пазла
	puzzle, err := fetching.DailyPuzzle()
	if err != nil {
		telegram.LogError(tgLogger, "Can't fetch daily puzzle", err)
		log.Fatal("Can't fetch daily puzzle:", err)
	}

	// Предуловия: pgn — строка, представляющую собой формат файла для сохранения шахматных партий: Portable Game Notation
	// Постусловия: валидный URL-адрес с корректным расположением доски и ошибка либо ее отсутствие
	picURL, err := pgn.GetPictureURL(puzzle.Game.Pgn)
	if err != nil {
		telegram.LogError(tgLogger, "Can't get picture from PGN", err)
		log.Fatal("Can't get picture from PGN:", err)
	}
    // ...
}
~~~

Стыковка спецификаций происходит между `DailyPuzzle` и `GetPictureURL`: постусловия первой функции становятся предусловиями второй (см. предыдущий пример).

## Пример 5

~~~go
// Предусловия: переданы корректные игровой ID и URL, которые не пустые и существуют в системе
// Постусловия: сформирован payload, преобразуемый в формат отправки сообщения
func ComposePayload(gameID, gamePicURL string) map[string]interface{} {
	if gameID == "" || gamePicURL == "" {
		return nil
	}

	gameURL := endpointPuzzleTraining + gameID

	pl := models.Payload{
		Username: Username,
		Text:     fmt.Sprintf(Message, gameURL),
		IconURL:  IconURL,
		Attachments: []map[string]interface{}{
			{
				"image_url": gamePicURL,
			},
		},
	}

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
}
~~~

Место вызова функции:
~~~go
func main(){
    // ...

	// Предуловия: pgn — строка, представляющую собой формат файла для сохранения шахматных партий: Portable Game Notation
	// Постусловия: валидный URL-адрес с корректным расположением доски и ошибка либо ее отсутствие
	picURL, err := pgn.GetPictureURL(puzzle.Game.Pgn)
	if err != nil {
		telegram.LogError(tgLogger, "Can't get picture from PGN", err)
		log.Fatal("Can't get picture from PGN:", err)
	}

	// Предусловия: переданы корректные игровой ID и URL, которые не пустые и существуют в системе
	// Постусловия: сформирован payload, преобразуемый в формат отправки сообщения
	payload := presentation.ComposePayload(puzzle.Puzzle.ID, picURL)
	if payload == nil {
		telegram.LogError(tgLogger, "Error composing payload", err)
		log.Fatal("Error composing payload:", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		telegram.LogError(tgLogger, "Error creating JSON payload", err)
		log.Fatal("Error creating JSON payload:", err)
	}
    // ...
}
~~~

Стыковка спецификаций происходит между `GetPictureURL` и `ComposePayload`: постусловия первой функции становятся предусловиями второй.
