# Избавляемся от зависимости от зависимостей 2.0

## 1. База данных PostgreSQL
**Семантика:**
- Статическая: проект компилируется с зависимостью на драйвер PostgreSQL.
- Динамическая: выполнение запросов к базе во время работы системы.
- Функциональная: возможность сменить базу данных (например, на MySQL) без изменения бизнес-логики.

**Ключевые свойства:**
- Надежность транзакций (ACID).
- Производительность на больших объемах данных.
- Безопасность хранения данных и управление доступом.

**Пространство допустимых изменений:**
- Добавление новых индексов или миграция базы с минимальным влиянием на систему.
- Смена на другую реляционную базу данных при сохранении SQL-запросов.
- Изменения, влияющие на производительность, но не на бизнес-логику.

**Супер-спецификация:**
- Изменения в API базы данных (например, изменение синтаксиса SQL-запросов) не должны требовать модификации бизнес-логики приложения.

## 2. Платежный процессор Stripe
**Семантика:**
- Статическая: проект зависит от библиотеки Stripe для работы с API.
- Динамическая: выполнение платежей через внешние сервисы.
- Функциональная: возможность сменить на другой процессор (например, PayPal) без изменений в коде.

**Ключевые свойства:**
- Надежность и скорость обработки платежей.
- Поддержка различных платежных методов.
- Соответствие стандартам безопасности PCI DSS.

**Пространство допустимых изменений:**
- Добавление новых методов оплаты в API.
- Смена процессора с минимальными изменениями в интерфейсе.
- Изменения в процессе валидации платежей, которые могут повлиять на логику обработки ошибок.

**Супер-спецификация:**
- Любые изменения API процессора должны быть обратными совместимыми, чтобы не нарушать интеграцию с другими сервисами.

## 3. Kafka (система обмена сообщениями)
**Семантика:**
- Статическая: система использует клиент Kafka для общения между микросервисами.
- Динамическая: обработка сообщений в очереди в реальном времени.
- Функциональная: возможность сменить Kafka на другой брокер сообщений без изменения логики обработки.

**Ключевые свойства:**
- Надежная доставка сообщений.
- Масштабируемость для высоконагруженных систем.
- Обработка сообщений в реальном времени.

**Пространство допустимых изменений:**
- Изменения конфигурации кластеров Kafka.
- Смена на другой брокер (RabbitMQ) без переписывания логики продюсера/консюмера.
- Изменение формата сообщений без влияния на целостность данных.

**Супер-спецификация:**
- Изменения в архитектуре очередей сообщений не должны нарушать гарантию доставки сообщений.

## 4. OpenTelemetry (мониторинг и трассировка)
**Семантика:**
- Статическая: подключение зависимости для сбора метрик и трассировок.
- Динамическая: сбор метрик во время выполнения системы.
- Функциональная: возможность сменить библиотеку на другой инструмент мониторинга (Prometheus, Jaeger).

**Ключевые свойства:**
- Сбор детализированных метрик о производительности.
- Масштабируемость мониторинга для распределенных систем.
- Поддержка различных back-end для хранения метрик.

**Пространство допустимых изменений:**
- Добавление новых метрик для мониторинга.
- Смена хранилища для метрик (Prometheus вместо Grafana).
- Изменения в формате сбора данных без нарушения совместимости.

**Супер-спецификация:**
- Любые изменения в сборе метрик должны быть совместимы с текущей инфраструктурой мониторинга и не нарушать работу системы.
