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