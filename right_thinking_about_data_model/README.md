# Как правильно думать над моделью данных

## 1. Избыточное хранение названий категорий в объявлениях

При получении списка объявлений клиентам нужны названия категорий. Обычно это реализуется через JOIN с таблицей категорий, что увеличивает время ответа и усложняет API (пользователи вынуждены делать отдельные запросы для получения категорий).

**Решение**:  
Добавить избыточное поле `category_name` в таблицу объявлений.  

- При создании/обновлении объявления копировать название категории из справочника категорий в поле `ad.category_name`.  
- При изменении названия категории триггером или фоновым джобом обновить все связанные объявления.

```go
// Структура объявления с избыточным полем
type Ad struct {
    ID           string
    Title        string
    CategoryID   string
    CategoryName string // Избыточное поле
}

// Метод создания объявления
func (s *AdService) CreateAd(ad Ad) error {
    // Получаем актуальное название категории
    category, err := s.categoryRepo.GetByID(ad.CategoryID)
    if err != nil {
        return err
    }
    ad.CategoryName = category.Name

    // Сохраняем объявление в БД
    return s.adRepo.Save(ad)
}
```

// Триггер в БД для автоматического обновления CategoryName при изменении категории
```sql
CREATE TRIGGER update_ad_category_name 
AFTER UPDATE ON categories 
FOR EACH ROW 
BEGIN
    UPDATE ads 
    SET category_name = NEW.name 
    WHERE category_id = NEW.id;
END;
```

**Упрощение API**:  
- Клиенты получают название категории сразу в ответе на запрос объявлений, без дополнительных JOIN или запросов.  
- Пример ответа API:  
  ```json
  {
      "id": "123",
      "title": "Велосипед",
      "category_name": "Спорт и отдых" // вместо "category_id": "789"
  }
  ```

## 2. Инвертированный индекс для поиска по ключевым словам

Поиск объявлений по ключевым словам (`LIKE '%велосипед%'`) работает медленно на больших данных.

**Решение**:  
Создать инвертированный индекс (мапу) в Redis или Elasticsearch, где ключ — слово, значение — список ID объявлений.  

- При добавлении/обновлении объявления парсить текст на ключевые слова и обновлять индекс.  
- Использовать асинхронную обработку, чтобы не блокировать основной поток.  

```go
// Инвертированный индекс в Redis
type SearchIndex struct {
    redisClient *redis.Client
}

// Обновление индекса при сохранении объявления
func (si *SearchIndex) UpdateIndex(ad Ad) {
    words := extractKeywords(ad.Title + " " + ad.Description)
    for _, word := range words {
        si.redisClient.SAdd("index:"+word, ad.ID)
    }
}

// Поиск объявлений по слову
func (si *SearchIndex) Search(word string) ([]string, error) {
    return si.redisClient.SMembers("index:" + word).Result()
}

// Асинхронный обработчик
func (s *AdService) CreateAd(ad Ad) error {
    // Сохранение в БД
    if err := s.adRepo.Save(ad); err != nil {
        return err
    }
    // Асинхронное обновление индекса
    go s.searchIndex.UpdateIndex(ad)
    return nil
}
```

**Упрощение API**:  
- Клиенты получают мгновенный поиск по ключевым словам через простой эндпоинт:  
  `GET /ads/search?query=велосипед` → возвращает список ID объявлений.  
- Нет необходимости в сложных SQL-запросах или полнотекстовом поиске на стороне клиента.

---

## 3. Кэширование геолокационных данных для «Объявлений рядом»

Поиск объявлений в радиусе 5 км (`ST_Distance` в PostgreSQL) требует вычислений и медленный на больших объемах.

**Решение**:  
Добавить избыточные поля `geohash` и `latitude/longitude` с пространственным индексом (например, PostGIS).  

- При создании объявления вычислять geohash и сохранять его в отдельное поле.  
- Использовать GiST-индекс для ускорения пространственных запросов.  

```go
type Ad struct {
    ID        string
    Latitude  float64
    Longitude float64
    GeoHash   string // Избыточное поле для быстрого поиска
}

// Обновление geohash при сохранении
func (s *AdService) CreateAd(ad Ad) error {
    ad.GeoHash = geohash.Encode(ad.Latitude, ad.Longitude)
    return s.adRepo.Save(ad)
}
~~~

// Поиск "объявлений рядом" через GiST-индекс
// SQL-запрос с использованием PostGIS
~~~sql
SELECT *
FROM ads
WHERE ST_DWithin(
    ST_MakePoint(ads.longitude, ads.latitude),
    ST_MakePoint(?, ?),
    5000 // 5 км
);
```

**Упрощение API**:  
- Клиенты получают эндпоинт с быстрым ответом:  
  `GET /ads/nearby?lat=55.75&lon=37.62&radius=5000`  
- Время ответа сокращается с 2 сек до 200 мс благодаря индексу.