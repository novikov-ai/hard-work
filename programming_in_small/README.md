# Как проектировать программы in small

### Пример 1

До:

~~~go
package main

import (
	"fmt"
)

type LogEntry struct {
	Level string
}

func aggregateLogs(logs []LogEntry) map[string]int {
	stats := make(map[string]int)
	for _, log := range logs {
		stats[log.Level]++
	}
	return stats
}

func main() {
	logs := []LogEntry{
		{"INFO"}, {"ERROR"}, {"INFO"}, {"DEBUG"}, {"ERROR"},
	}
	fmt.Println(aggregateLogs(logs))
}
~~~

После:
~~~go
package main

import (
	"fmt"
)

type LogEntry struct {
	Level string
}

func aggregateLogsCorecursive(logs []LogEntry) map[string]int {
	return aggregate(logs, make(map[string]int))
}

func aggregate(remaining []LogEntry, acc map[string]int) map[string]int {
	if len(remaining) == 0 {
		return acc
	}

	acc[remaining[0].Level]++
	return aggregate(remaining[1:], acc) // заменили итерацию на рекурсию, обрабатываем 1 лог за 1 раз
}

func main() {
	logs := []LogEntry{
		{"INFO"}, {"ERROR"}, {"INFO"}, {"DEBUG"}, {"ERROR"},
	}
	fmt.Println(aggregateLogsCorecursive(logs))
}
~~~

### Пример 2

До:

~~~go
package main

import (
	"fmt"
)

type APIResponse struct {
	Items  []string
	Cursor string
}

func fetchAPI(cursor string) APIResponse {
	// Мокаем API-ответ
	if cursor == "end" {
		return APIResponse{}
	}
	return APIResponse{
		Items:  []string{"item1", "item2"},
		Cursor: "end",
	}
}

func getAllItems() []string {
	var allItems []string
	cursor := ""
	for {
		resp := fetchAPI(cursor)
		allItems = append(allItems, resp.Items...)
		if resp.Cursor == "" {
			break
		}
		cursor = resp.Cursor
	}
	return allItems
}

func main() {
	fmt.Println(getAllItems())
}
~~~

После:

~~~go
package main

import (
	"fmt"
)

type APIResponse struct {
	Items  []string
	Cursor string
}

func fetchAPI(cursor string) APIResponse {
	// Мокаем API-ответ
	if cursor == "end" {
		return APIResponse{}
	}
	return APIResponse{
		Items:  []string{"item1", "item2"},
		Cursor: "end",
	}
}

// Рекурсивный сбор всех элементов с пагинацией
func getAllItemsCorecursive(cursor string, acc []string) []string {
	resp := fetchAPI(cursor)
	acc = append(acc, resp.Items...)
	if resp.Cursor == "" {
		return acc
	}
	return getAllItemsCorecursive(resp.Cursor, acc)
}

func main() {
	fmt.Println(getAllItemsCorecursive("", []string{}))
}
~~~

### Пример 3

До:

~~~go
package main

import (
	"fmt"
	"sort"
)

type Event struct {
	PlayerID string
	Score    int
}

type Player struct {
	ID    string
	Score int
}

func updateLeaderboard(events []Event) []Player {
	leaderboard := make(map[string]int)
	for _, event := range events {
		leaderboard[event.PlayerID] += event.Score
	}

	// Convert to sorted list
	var sortedLeaderboard []Player
	for id, score := range leaderboard {
		sortedLeaderboard = append(sortedLeaderboard, Player{ID: id, Score: score})
	}

	sort.Slice(sortedLeaderboard, func(i, j int) bool {
		return sortedLeaderboard[i].Score > sortedLeaderboard[j].Score
	})

	return sortedLeaderboard
}

func main() {
	events := []Event{
		{"player1", 100}, {"player2", 50}, {"player1", -30}, {"player3", 200},
	}
	fmt.Println(updateLeaderboard(events))
}
~~~

После:

~~~go
package main

import (
	"fmt"
	"sort"
)

type Event struct {
	PlayerID string
	Score    int
}

type Player struct {
	ID    string
	Score int
}

// Итерация заменена корекурсией: каждое событие обрабатывается по одному, используется аккумулятор
func updateLeaderboardCorecursive(events []Event) []Player {
	return update(events, make(map[string]int))
}

// Построение результата следуется за выходными данными: после обработки всех событий аккумулятор конвертируется в сортированный список
// Сортировка и преобразование выполняется один раз в конце
func update(remaining []Event, acc map[string]int) []Player {
	if len(remaining) != 0 {
		acc[remaining[0].PlayerID] += remaining[0].Score
		return update(remaining[1:], acc)
	}

	var sortedLeaderboard []Player
	for id, score := range acc {
		sortedLeaderboard = append(sortedLeaderboard, Player{ID: id, Score: score})
	}
	sort.Slice(sortedLeaderboard, func(i, j int) bool {
		return sortedLeaderboard[i].Score > sortedLeaderboard[j].Score
	})
	
	return sortedLeaderboard
}

func main() {
	events := []Event{
		{"player1", 100}, {"player2", 50}, {"player1", -30}, {"player3", 200},
	}
	fmt.Println(updateLeaderboardCorecursive(events))
}
~~~

