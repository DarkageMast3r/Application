package models

import "fmt"

type FinancieringsDossier struct {
	DossierID     int
	ClientID      int
	ZorgTechID    int
	AanvraagDatum float64 //no datetime type afaik
	Facturen      []Factuur
	Budget        Budget
}

func (f *FinancieringsDossier) Constructor(dossierID int, clientID int, zorgtechID int) {
	f.DossierID = dossierID
	f.ClientID = clientID
	f.ZorgTechID = zorgtechID
}

func (f *FinancieringsDossier) VraagBudgetAan(bedrag float64) {
	f.Budget.BugdetConstructor(bedrag)
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
