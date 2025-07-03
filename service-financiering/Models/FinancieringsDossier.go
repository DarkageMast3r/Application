package models

import (
	"fmt"
	"strconv"
	"time"
)

type FinancieringsDossier struct {
	DossierID     int
	ClientID      int
	ZorgTechID    int
	AanvraagDatum string
	Budget        Budget
	Facturen      []Factuur
}

func (f *FinancieringsDossier) NieuwDossier(clientID int, zorgtechID int) {
	f.ClientID = clientID
	f.ZorgTechID = zorgtechID
}

func (f *FinancieringsDossier) VraagBudgetAan(bedrag float64) {
	f.Budget.NieuwBudget(bedrag)
	t := new(time.Time)
	year, month, day := t.Date()
	f.AanvraagDatum = strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
}

func (f *FinancieringsDossier) VerwerkGoedkeuring(Goedgekeurd bool) {
	if Goedgekeurd {
		f.Budget.BudgetStatus = "Goedgekeurd"
	} else {
		f.Budget.BudgetStatus = "Afgewezen"
	}
}

func (f *FinancieringsDossier) ReserveerBudget() {
	f.Budget.BudgetStatus = "Gereserveerd"
}

func (f *FinancieringsDossier) VerwerkFactuur(factuur Factuur) {
	f.Budget.UpdateBudget(factuur.Bedrag)
	f.Budget.BudgetStatus = "Gefactureerd"
	fmt.Println(factuur)
	fmt.Println(f.Facturen)
}
