# Сила низкоуровневых подходов

Используется собственный проприетарный формат поверх gRPC для межсервисного взаимодействия. 

Заменить другим, вероятно, можно было бы, но это вряд ли было бы более эффективно в большинстве случаев.

Однако, если необходимы кастомные интеграции с клиентами за пределами компании, то без форматов JSON, XML не обойтись.

Также для сервисов, которым не важна производительность и которые взаимодействуют с легаси компонентами, существует возможность заменить проприетарный формат на более простой.