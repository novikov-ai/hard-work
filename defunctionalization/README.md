## Дефункционализация

### 1. **Сохранение фильтров поиска**

Пользователи часто сохраняют сложные условия поиска (например, «квартиры в Москве до 10 млн, с ремонтом, от собственника»). Передавать такие фильтры между сервисами (frontend, backend, история поисков) как функции нельзя — их невозможно сериализовать.

Как вариант - определить тип данных, который описывает все возможные условия фильтрации.

```go
// Тип для фильтра объявлений
type Filter struct {
    City       string
    MaxPrice   int
    MinRooms   int
    IsOwner    bool
    HasPhoto   bool
    // ... другие параметры
}

// Метод Apply проверяет, подходит ли объявление под фильтр
func (f Filter) Apply(ad Advertisement) bool {
    return ad.City == f.City &&
           ad.Price <= f.MaxPrice &&
           ad.Rooms >= f.MinRooms &&
           (f.IsOwner && ad.IsOwner) &&
           (f.HasPhoto && len(ad.Photos) > 0)
}
```

### 2. **Асинхронная модерация объявлений**

После отправки объявления на модерацию, система должна:  
1. Проверить текст на запрещенные слова.  
2. Запросить подтверждение телефона, если это новый пользователь.  
3. Уведомить пользователя о результате.  

Если использовать цепочку коллбэков, код станет сложным, а состояние процесса будет неявным.

Можно представить каждое состояние модерации в виде структуры данных.

```go
// Тип для состояния модерации
type ModerationState struct {
    AdID         string
    Step         string // "initial_check", "phone_confirm", "finalize"
    PhoneConfirmToken string
}

// Обработка состояния
func HandleModeration(state ModerationState) {
    switch state.Step {
    case "initial_check":
        if hasBannedWords(state.AdID) {
            notifyUser(state.AdID, "rejected")
        } else {
            token := generatePhoneToken()
            saveState(ModerationState{AdID: state.AdID, Step: "phone_confirm", PhoneConfirmToken: token})
            sendSMS(state.UserPhone, token)
        }
    case "phone_confirm":
        if validateToken(state.PhoneConfirmToken) {
            publishAd(state.AdID)
            notifyUser(state.AdID, "approved")
        }
    // ...
    }
}
```

### 3. **Работа с геоданными**

При поиске объявлений в радиусе 5 км от пользователя, функция-предикат для фильтрации зависит от текущих координат, которые известны только в момент запроса.

Решение - представить геозапрос как структуру с параметрами, которые можно кешировать.

```go
// Гео-фильтр
type GeoFilter struct {
    Lat    float64
    Lon    float64
    Radius int // в метрах
}

// Проверка расстояния
func (g GeoFilter) Apply(ad Advertisement) bool {
    return distance(g.Lat, g.Lon, ad.Lat, ad.Lon) <= g.Radius
}

// Пример использования
filter := GeoFilter{Lat: 55.751244, Lon: 37.618423, Radius: 5000}
results := filterAds(ads, filter)
```