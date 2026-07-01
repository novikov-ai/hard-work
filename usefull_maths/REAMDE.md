# Полезная математика для программистов

## Пример 1. Сумма списка и нейтральный элемент

```go
func add(a, b int) int {
	return a + b
}

func Sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total = add(total, n)
	}
	return total
}
```

total — нейтральный элемент: `add(0, x) == x`. Он нужен, чтобы было с чего начинать
цикл, и чтобы не обрабатывать пустой список отдельным edge case'ом.

## Пример 2. Убираем дублирование — обобщённая функция

```go
func Fold(nums []int, identity int, op func(a, b int) int) int {
	acc := identity
	for _, n := range nums {
		acc = op(acc, n)
	}
	return acc
}

// использование:
sum := Fold(nums, 0, add)
max := Fold(nums, math.MinInt, maxOp)
```

Пара (нейтральный элемент, ассоциативная операция) передаётся как параметр в
обобщённый алгоритм свёртки — это и есть "моноид" на практике.

## Пример 3. Generics — обобщаем на любой тип

```go
func Fold[T any](items []T, identity T, op func(a, b T) T) T {
	acc := identity
	for _, x := range items {
		acc = op(acc, x)
	}
	return acc
}

// конкатенация строк
concat := Fold([]string{"foo", "bar", "baz"}, "", func(a, b string) string { return a + b })

// объединение срезов
merged := Fold([][]int{{1, 2}, {3}, {4, 5}}, []int{}, func(a, b []int) []int {
	return append(a, b...)
})

// логическое ИЛИ
anyTrue := Fold([]bool{false, false, true}, false, func(a, b bool) bool { return a || b })
```

Одна функция `Fold`, а конкретное поведение задаётся операцией и нейтральным
элементом.
