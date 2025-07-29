package models

type FinancieringsDossier struct {
	DossierID     int
	ClientID      int
	ZorgTechID    int
	AanvraagDatum string
	Budget        Budget
	Facturen      []Factuur
}
