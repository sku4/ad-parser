package kufar

import "time"

type Resp struct {
	Ads []struct {
		AccountID         int `json:"account_id"`
		AccountParameters []struct {
			Pl string `json:"pl"`
			Vl string `json:"vl"`
			P  string `json:"p"`
			V  string `json:"v"`
			Pu string `json:"pu"`
		} `json:"account_parameters"`
		AdID         int    `json:"ad_id"`
		AdLink       string `json:"ad_link"`
		AdParameters []struct {
			Pl string      `json:"pl"`
			Vl interface{} `json:"vl"`
			P  string      `json:"p"`
			V  interface{} `json:"v"`
			Pu string      `json:"pu"`
		} `json:"ad_parameters"`
		Body      string `json:"body"`
		Category  string `json:"category"`
		CompanyAd bool   `json:"company_ad"`
		Currency  string `json:"currency"`
		Images    []struct {
			ID           string `json:"id"`
			MediaStorage string `json:"media_storage"`
			YamsStorage  bool   `json:"yams_storage"`
			Path         string `json:"path,omitempty"`
		} `json:"images"`
		ListID       int       `json:"list_id"`
		ListTime     time.Time `json:"list_time"`
		MessageID    string    `json:"message_id"`
		PaidServices struct {
			Halva     bool        `json:"halva"`
			Highlight bool        `json:"highlight"`
			Polepos   bool        `json:"polepos"`
			Ribbons   interface{} `json:"ribbons"`
		} `json:"paid_services"`
		PhoneHidden      bool    `json:"phone_hidden"`
		PriceByn         string  `json:"price_byn"`
		PriceUsd         *string `json:"price_usd"`
		RemunerationType string  `json:"remuneration_type"`
		Subject          string  `json:"subject"`
		Type             string  `json:"type"`
	} `json:"ads"`
	Pagination struct {
		Pages []struct {
			Label string  `json:"label"`
			Num   int     `json:"num"`
			Token *string `json:"token"`
		} `json:"pages"`
	} `json:"pagination"`
	Total int `json:"total"`
}
