package main

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
