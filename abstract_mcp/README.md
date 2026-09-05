# Абстракции против паттерна MVC

## Пример 1: Репозиторий объявлений — карточка товара зависит от всего

Бэкенд классифайда. На странице объявления (карточка товара) нужно достать объявление по ID и показать. Но интерфейс репозитория — один:

```go
type ListingRepo interface {
    GetByID(id int64) (*Listing, error)
    GetBySellerID(sellerID int64) ([]*Listing, error)
    Create(l *Listing) error
    Update(l *Listing) error
    Delete(id int64) error
    Bump(id int64) error
    Promote(id int64, plan string) error
}
```

Контроллер карточки вызывает только `GetByID`. Что вышло на практике:

- Добавили `Promote` — изменили сигнатуру, упали тесты карточки, которая к промо вообще не относится
- Мок карточки — семь заглушек ради одного `GetByID`
- Через полгода никто не помнит, зачем карточке `Bump` — потому что не зачем, она его просто не вызывает

**Идеально — ролевые интерфейсы:**

```go
type ListingReader interface {
    GetByID(id int64) (*Listing, error)
    GetBySellerID(sellerID int64) ([]*Listing, error)
}

type ListingWriter interface {
    Create(l *Listing) error
    Update(l *Listing) error
    Delete(id int64) error
}

type ListingPromoter interface {
    Bump(id int64) error
    Promote(id int64, plan string) error
}
```

Контроллер карточки — `ListingReader`. Личный кабинет продавца — `ListingReader` + `ListingWriter`. Раздел промо — `ListingPromoter`. Каждый видит своё.

---

## Пример 2: Сервис продавца — поддержка видит лишнее

У продавца в классифайде есть профиль, статистика, объявления, а ещё бан и верификация. Всё в одном интерфейсе:

```go
type SellerService interface {
    GetProfile(id int64) (*Seller, error)
    GetListings(sellerID int64) ([]*Listing, error)
    UpdateProfile(s *Seller) error
    VerifySeller(id int64) error
    BanSeller(id int64, reason string) error
    GetStats(sellerID int64) (*SellerStats, error)
}
```

Контроллер поддержки (оператор смотрит профиль продавца в чате) зависит от этого интерфейса. Ему нужно `GetProfile` и `GetStats`. А он «видит» `BanSeller`, `VerifySeller`, `UpdateProfile` — и при тестировании мокает всё подряд. Из-за этого:

- Оператор по ошибке наткнулся на метод `BanSeller` и попросил кнопку бана в интерфейсе поддержки — хотя бан должен проходить через отдельный flow с аудитом
- При добавлении OTP-верификации в `VerifySeller` — пришлось трогать моки поддержки, которые к верификации не причастны

**Идеально — разделение по ролям:**

```go
type SellerQuerier interface {
    GetProfile(id int64) (*Seller, error)
    GetStats(sellerID int64) (*SellerStats, error)
    GetListings(sellerID int64) ([]*Listing, error)
}

type SellerManager interface {
    UpdateProfile(s *Seller) error
}

type SellerModerator interface {
    VerifySeller(id int64) error
    BanSeller(id int64, reason string) error
}
```

Поддержка — `SellerQuerier`. Сам продавец — `SellerQuerier` + `SellerManager`. Модерация — `SellerModerator`. Никто не видит чужого.

---

## Пример 3: CRUD-контроллер избранного — публичный API реализует лишнее

У избранного (сохранённые объявления) был один контроллер:

```go
type FavoriteController interface {
    Add(w http.ResponseWriter, r *http.Request)
    Remove(w http.ResponseWriter, r *http.Request)
    List(w http.ResponseWriter, r *http.Request)
}
```

Публичный мобильный API использует `Add` и `List`. От `Remove` он тоже зависит, потому что интерфейс один. Проблемы:

- На ранней стадии `Remove` просто не был реализован (возвращал 501), но интерфейс этого не показывает — компилятор молчал
- Когда добавили `ClearAll` (очистка избранного) — мобильный API сломался на компиляции, хотя `ClearAll` нужен только в админке

**Идеально — контроллеры делятся на Read и Write:**

```go
type FavoriteReader interface {
    List(w http.ResponseWriter, r *http.Request)
}

type FavoriteWriter interface {
    Add(w http.ResponseWriter, r *http.Request)
    Remove(w http.ResponseWriter, r *http.Request)
}

type FavoriteAdmin interface {
    ClearAll(w http.ResponseWriter, r *http.Request)
}
```

Мобильный клиент — `FavoriteReader` + `FavoriteWriter`. Админка — всё разом. Добавление `ClearAll` не трогает мобильный API.