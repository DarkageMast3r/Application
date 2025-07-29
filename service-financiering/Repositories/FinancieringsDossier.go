package repositories

import (
	m "Financiering/Models"
	"log"
	"strconv"
	"time"
)

func NieuwDossier(f *m.FinancieringsDossier, clientID int, zorgtechID int) error {
	f.ClientID = clientID
	f.ZorgTechID = zorgtechID
	return InsertDossier(f.ClientID, f.ZorgTechID)
}

func VraagBudgetAan(f *m.FinancieringsDossier, bedrag float64) {
	NieuwBudget(&f.Budget, bedrag)
	t := new(time.Time)
	year, month, day := t.Date()
	f.AanvraagDatum = strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	err := ConnectDossier(f.Budget.ID, f.DossierID) //unknown if those variables are known at this point, otherwise im gonna have to manually retrieve those
	log.Println(err)
}

func VerwerkGoedkeuring(f *m.FinancieringsDossier, Goedgekeurd bool) {
	if Goedgekeurd {
		f.Budget.BudgetStatus = "Goedgekeurd"
	} else {
		f.Budget.BudgetStatus = "Afgewezen"
	}
}

func ReserveerBudget(f *m.FinancieringsDossier) {
	f.Budget.BudgetStatus = "Gereserveerd"
}

func VerwerkFactuur(f *m.FinancieringsDossier, factuur m.Factuur) {
	err := UpdateBudget(&f.Budget, factuur.Bedrag)
	log.Println(err)
	// not sure what to do with this error :KumiThink:
}
