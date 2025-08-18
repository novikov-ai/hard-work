# О циклах по умному

Выражаем намерение через абстракции, а не механичный-ручной перебор элементов списка.

## Пример 1

~~~go
func FilterInvalid(txs []Transaction) []Transaction {
    valid := make([]Transaction, 0)
    for _, tx := range txs {
        if tx.IsValid() {
            valid = append(valid, tx)
        }
    }
    return valid
}
~~~

~~~go
func FilterInvalid(txs []Transaction) []Transaction {
    return slices.DeleteFunc(txs, func(tx Transaction) bool {
        return !tx.IsValid()
    })
}
~~~

## Пример 2

~~~go
func FindUser(users []User, id string) *User {
    for _, u := range users {
        if u.ID == id {
            return &u
        }
    }
    return nil
}
~~~

~~~go
func FindUser(users []User, id string) *User {
    if idx := slices.IndexFunc(users, func(u User) bool {
        return u.ID == id
    }); idx != -1 {
        return &users[idx]
    }
    return nil
}
~~~

## Пример 3

~~~go
func TotalAmount(payments []Payment) float64 {
    total := 0.0
    for _, p := range payments {
        total += p.Amount
    }
    return total
}
~~~

~~~go
func TotalAmount(payments []Payment) float64 {
    var total atomic.Float64
    slices.Clip(payments) // Оптимизация памяти
    for _, p := range payments {
        total.Add(p.Amount)
    }
    return total.Load()
}
~~~

## Пример 4

~~~go
func ToDTOs(users []*pb.User) []UserDTO {
    dtos := make([]UserDTO, len(users))
    for i, u := range users {
        dtos[i] = UserDTO{ID: u.Id, Name: u.Name}
    }
    return dtos
}
~~~

~~~go
func Transform[S any, D any](src []S, mapFn func(S) D) []D {
    dst := make([]D, len(src))
    for i, v := range src {
        dst[i] = mapFn(v)
    }
    return dst
}

// Вызов:
dtos := Transform(pbUsers, func(u *pb.User) UserDTO {
    return UserDTO{ID: u.Id, Name: u.Name}
})
~~~

## Пример 5

~~~go
func ProcessBatch(tasks []Task) {
    var wg sync.WaitGroup
    for _, t := range tasks {
        wg.Add(1)
        go func(task Task) {
            defer wg.Done()
            task.Process()
        }(t)
    }
    wg.Wait()
}
~~~

~~~go
func ProcessBatch(tasks []Task) error {
    g, ctx := errgroup.WithContext(context.Background())
    for _, task := range tasks {
        task := task
        g.Go(func() error {
            return task.Process(ctx)
        })
    }
    return g.Wait()
}
~~~