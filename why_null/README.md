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