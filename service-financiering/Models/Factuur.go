package models

type Factuur struct {
	FactuurID int
	ProductID int
	Bedrag    float64
	Betaald   bool
}