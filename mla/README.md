# Формализуем многослойную (многоуровневую) архитектуру

В проекте разрабатываемого сервиса поддерживается многоуровневная архитектура. 

Существует четыре основных уровней.

## Уровень 1 - Service & Worker

На этом уровне расположена инициализация сервиса и регистрация служебных утилит (регистрация хендлеров, инициализация метрик и пр.)

В параллель с инициализацией сервиса существует воркер, который также инициирует собственную работу и служебные утилиты. 

## Уровень 2 - Handlers

Второй уровень отвечает за бизнес-логику всех обработчиков проекта. 

## Уровень 3 - Main Usecases

На этом уровне располагаются основные сценарии использования сервиса. Один из них отвечает за CRUD-операции. 

## Уровень 4 - Storage и служебные Usecases

На уровне хранилища происходит инкапсуляция по работе с базами данных. 

## Диаграмма 

+--------------------------+           +--------------------------+
|          Service         |           |         Worker           |
+--------------------------+           +--------------------------+
             |                                      |
             v                                      v
+--------------------------+           +--------------------------+
|         Handlers         |           |         Handlers         |
+--------------------------+           +--------------------------+
            |                                       |
            v                                       v
+--------------------------+           +--------------------------+
|     Usecases (CRUD)      |           |     Usecases (CRUD)      |
+--------------------------+           +--------------------------+
            |                                       |
            v                                       v
+--------------------------+           +--------------------------+
|     Storage & Usecases   |           |     Storage & Usecases   |
+--------------------------+           +--------------------------+