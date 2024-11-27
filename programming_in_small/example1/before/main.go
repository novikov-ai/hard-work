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