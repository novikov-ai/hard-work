# Ускоряем код фреймворков в 100 раз

### Примеры ошибочного предположения касательно работы чужого кода 

1. Используя GORM, предположил, что запрос к связанным таблицам всегда выполняется "жадно". Оказывается, это не так. GORM использует "ленивую" загрузку. Чтобы это обойти, можно использовать `Preload` для явного указания, какие связи нужно загрузить. Это не баг, но мое недопонимание поведения по умолчанию.
2. Встречал ситуацию с ответом от стороннего API, когда результат не был отсортирован нужным мне способом. Я полагался на сортировку по первичному ключу, но это было ошибочно и не было задокументировано. 
3. Когда работал с примитивами синхронизации, предположил, что блокировка работает эффективно среди горутин без ухудшения производительности. Но при большой конкурентности заметил, что есть сильная деградация производительности, так как стандартные мьютексы по умолчанию не были оптимизированы для работы под высокой нагрузкой. Эта ситуация явно не была описана в документации.

### Примеры продуманного типа результата, исключающего нежелательные формы обработки

1. Вместо возврата из функции массива слайсов можно использовать неотсортированную мапу, чтобы потребители не думали полагаться на сортировку значений, так как она может быть произвольной.

~~~go
type UserSet map[int]User

func GetUsersByRole(db *gorm.DB, roleName string) (UserSet, error) {
    users := UserSet{}
    err := db.Raw("SELECT * FROM users WHERE role_name = ?", roleName).Scan(&users).Error
    return users, err
}
~~~

2. Явная обработка ошибок. Вместо использования стандартного типа "error" я определяю явный тип ошибки, чтобы сделать обработку более продуманной.
~~~go
type NotFoundError struct {
    message string
}

func (e *NotFoundError) Error() string {
    return e.message
}
~~~

3. Вместо возврата "сырого" результата от ответа API возвращаем обертку над ним, которая валидирует и контролирует как данные используются.

~~~go
// APIResponse is the raw response from an external API.
type APIResponse struct {
    Name  *string `json:"name"`
    Email *string `json:"email"`
}

type User struct {
    Name  string
    Email string
}

func NewUser(apiResp APIResponse) (*User, error) {
    if apiResp.Name == nil || apiResp.Email == nil {
        return nil, fmt.Errorf("invalid user data: missing name or email")
    }

    return &User{
        Name:  *apiResp.Name,
        Email: *apiResp.Email,
    }, nil
}

func main() {
    apiResp := func() APIResponse {
        // implement logic
        return APIResponse{}
    }()

    user, err := NewUser(apiResp)
    if err != nil {
        log.Fatalf("Error creating user: %v", err)
    }

    fmt.Println("User:", user.Name, user.Email)
}
~~~