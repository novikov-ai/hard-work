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