# Логика Хоара для программистов-2

## Пример 1

До:
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

После:
- ослабили предусловние: оставили "любые" пути на файлы
- усилили постусловие: добавили ошибку при проверке путей к файлам
~~~go
// Предусловия: в качестве массива строк передаются пути на файлы
// Постусловия: если пути к файлам некорректные, то ошибка; если в файлах были найдены переменные окружения, то они были установлены
func Load(filenames ...string) (err error) {
    // усиление постусловия
    if err := validatePaths(filenames); err != nil{
        return err
    }

	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}
~~~

## Пример 2

До:
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

После:
- ослабили предусловние: убрали его совсем
- усилили постусловие: если не было установлено переменных окружений, то вернули ошибку
~~~go
// Постусловия: если нет переменных окружений, то завершаем работу с ошибкой; получен URL веб-хука или завершение работы программы
func getWebhookURL(tgLogger *slog.Logger) (string, error) {
    // усиление постусловия
    if !isEnvSetUp{
         return "", fmt.Errorf("environment not configured")
    }

	webhookURL := ""

	switch envFlag {
	case "prod":
		webhookURL = os.Getenv("MATTERMOST_WEBHOOK_URL")
	case "debug":
		webhookURL = os.Getenv("MATTERMOST_WEBHOOK_URL_DEBUG")
	}

	if webhookURL == "" {
		return "", fmt.Errorf("webhook url is empty for env=%s", envFlag)
	}

	return webhookURL, nil
}
~~~

## Пример 3

До:
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

После:
- ослабили предусловние: убрали предусловие
- усилили постусловие: добавили классификацию ошибок и гарантию, что вернется не пустой puzzleID
~~~go
// Постусловия: классифицированная ошибка, если есть проблемы; получена моделька пазла с непустым ID, если все хорошо
func DailyPuzzle() (models.Puzzle, error) {
	resp, err := http.Get(apiPuzzleDaily)
	if err != nil {
		return models.Puzzle{}, fmt.Errorf("network: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Puzzle{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Puzzle{}, fmt.Errorf("read body: %w", err)
	}

	var puzzle models.Puzzle
	err = json.Unmarshal(body, &puzzle)
	if err != nil {
		return models.Puzzle{}, fmt.Errorf("parse: %w", err)
	}

    // усиление постусловия
    if puzzle.ID == "" {
        return models.Puzzle{}, fmt.Errorf(
            "parse: puzzle has empty ID",
        )
    }

	return puzzle, nil
}
~~~

## Пример 4

До:
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

После:
- ослабили предусловние: pgn -> строка в произвольном формате
- усилили постусловие: возвращаем ошибку валидации, если некорректная pgn-строка
~~~go
// Предуловия: pgn — строка в производном формате
// Постусловия: при некорректной pgn-строке — ошибка валидации; валидный URL-адрес с корректным расположением доски и ошибка либо ее отсутствие 
func GetPictureURL(pgn string) (string, error) {
    // усиление постусловия
    if err := validate(pgn); err != nil{
        return "", err
    }

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

## Пример 5

До:
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

После:
- ослабили предусловние: убрали "корректность" передаваемого URL
- усилили постусловие: добавили валидацию URL и стали возвращать ошибку
~~~go
// Предусловия: передан корректный игровой ID и произвольный URL
// Постусловия: вернули ошибку, если URL некорректен; сформирован payload, преобразуемый в формат отправки сообщения
func ComposePayload(gameID, gamePicURL string) (map[string]interface{}, error) {
	// усиление постусловия
    if err := urlValid(gamePicURL); err != nil{
        return nil, err
    }

    if gameID == "" {
		return nil, errors.New("empty game ID")
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
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(plEncoded, &result)
	if err != nil {
		return nil
	}

	return result, nil
}
~~~