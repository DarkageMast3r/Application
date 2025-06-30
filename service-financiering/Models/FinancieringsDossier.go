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

func (f *FinancieringsDossier) VerwerkGoedkeuring(Goedgekeurd bool) {
	if Goedgekeurd {
		f.Budget.BudgetGoedgekeurd()
	} else {
		f.Budget.BudgetAfgewezen()
	}
}

func (f *FinancieringsDossier) ReserveerBudget() {
	f.Budget.BudgetStatus = f.Budget.BudgetStatus.GetStatus("Gereserveerd")
}

func (f *FinancieringsDossier) VerwerkFactuur(factuur Factuur) {
	f.Budget.UpdateBudget(factuur.Bedrag)
	fmt.Println(factuur)
	fmt.Println(f.Facturen)
}
