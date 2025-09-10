// package main runs the entire codebase
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type application struct {
	client          *http.Client
	logger          *slog.Logger
	telegramToken   string
	telegramChatID  string
	baseURL         string
	telegramMessage string
}

func NewApp() (*application, error) {
	token := os.Getenv("TELEGRAM_HOMELAB_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	baseURL := os.Getenv("BASE_URL")

	if token == "" || chatID == "" || baseURL == "" {
		return &application{}, errors.New("Missing required environment variables: TELEGRAM_HOMELAB_TOKEN, TELEGRAM_CHAT_ID, BASE_URL")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &application{
		client:         client,
		logger:         logger,
		telegramToken:  token,
		telegramChatID: chatID,
		baseURL:        baseURL,
	}, nil
}

func (a application) sendMessage() error {
	var url = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.telegramToken)
	body, _ := json.Marshal(map[string]string{
		"chat_id": a.telegramChatID,
		"text":    a.telegramMessage,
	})
	req, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return errors.New("failed to send message")
	}
	defer req.Body.Close()
	return nil
}

func (a application) Init() error {
	today := time.Now().Format("1-2-2006")
	a.logger.Info(fmt.Sprintf("Fetching data for %s", today))
	var url = fmt.Sprintf("%s/%s/%s/0", a.baseURL, today, today)

	req, err := a.client.Get(url)
	if err != nil {
		return err
	}

	if req.StatusCode != 200 {
		a.logger.Error("Failed to fetch meal data", slog.Int("status_code", req.StatusCode))
		return err
	}

	var response response
	data, err := io.ReadAll(req.Body)

	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	defer req.Body.Close()

	var message mealList

	for _, menu := range response.MenuSchedules {
		for _, block := range menu.MenuBlocks {
			for _, line := range block.CafeteriaLineList.Data {
				for _, item := range line.FoodItemList.Data {
					if item.LocationName == "CRES- Alternate" {
						continue
					}
					switch block.BlockName {
					case "Breakfast":
						breakfast := meal{
							Type: "Breakfast",
							Item: item.ItemName,
						}
						message.Meals = append(message.Meals, breakfast)
					case "Lunch":
						lunch := meal{
							Type: "Lunch",
							Item: item.ItemName,
						}
						message.Meals = append(message.Meals, lunch)
					}
				}
			}
		}
	}

	var payload string
	payload += fmt.Sprintf("Today is: %s\n", time.Now().Format("2006-01-02"))
	for _, meal := range message.Meals {
		payload += fmt.Sprintf("For %s: %s\n", meal.Type, meal.Item)
	}
	a.telegramMessage = payload
	a.sendMessage()
	return nil
}

func main() {
	app, err := NewApp()
	if err != nil {
		panic(err)
	}
	app.logger.Info("Starting meal fetcher")
	err = app.Init()
	if err != nil {
		app.logger.Error(error.Error(err))
	}
	app.logger.Info("Completed meal fetcher")
}
