package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

type Project struct {
	UserID      int    `json:"userId"`
	ProjectID   int    `json:"projectId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ProjectName string `json:"projectName"`
}

func readAPI() (string, error) {
	file, err := os.Open("secret.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return "", fmt.Errorf("file is empty")
	}
	return scanner.Text(), nil
}

func fetchProjects(apiKey string) ([]Project, error) {
	url := "https://app.repeatlab.ru/backend/api/v1/public/project/getAll"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Api-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func HTMLHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := readAPI()
	if err != nil {
		http.Error(w, "Read API error", http.StatusInternalServerError)
		return
	}

	projects, err := fetchProjects(apiKey)
	if err != nil {
		http.Error(w, "Fetch projects error", http.StatusBadGateway)
		return
	}

	tmpl, err := template.ParseFiles("templates/projects.html")
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, projects)
}

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := readAPI()
	if err != nil {
		http.Error(w, "Read API error", http.StatusInternalServerError)
		return
	}

	projects, err := fetchProjects(apiKey)
	if err != nil {
		http.Error(w, "Fetch projects error", http.StatusBadGateway)
		return
	}

	jsonData, err := json.MarshalIndent(projects, "", " ")
	if err != nil {
		http.Error(w, "JSON marshal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
