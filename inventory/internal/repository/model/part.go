package model

import "time"

type Category string

const (
	CategoryUnspecified Category = "UNSPECIFIED"
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryPorthole    Category = "PORTHOLE"
	CategoryWing        Category = "WING"
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type MetadataValue struct {
	StringValue *string  `bson:"string_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
}

type Part struct {
	UUID          string                   `bson:"uuid"`
	Name          string                   `bson:"name"`
	Description   string                   `bson:"description"`
	Price         float64                  `bson:"price"`
	StockQuantity int64                    `bson:"stock_quantity"`
	Category      Category                 `bson:"category"`
	Dimensions    Dimensions               `bson:"dimensions"`
	Manufacturer  Manufacturer             `bson:"manufacturer"`
	Tags          []string                 `bson:"tags"`
	Metadata      map[string]MetadataValue `bson:"metadata"`
	CreatedAt     time.Time                `bson:"created_at"`
	UpdatedAt     time.Time                `bson:"updated_at"`
}

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
