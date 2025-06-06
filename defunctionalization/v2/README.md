## Дефункционализация + CPS

 Комбинация CPS и дефункционализации — мощный инструмент для построения устойчивых асинхронных систем, который легко может найти применение в моей работе. Рассмотрим типичные сценарии в Классифайде. 

### 1. **Многошаговая публикация объявления**  
**Проблема:**  
Процесс публикации включает:  
1. Валидацию данных → 2. Загрузку медиа → 3. Модерацию → 4. Опубликование  
При сбое на любом этапе нужно восстановить процесс.

**Решение CPS + дефункционализация:**  
```go
// Продолжения как структуры
type Continuation struct {
    Action string
    AdID   string
    Data   map[string]interface{}
}

// Обработчик в CPS-стиле
func PublishAd(ad Advertisement, next Continuation) {
    // Дефункционализация: сохраняем продолжение в БД
    saveState(Continuation{
        Action: "validate",
        AdID:   ad.ID,
        Data:   map[string]interface{}{"step": 1},
    })
    ValidateAd(ad, validateContinuation(ad))
}

func validateContinuation(ad Advertisement) Continuation {
    return Continuation{
        Action: "upload_media",
        AdID:   ad.ID,
        Data:   map[string]interface{}{"step": 2},
    }
}

// Восстановление после сбоя
func Resume(adID string) {
    state := loadState(adID) // Загружаем из БД
    switch state.Action {
    case "validate":
        ValidateAd(loadAd(adID), validateContinuation())
    case "upload_media":
        UploadMedia(adID, mediaContinuation())
    // ...
    }
}
```

**Преимущества:**  
- Состояние процесса сохраняется в БД при каждом шаге  
- После сбоя сервера можно продолжить с прерванного места  
- Нет потери данных при обновлении системы  

### 2. **Асинхронная обработка изображений**  
**Проблема:**  
После загрузки фото нужно:  
1. Сжать → 2. Сгенерировать превью → 3. Распознать текст → 4. Сохранить в CDN  

**Решение:**  
```go
// Дефункционализированные продолжения
type ImageTask struct {
    Steps []ImageOp // [Resize(800), Thumbnail(200), OCR(), ...]
}

// Интерпретатор операций
func ProcessImage(img []byte, task ImageTask) {
    for _, op := range task.Steps {
        img = applyOp(img, op) // Применяем операцию
        saveTaskProgress(task) // Сохраняем прогресс
    }
}

// Пример использования
task := ImageTask{
    Steps: []ImageOp{
        {Type: "resize", Width: 800},
        {Type: "thumbnail", Size: 200},
        {Type: "ocr"},
    },
}
ProcessImage(rawImage, task)
```

**Преимущества:**  
- Задачи можно ставить в очередь (RabbitMQ/Kafka)  
- Возможность перезапуска с любого этапа  
- Добавление новых операций без изменения ядра  


### 3. **Цепочки рекомендаций**  
**Проблема:**  
После просмотра объявления:  
1. Обновить историю → 2. Пересчитать рекомендации → 3. Отправить уведомление  

**CPS + дефункционализация:**  
```go
type RecommendationFlow struct {
    UserID string
    Steps  []RecStep
}

func (f *RecommendationFlow) Next() {
    if len(f.Steps) == 0 { return }
    step := f.Steps[0]
    f.Steps = f.Steps[1:]
    saveFlow(f) // Сохраняем состояние
    
    switch step.Type {
    case "update_history":
        updateHistory(f.UserID, step.ItemID, f.Next)
    case "recalculate":
        recalculate(f.UserID, f.Next)
    case "notify":
        sendNotification(f.UserID, f.Next)
    }
}

// Инициализация
flow := RecommendationFlow{
    UserID: "123",
    Steps: []RecStep{
        {Type: "update_history", ItemID: "item789"},
        {Type: "recalculate"},
        {Type: "notify", Template: "new_recommendations"},
    },
}
flow.Next()
```

### Польза:
1. **Отказоустойчивость**  
   Состояние процесса хранится в БД → перезапуск сервера не прерывает операции  

2. **Масштабируемость**  
   Дефункционализированные задачи можно передавать между сервисами через очередь сообщений  

3. **Аналитика процессов**  
   Сохранённые состояния позволяют:  
   - Отлаживать сложные цепочки  
   - Считать метрики выполнения этапов  
   - Визуализировать workflow  

4. **Безопасность**  
   Нет передачи/сериализации исполняемого кода → защита от инъекций  