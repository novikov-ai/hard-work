# Полезная математика для программистов

## Пример 1. Операция и ассоциативность

```go
func add(a, b int) int {
	return a + b
}

fmt.Println(add(add(1, 2), 3)) // 6
fmt.Println(add(1, add(2, 3))) // 6
```

Ассоциативность — это факт: `(1+2)+3 == 1+(2+3)`. Результат не зависит от того,
как расставить скобки, то есть от порядка группировки элементов.

## Пример 2. Сумма списка и нейтральный элемент

```go
func Sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total = add(total, n)
	}
	return total
}
```

`0` — нейтральный элемент: `add(0, x) == x`. Он нужен, чтобы было с чего начинать
цикл, и чтобы не обрабатывать пустой список отдельным edge case'ом.

## Пример 3. Параллельная сумма

```go
func SumParallel(nums []int) int {
	mid := len(nums) / 2
	left, right := nums[:mid], nums[mid:]

	var leftSum, rightSum int
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		leftSum = Sum(left)
	}()
	go func() {
		defer wg.Done()
		rightSum = Sum(right)
	}()

	wg.Wait()
	return add(leftSum, rightSum)
}
```

Раз `(1+2)+3 == 1+(2+3)`, сумму можно считать параллельно: разбить список на
части, посчитать каждую в своей горутине, а частичные результаты сложить —
ответ будет тот же самый. Если бы операция была не ассоциативной (например,
вычитание), так делать было бы нельзя.

## Пример 4. Та же схема для максимума

```go
func maxOp(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Max(nums []int) int {
	m := math.MinInt // нейтральный элемент для max
	for _, n := range nums {
		m = maxOp(m, n)
	}
	return m
}
```

Код дословно повторяет структуру `Sum`. Отличаются только операция (`maxOp`
вместо `add`) и нейтральный элемент (`math.MinInt` вместо `0`).

## Пример 5. Убираем дублирование — обобщённая функция

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

## Пример 6. Generics — обобщаем на любой тип

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
