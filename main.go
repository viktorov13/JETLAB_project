package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("secret.txt")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		fmt.Println("Файл пуст")
		return
	}

	apiKey := scanner.Text()

	url := "https://app.repeatlab.ru/backend/api/v1/public/project/getAll"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return
	}

	req.Header.Set("Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
