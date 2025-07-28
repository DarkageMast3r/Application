package models

import (
	"log"
	r "repositories"
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
	err := r.InsertDossier(f.ClientID, f.ZorgTechID)
	log.Println(err)
}

func (f *FinancieringsDossier) VraagBudgetAan(bedrag float64) {
	f.Budget.NieuwBudget(bedrag)
	t := new(time.Time)
	year, month, day := t.Date()
	f.AanvraagDatum = strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	err := r.ConnectDossier(f.Budget.ID, f.DossierID) //unknown if those variables are known at this point, otherwise im gonna have to manually retrieve those
	log.Println(err)
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
	err := f.Budget.UpdateBudget(factuur.Bedrag)
	log.Println(err)
	// not sure what to do with this error :KumiThink:
}
