package onliner

import "time"

type Resp struct {
	Apartments []struct {
		ID       int `json:"id"`
		AuthorID int `json:"author_id"`
		Location struct {
			Address     string  `json:"address"`
			UserAddress string  `json:"user_address"`
			Latitude    float64 `json:"latitude"`
			Longitude   float64 `json:"longitude"`
		} `json:"location"`
		Price struct {
			Amount    string `json:"amount"`
			Currency  string `json:"currency"`
			Converted struct {
				BYN struct {
					Amount   string `json:"amount"`
					Currency string `json:"currency"`
				} `json:"BYN"`
				USD struct {
					Amount   *string `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"USD"`
			} `json:"converted"`
		} `json:"price"`
		Photo          string `json:"photo"`
		Resale         bool   `json:"resale"`
		NumberOfRooms  int    `json:"number_of_rooms"`
		Floor          int    `json:"floor"`
		NumberOfFloors int    `json:"number_of_floors"`
		Area           struct {
			Total   *float64 `json:"total"`
			Living  *float64 `json:"living"`
			Kitchen *float64 `json:"kitchen"`
		} `json:"area"`
		Seller struct {
			Type string `json:"type"`
		} `json:"seller"`
		CreatedAt     time.Time `json:"created_at"`
		LastTimeUp    time.Time `json:"last_time_up"`
		UpAvailableIn int       `json:"up_available_in"`
		URL           string    `json:"url"`
		AuctionBid    *struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"auction_bid"`
	} `json:"apartments"`
	Total int `json:"total"`
	Page  struct {
		Limit   int `json:"limit"`
		Items   int `json:"items"`
		Current int `json:"current"`
		Last    int `json:"last"`
	} `json:"page"`
}
