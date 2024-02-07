package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type response struct {
	MenuSchedules []struct {
		MenuBlocks []struct {
			BlockName         string `json:"blockName"`
			ScheduledDate     string `json:"scheduledDate"`
			CafeteriaLineList struct {
				Data []struct {
					Name         string `json:"name"`
					FoodItemList struct {
						Data []struct {
							LocationName string `json:"location_name"`
							ItemName     string `json:"item_Name"`
							Description  string `json:"description"`
						} `json:"data"`
					} `json:"foodItemList"`
				} `json:"data"`
			} `json:"cafeteriaLineList"`
		} `json:"menuBlocks"`
	} `json:"menuSchedules"`
}

type meal struct {
	Type        string
	Date        string
	School      string
	Item        string
	Description string
}

type mealList struct {
	Meals []meal
}

func sendMessage(msg string) {
	token := os.Getenv("TELEGRAM_HOMELAB_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if token == "" || chatID == "" {
		fmt.Println("Missing token or chat id")
		return
	}

	var url = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	body, _ := json.Marshal(map[string]string{
		"chat_id": chatID,
		"text":    msg,
	})
	req, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	defer req.Body.Close()
}

func getMeals() string {
	today := time.Now().Format("1-2-2006")
	baseURL := os.Getenv("BASE_URL")
	if !strings.HasPrefix(baseURL, "http") {
		return "URL Not provided"
	}

	fmt.Println(baseURL)
	var url = fmt.Sprintf("%s/%s/%s/0", baseURL, today, today)

	req, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if req.StatusCode != 200 {
		fmt.Println("Error: ", req.StatusCode)
		return ""
	}

	var response response
	data, err := io.ReadAll(req.Body)

	json.Unmarshal(data, &response)

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
	sendMessage(payload)
	return payload
}

func main() {
	meals := getMeals()
	fmt.Println(meals)
}
