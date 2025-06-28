package models

import "fmt"

type FinancieringsDossier struct {
	DossierID         int
	ClientID          int
	ZorgTechID        int
	BedragGoedgekeurd float64
	AanvraagDatum     float64 //no datetime type afaik
	Facturen          []Factuur
	Budget            Budget
}

// func (f FinancieringsDossier) Constructor(dossierID int, clientID int, zorgtechID int, budgetStatus BudgetStatus, bedragAangevraagd float64, bedragGoedgekeurd float64)

func (f *FinancieringsDossier) VraagBudgetAan(bedrag float64) {
	f.Budget = f.Budget.Constructor(bedrag)
}

func (f *FinancieringsDossier) VerwerkGoedkeuring(bedrag float64) {
	// On recieving approval do stuff(idk what)
}

func (f *FinancieringsDossier) ReserveerBudget() {
	// Finalise the budget
}

func (f *FinancieringsDossier) VerwerkFactuur(factuur Factuur) {
	// Stuur de factuur/betaal or something, idk the business logic
	// make the factuur turn over a bool i think
	// does this need to be a method like this?
	fmt.Println(factuur)
}
