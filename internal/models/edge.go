package models

type Edge struct {
	Id                   string
	To                   VANT
	Weight               float64
	TransmitionRate      float64
	TransmitionAvailable float64
}
