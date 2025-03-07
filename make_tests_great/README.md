# Делаем тесты хорошими

## 1. Тестирование HTTP-обработчика с зависимостью от сервиса
Сценарий: Обработчик /users/{id}, который зависит от UserService.

~~~go
// Мок интерфейса UserService
//go:generate mockgen -destination=mocks/mock_user_service.go -package=mocks . UserService
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
}

func TestUserHandler_GetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserService := mocks.NewMockUserService(ctrl)
    handler := NewUserHandler(mockUserService)

    t.Run("User not found", func(t *testing.T) {
        mockUserService.EXPECT().
            GetUser(gomock.Any(), "invalid_id").
            Return(nil, ErrUserNotFound)

        req := httptest.NewRequest("GET", "/users/invalid_id", nil)
        w := httptest.NewRecorder()

        handler.GetUser(w, req)
        assert.Equal(t, http.StatusNotFound, w.Code) // Проверяем абстрактный эффект: статус 404
    })
}
~~~

Что проверяем:
Абстрактный эффект — корректный HTTP-статус при ошибке.
Явность: Интерфейс UserService явно объявлен, мокируется логика доступа к данным.

## 2. Тестирование потребителя сообщений Kafka
Сценарий: Consumer, обрабатывающий события из Kafka и сохраняющий их в репозиторий.

~~~go
// Мок интерфейса KafkaClient и EventRepository
//go:generate mockgen -destination=mocks/mock_kafka.go -package=mocks . KafkaClient
type KafkaClient interface {
    ReadMessage() (*kafka.Message, error)
    CommitMessage(message *kafka.Message) error
}

func TestEventProcessor_HandleMessage(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockKafka := mocks.NewMockKafkaClient(ctrl)
    mockRepo := mocks.NewMockEventRepository(ctrl)
    processor := NewEventProcessor(mockKafka, mockRepo)

    testMsg := &kafka.Message{Value: []byte(`{"event":"payment_succeeded"}`)}

    t.Run("Success processing", func(t *testing.T) {
        mockKafka.EXPECT().ReadMessage().Return(testMsg, nil)
        mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
        mockKafka.EXPECT().CommitMessage(testMsg).Return(nil)

        err := processor.Handle(context.Background())
        assert.NoError(t, err) // Абстрактный эффект: сообщение обработано и закоммичено
    })
}
~~~

Что проверяем:
Абстрактный эффект — успешная обработка и коммит сообщения.
Явность: Интерфейсы KafkaClient и EventRepository инкапсулируют логику работы с инфраструктурой.

## 3. Тестирование кэширующего прокси для внешнего API
Сценарий: Сервис кэширует ответы от внешнего API.

~~~go
// Мок интерфейсов APIClient и CacheProvider
func TestCachedAPIService_GetData(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockAPI := mocks.NewMockAPIClient(ctrl)
    mockCache := mocks.NewMockCacheProvider(ctrl)
    service := NewCachedAPIService(mockAPI, mockCache)

    t.Run("Cache hit", func(t *testing.T) {
        mockCache.EXPECT().
            Get("cache_key").
            Return([]byte(`cached_data`), nil)

        data, err := service.GetData(context.Background())
        assert.NoError(t, err)
        assert.Equal(t, "cached_data", string(data)) // Проверяем абстрактный эффект: данные из кэша
    })
}
~~~

Что проверяем:
Абстрактный эффект — использование кэша вместо вызова внешнего API.
Явность: Интерфейсы APIClient и CacheProvider декларируют контракты взаимодействия.

## 4. Тестирование межсервисной аутентификации
Сценарий: Middleware для проверки JWT-токена через AuthService.

~~~go
// Мок интерфейса AuthServiceClient (gRPC)
func TestAuthMiddleware_ValidToken(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockAuth := mocks.NewMockAuthServiceClient(ctrl)
    middleware := NewAuthMiddleware(mockAuth)

    t.Run("Valid token", func(t *testing.T) {
        mockAuth.EXPECT().
            ValidateToken(gomock.Any(), &auth.TokenRequest{Token: "valid_token"}).
            Return(&auth.TokenResponse{Valid: true}, nil)

        req := httptest.NewRequest("GET", "/", nil)
        req.Header.Set("Authorization", "Bearer valid_token")
        w := httptest.NewRecorder()

        handler := middleware.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
        handler.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code) // Абстрактный эффект: доступ разрешён
    })
}
~~~

Что проверяем:
Абстрактный эффект — пропуск авторизованного запроса.
Явность: gRPC-интерфейс AuthServiceClient явно определён в protobuf-контракте.

## 5. Тестирование транзакционного репозитория
Сценарий: Транзакция в БД с откатом при ошибке.

~~~go
// Мок интерфейса TxRepository
func TestOrderService_CreateOrderWithRollback(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockTxRepository(ctrl)
    service := NewOrderService(mockRepo)

    t.Run("Rollback on error", func(t *testing.T) {
        mockRepo.EXPECT().BeginTx(gomock.Any()).Return(mockTx, nil)
        mockRepo.EXPECT().SaveOrder(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
        mockTx.EXPECT().Rollback().Return(nil) // Проверяем абстрактный эффект: откат транзакции

        err := service.CreateOrder(context.Background(), &Order{})
        assert.Error(t, err)
    })
}
~~~

Что проверяем:
Абстрактный эффект — корректный откат транзакции при ошибке.
Явность: Интерфейс TxRepository явно объявляет методы управления транзакциями.

## Использование моков

Используется gomock. Генерация моков через go generate, строгий контроль ожидаемых вызовов с EXPECT().