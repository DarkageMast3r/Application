package models

type BudgetStatus int

const (
	Aangevraagd BudgetStatus = iota
	Goedgekeurd
	Afgewezen
	Gereserveerd
	Gefactureerd
)

var statusNaam = map[BudgetStatus]string{
	Aangevraagd:  "Aangevraagd",
	Goedgekeurd:  "Goedgekeurd",
	Afgewezen:    "Afgewezen",
	Gereserveerd: "Gereserveerd",
	Gefactureerd: "Gefactureerd",
}

func (status BudgetStatus) GetStatus() string {
	return statusNaam[status]
}
