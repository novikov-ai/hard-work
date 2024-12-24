# Что такое абстракция

### Пять точных абстракций - что во что в соответствующих частях кода отображается.

Данные объявления (структура в базе данных) отображаются в поисковый индекс, подходящий для эффективного поиска.

~~~go
// Структура объявления в базе данных
type AdDBModel struct {
    ID          string
    Title       string
    Description string
    CategoryID  int
    CreatedAt   time.Time
}

// Структура для поискового индекса
type AdIndexModel struct {
    ID        string
    Title     string
    Content   string // Объединённые текстовые поля
    Tags      []string
    Timestamp int64
}

// Трансформация: AdDBModel -> AdIndexModel
func TransformAdToIndexModel(ad AdDBModel) AdIndexModel {
    return AdIndexModel{
        ID:        ad.ID,
        Title:     ad.Title,
        Content:   fmt.Sprintf("%s %s", ad.Title, ad.Description),
        Tags:      mapCategoryToTags(ad.CategoryID),
        Timestamp: ad.CreatedAt.Unix(),
    }
}

// Вспомогательная функция
func mapCategoryToTags(categoryID int) []string {
    // Преобразование категории в список тегов
    categories := map[int][]string{
        1: {"electronics", "gadgets"},
        2: {"home", "furniture"},
    }
    return categories[categoryID]
}
~~~

Исторические данные о ценах (промоакции, скидки, базовые тарифы) отображаются в итоговую стоимость, актуальную для конкретного момента.

~~~go
// Входные данные: история изменений цен
type PricingHistory struct {
    BasePrice float64
    Discounts []float64
    Promotions []float64
}

// Итоговая стоимость
type FinalPrice struct {
    Amount float64
}

// Трансформация: PricingHistory -> FinalPrice
func CalculateFinalPrice(history PricingHistory) FinalPrice {
    price := history.BasePrice

    // Применяем скидки
    for _, discount := range history.Discounts {
        price -= discount
    }

    // Добавляем наценки за промо
    for _, promo := range history.Promotions {
        price += promo
    }

    // Гарантируем, что цена не уйдёт в отрицательное значение
    if price < 0 {
        price = 0
    }

    return FinalPrice{Amount: price}
}
~~~

История взаимодействий пользователя с платформой отображается в рекомендации объявлений.

~~~go
// Входные данные: взаимодействия пользователя
type UserInteractions struct {
    ViewedAds   []string
    ClickedAds  []string
    PurchasedAds []string
}

// Рекомендация
type Recommendation struct {
    AdID    string
    RelevanceScore float64
}

// Трансформация: UserInteractions -> []Recommendation
func GenerateRecommendations(interactions UserInteractions) []Recommendation {
    relevanceMap := make(map[string]float64)

    // Увеличиваем релевантность для просмотренных объявлений
    for _, adID := range interactions.ViewedAds {
        relevanceMap[adID] += 1.0
    }

    // Ещё выше для кликнутых
    for _, adID := range interactions.ClickedAds {
        relevanceMap[adID] += 2.0
    }

    // Максимум для купленных
    for _, adID := range interactions.PurchasedAds {
        relevanceMap[adID] += 3.0
    }

    // Преобразуем карту в список рекомендаций
    recommendations := []Recommendation{}
    for adID, score := range relevanceMap {
        recommendations = append(recommendations, Recommendation{AdID: adID, RelevanceScore: score})
    }

    // Сортируем по релевантности
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].RelevanceScore > recommendations[j].RelevanceScore
    })

    return recommendations
}
~~~

Внутренние транзакции сервиса отображаются в запросы к внешней платёжной системе.

~~~go
// Внутренняя транзакция
type InternalTransaction struct {
    UserID string
    Amount float64
    Currency string
}

// Внешний запрос
type PaymentRequest struct {
    Account string
    Total   float64
    Currency string
}

// Трансформация: InternalTransaction -> PaymentRequest
func TransformToPaymentRequest(tx InternalTransaction) PaymentRequest {
    return PaymentRequest{
        Account: tx.UserID,
        Total:   tx.Amount,
        Currency: tx.Currency,
    }
}

// Пример интерфейса платёжной системы
type PaymentGateway interface {
    ProcessPayment(request PaymentRequest) (bool, error)
}
~~~

Лог данных о просмотрах и кликах на объявления отображается в агрегированные метрики.

~~~go
// Лог взаимодействий
type InteractionLog struct {
    AdID   string
    Event  string // "view", "click", "purchase"
    Timestamp time.Time
}

// Агрегированные метрики
type AdMetrics struct {
    AdID    string
    Views   int
    Clicks  int
    Purchases int
}

// Трансформация: []InteractionLog -> map[string]AdMetrics
func AggregateMetrics(logs []InteractionLog) map[string]AdMetrics {
    metrics := make(map[string]AdMetrics)

    for _, log := range logs {
        adMetrics := metrics[log.AdID]

        if log.Event == "view" {
            adMetrics.Views++
        } else if log.Event == "click" {
            adMetrics.Clicks++
        } else if log.Event == "purchase" {
            adMetrics.Purchases++
        }

        metrics[log.AdID] = adMetrics
    }

    return metrics
}
~~~

### Определение Дейкстры применительно к своей практике

Определение Дейкстры понимаю буквально: абстракция помогает нам создавать дополнительный слой, на котором мы мыслим не на уровне деталей реализации, но на уровне смыслов, оперируя которыми мы можем задавать точную спецификацию для своей системы, от которой будет невозможно отклониться.

С помощью такого слоя на практике появляется возможность мыслить не конкретным синтаксисом и деталями реализации, а сразу задавать нужное поведение при помощи точной спецификации, помогающей проектировать и упорядочивать различные части системы.





