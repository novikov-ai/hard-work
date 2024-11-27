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

// Обновление статистики лидеров с использованием корекурсии
func updateLeaderboardCorecursive(events []Event) []Player {
	return update(events, make(map[string]int))
}

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
