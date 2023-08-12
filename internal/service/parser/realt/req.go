package realt

type ReqGraphQL struct {
	OperationName string       `json:"operationName"`
	Query         string       `json:"query"`
	Variables     ReqVariables `json:"variables"`
}

type ReqVariables struct {
	Data ReqData `json:"data"`
}

type ReqData struct {
	Where             ReqWhere      `json:"where"`
	Pagination        ReqPagination `json:"pagination"`
	Sort              []ReqSort     `json:"sort"`
	ExtraFields       interface{}   `json:"extraFields"`
	IsReactAdaptiveUA bool          `json:"isReactAdaptiveUA"`
}

type ReqWhere struct {
	Category int    `json:"category"`
	Geo      ReqGeo `json:"geo"`
}

type ReqGeo struct {
	Bbox      [][]float64 `json:"bbox"`
	GeoHashes []string    `json:"geoHashes"`
}

type ReqPagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type ReqSort struct {
	By    string `json:"by"`
	Order string `json:"order"`
}
