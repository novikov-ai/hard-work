# Про раздутость кода

### Пример 1

~~~go
func resolvedField(field Field) Field {
	// ...
	for _, constraint := range field.Constraints {
		existingConstraint, ok := uniqueConstraints[constraint.GetCode()]
		if !ok {
			uniqueConstraints[constraint.GetCode()] = constraint
			continue
		}

		// "Неименованный" кусок кода:
        // Происходит проверка на значимость одного Constraint над другим в зависимости от количества аргументов
		if len(existingConstraint.GetArguments()) > len(constraint.GetArguments()) {
			uniqueConstraints[constraint.GetCode()] = existingConstraint
		} else {
			uniqueConstraints[constraint.GetCode()] = constraint
		}
	}

	// ...
}
~~~

### Пример 2, 3 
~~~go
func (gc *gpsClient) GetTripsByVehicle(ctx context.Context, vehicleID int64, start, end string) ([]models.Trip, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

    // "Неименованный" кусок кода:
    // Инициирование http-клиента вполне семантически автономно
	cl := http.DefaultClient
	respGeo, err := cl.Do(request)
	if err != nil {
		return []models.Trip{}, err
	}
	defer respGeo.Body.Close()

    // "Неименованный" кусок кода:
    // Определение ОК-статуса
	if respGeo.StatusCode == http.StatusOK {
		body, err := io.ReadAll(respGeo.Body)
		if err != nil {
			return []models.Trip{}, err
		}

		var decoder models.GeoDecoder
		err = json.Unmarshal(body, &decoder)
		if err != nil {
			return []models.Trip{}, err
		}

		k := 0
		for i, r := range decoder.Results {
			// 
		}
	}

	return trips, nil
}
~~~

### Пример 4 

~~~go
func Process(){
    // ...

    // "Неименованный" кусок кода:
    // Определение наличия черновика
    if draft != nil && draft.SessionID != "" {
		publishSessionID = draft.SessionID
	}

	if (data.UserData.ItemId != nil && item != nil && item.UserID != nil) && (employee == nil && *item.UserID != user.ID) {
		return nil, errors.New("not found")
	}

	categoryID := p.getTarget(ctx, draft, &data.UserData.Navigation)
	p.expandConfig(
		ctx,
		data,
		draft,
	)

    // ...
}
~~~

### Пример 5 

~~~go
func Exec(){
    // ...
    if err != nil {
        // "Неименованный" кусок кода:
        // Определение типа ошибки
			switch {
			case errors.Is(err, draftStorage.ErrNotFound):
			case errors.Is(err, draftStorage.ErrForbidden):
				return nil, errors.New(ErrDraftLimitExceeded)
			default:
				return nil, errors.New(ErrDraftDefault)
			}
		}
    // ... 
}
~~~