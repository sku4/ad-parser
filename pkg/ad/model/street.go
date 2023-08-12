package model

type StreetTnt struct {
	ID   uint64 `mapstructure:"id" json:"id"`
	Name string `mapstructure:"name" json:"name"`
	Type uint8  `mapstructure:"type" json:"type"`
}
