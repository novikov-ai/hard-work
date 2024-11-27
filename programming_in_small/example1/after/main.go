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
	return aggregate(remaining[1:], acc)
}

func main() {
	logs := []LogEntry{
		{"INFO"}, {"ERROR"}, {"INFO"}, {"DEBUG"}, {"ERROR"},
	}
	fmt.Println(aggregateLogsCorecursive(logs))
}