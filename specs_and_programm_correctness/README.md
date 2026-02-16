# Спецификации и корректность программы

## Пример 1

Комментарий:

`CreateUser` создает пользователя (сохраняет в БД). Ошибкой модульного рассуждения в этом примере
могут служить строки, изменяющие входящую модель: 
`u.CreatedAt = time.Now()` и `u.ID = id`. На такое поведение может кто-то завязаться, а это не гарантируется спецификацией. 

Код:
~~~go
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Предусловие: u.Name и u.Email не пустые, u.Email уникален.
// Постусловие: пользователь сохранён; возвращает ошибку при неудаче.
func (r *UserRepo) CreateUser(u *model.User) error {
	u.CreatedAt = time.Now()

	result, err := r.db.Exec(
		"INSERT INTO users (name, email, role, created_at) VALUES (?, ?, ?, ?)",
		u.Name, u.Email, u.Role, u.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	u.ID = id

	return nil
}
~~~

## Пример 2

Комментарий: 

`GetLatestUserOrder` по спецификации должен возвращать последний заказ, однако в реализации мы видим, что возвращается самый первый в массиве (берется по нулевому индексу). При изменении сортировки метода `GetOrdersByUser` могут возникнуть проблему, если ожидаем, что всегда будет браться элемент по нулевому индексу. 

Код:

~~~go
package usecase

import (
	"fmt"

	"internal/model"
	"internal/repository"
)

type OrderUseCase struct {
	orderRepo   *repository.OrderRepo
	productRepo *repository.ProductRepo
	userRepo    *repository.UserRepo
}

func NewOrderUseCase(
	or *repository.OrderRepo,
	pr *repository.ProductRepo,
	ur *repository.UserRepo,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   or,
		productRepo: pr,
		userRepo:    ur,
	}
}

// Предусловие: userID > 0.
// Постусловие: возвращает последний заказ или ошибку, если заказов нет.
func (uc *OrderUseCase) GetLatestUserOrder(userID int64) (*model.Order, error) {
	orders, err := uc.orderRepo.GetOrdersByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found for user %d", userID)
	}

	// Первый элемент — самый свежий
	return &orders[0], nil
}
~~~