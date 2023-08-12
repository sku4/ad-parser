package realt

import "time"

type RespClustered struct {
	Data struct {
		Type     string `json:"type"`
		Features []struct {
			Type     string      `json:"type"`
			ID       string      `json:"id"`
			Bbox     [][]float64 `json:"bbox"`
			Number   int         `json:"number"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				Data struct {
					GeoHash             string        `json:"geoHash"`
					MinPrice            int           `json:"minPrice"`
					MinPricePerM2       int           `json:"minPricePerM2"`
					MinPricePerPerson   int           `json:"minPricePerPerson"`
					Uuids               []interface{} `json:"uuids"`
					PointUUID           string        `json:"pointUuid"`
					AddressIntersection bool          `json:"addressIntersection"`
					Category            int           `json:"category"`
				} `json:"data"`
			} `json:"properties"`
			UUID string `json:"uuid"`
		} `json:"features"`
	} `json:"data"`
}

type RespGraphQL struct {
	Data struct {
		SearchObjects struct {
			Body struct {
				Results []struct {
					UUID              string       `json:"uuid"`
					Title             string       `json:"title"`
					Description       string       `json:"description"`
					Headline          interface{}  `json:"headline"`
					CreatedAt         time.Time    `json:"createdAt"`
					UpdatedAt         time.Time    `json:"updatedAt"`
					MetroTime         interface{}  `json:"metroTime"`
					MetroTimeType     interface{}  `json:"metroTimeType"`
					Price             *float64     `json:"price"`
					PriceCurrency     *float64     `json:"priceCurrency"`
					PricePerM2        *float64     `json:"pricePerM2"`
					PricePerM2Max     *float64     `json:"pricePerM2Max"`
					PricePerPerson    *float64     `json:"pricePerPerson"`
					PriceMin          *float64     `json:"priceMin"`
					PriceMax          *float64     `json:"priceMax"`
					Storeys           *uint8       `json:"storeys"`
					Storey            *uint8       `json:"storey"`
					Rooms             *uint8       `json:"rooms"`
					ContactPhones     []string     `json:"contactPhones"`
					Images            []string     `json:"images"`
					AreaTotal         *float64     `json:"areaTotal"`
					AreaLiving        *float64     `json:"areaLiving"`
					AreaKitchen       *float64     `json:"areaKitchen"`
					AreaMax           *interface{} `json:"areaMax"`
					AreaMin           *interface{} `json:"areaMin"`
					AreaLand          *interface{} `json:"areaLand"`
					ObjectType        interface{}  `json:"objectType"`
					Code              int          `json:"code"`
					StateRegionName   string       `json:"stateRegionName"`
					StateDistrictName string       `json:"stateDistrictName"`
					TownType          int          `json:"townType"`
					TownName          string       `json:"townName"`
					StreetName        *string      `json:"streetName"`
					Address           *string      `json:"address"`
					ContactName       string       `json:"contactName"`
					AgencyName        string       `json:"agencyName"`
					MetroStationName  interface{}  `json:"metroStationName"`
					MetroLineID       interface{}  `json:"metroLineId"`
					HouseNumber       *int         `json:"houseNumber"`
					BuildingNumber    *interface{} `json:"buildingNumber"`
					PaymentStatus     int          `json:"paymentStatus"`
					Comments          string       `json:"comments"`
					IsFavorite        bool         `json:"isFavorite"`
					Category          int          `json:"category"`
					Has3DTour         bool         `json:"has3dTour"`
					HasVideo          bool         `json:"hasVideo"`
					StateRegionUUID   string       `json:"stateRegionUuid"`
					NumberOfBeds      interface{}  `json:"numberOfBeds"`
					DirectionName     interface{}  `json:"directionName"`
					TownDistance      interface{}  `json:"townDistance"`
					CustomSorting     int          `json:"customSorting"`
					SpecialComment    interface{}  `json:"specialComment"`
					Location          []float64    `json:"location"`
					Year              *uint16      `json:"buildingYear"`
					Toilet            *int         `json:"toilet"`
				} `json:"results"`
				Pagination struct {
					Page       int `json:"page"`
					PageSize   int `json:"pageSize"`
					TotalCount int `json:"totalCount"`
				} `json:"pagination"`
				Rates []struct {
					From int     `json:"from"`
					To   int     `json:"to"`
					Rate float64 `json:"rate"`
				} `json:"rates"`
				ExtraFields interface{} `json:"extraFields"`
			} `json:"body"`
			Success bool          `json:"success"`
			Errors  []interface{} `json:"errors"`
		} `json:"searchObjects"`
	} `json:"data"`
}
