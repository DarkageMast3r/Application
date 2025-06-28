package models

import (
	"fmt"
	"time"
)

type FinancieringsDossier struct {
	DossierID     int
	ClientID      int
	ZorgTechID    int
	AanvraagDatum Date
	Facturen      []Factuur
	Budget        Budget
}

func (f *FinancieringsDossier) NieuwDossier(dossierID int, clientID int, zorgtechID int) {
	f.DossierID = dossierID
	f.ClientID = clientID
	f.ZorgTechID = zorgtechID
}

func (f *FinancieringsDossier) VraagBudgetAan(bedrag float64) {
	f.Budget.NieuwBudget(bedrag)
	t := new(time.Time)
	f.AanvraagDatum.Year, f.AanvraagDatum.Month, f.AanvraagDatum.Day = t.Date()
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
