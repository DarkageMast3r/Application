package models

type BudgetStatus int

const (
	Aangevraagd BudgetStatus = iota
	Goedgekeurd
	Afgewezen
	Gereserveerd
	Gefactureerd
)

var statusNaam = map[string]BudgetStatus{
	"Aangevraagd":  Aangevraagd,
	"Goedgekeurd":  Goedgekeurd,
	"Afgewezen":    Afgewezen,
	"Gereserveerd": Gereserveerd,
	"Gefactureerd": Gefactureerd,
}

func (status BudgetStatus) GetStatus(naam string) BudgetStatus {
	return BudgetStatus(statusNaam[naam])
}
