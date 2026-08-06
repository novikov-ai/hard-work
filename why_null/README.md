# Зачем вам NULL

## Пример 1

Существует табличка `element` со следующей схемой:
~~~sql
create table element (
    id bigserial primary key not null,
    hash text not null unique,
    title text,
    -- ... 
);
~~~

Поле `title` задано как null, предполагая, что оно может быть не заполнено. Лучше сделать поле NOT NULL и договориться, что отсутствие значения означает то же, что и null.

## Пример 2

Существует табличка `element` со следующей схемой:
~~~sql
create table element (
    id bigserial primary key not null,
    hash text not null unique,
    title text not null,
    price_type text,
    -- ...
);
~~~

Поле `price_type` задано как null, предполагая, что оно может быть не заполнено. 

Вместо null-поля лучше было бы использовать NOT NULL-поле + добавить поле-справочник `price_type_code` int not null, где будет находиться код единицы измерения.

## Пример 3

В схеме issue-трекеров вроде GitLab/Redmine таблица issues обычно выглядит так:

~~~sql
create table issues (
    id bigserial primary key not null,
    project_id bigint not null references projects(id),
    title text not null,
    assignee_id bigint references users(id), -- nullable
    due_date date,                            -- nullable
    -- ...
);
~~~~

assignee_id и due_date — nullable, потому что задача может быть не назначена и не иметь дедлайна.

Лучше вынести назначение в отдельную таблицу-компаньон (паттерн "отдельная таблица" из прошлого ответа — отсутствие строки = отсутствие назначения):

~~~sql
create table issues (
    id bigserial primary key not null,
    project_id bigint not null references projects(id),
    title text not null,
    -- assignee_id и due_date убраны отсюда
);

create table issue_assignments (
    issue_id bigint primary key references issues(id),
    assignee_id bigint not null references users(id),
    assigned_at timestamptz not null default now(),
    unassigned_reason_id int not null default 0
        references assignment_reason(id)
);
~~~

## Пример 4

Номер телефона продавца в объявлении:

~~~sql
create table ads (
    id bigserial primary key not null,
    seller_id bigint not null references sellers(id),
    phone text, -- nullable: "продавец не указал телефон"
    -- ...
);
~~~

phone = NULL смешивает три разных случая: продавец не ввёл номер, продавец скрыл номер намеренно (только чат), номер временно недоступен из-за модерации/бана. Все три в коде обрабатываются одним и тем же if ad.Phone != nil, и продукту сложно
потом различить "нет номера" от "не хочет показывать".

~~~sql
create table ads (
    id bigserial primary key not null,
    seller_id bigint not null references sellers(id),
    phone text not null default '',
    contact_method_id int not null default 1
        references contact_method(id)
    -- 1 = только телефон, 2 = только чат, 3 = телефон и чат, 4 = не указан
);
~~~

Теперь Phone string в Go — без указателя, а вся логика "что показывать в UI" завязана на ContactMethodID, а не на проверку строки на пустоту вместе с NULL.

## Пример 5

Обложка (первое фото) объявления:

~~~sql
create table ads (
    id bigserial primary key not null,
    seller_id bigint not null references sellers(id),
    cover_photo_url text, -- nullable: "фото ещё не загружено"
    -- ...
);
~~~

Здесь NULL используют как "фото нет", но в объявлениях без фото это законное бизнес-состояние (например, объявление услуг), а не ошибка ввода — и такие объявления часто нужно занижать в поисковой выдаче. Проверка WHERE cover_photo_url IS NOT
NULL работает, но плохо сочетается с индексами и фильтрами по нескольким условиям сразу.

~~~sql
create table ads (
    id bigserial primary key not null,
    seller_id bigint not null references sellers(id),
    has_photos boolean not null default false
    -- cover_photo_url убрано отсюда
);

create table ad_photos (
    id bigserial primary key not null,
    ad_id bigint not null references ads(id),
    url text not null,
    sort_order int not null default 0,
    unique (ad_id, sort_order)
);
~~~

has_photos в ads — денормализованный, но NOT NULL флаг для быстрой фильтрации в листинге (WHERE has_photos = true), а реальные фото и их порядок живут в ad_photos; отсутствие строк там = отсутствие фото, без единого NULL в обеих таблицах.