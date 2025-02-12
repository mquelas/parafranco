package dtos

type HotelDto struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	Amenities   []string `json:"amenities"`
	Photos      []string `json:"photos"`
}
