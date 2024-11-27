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