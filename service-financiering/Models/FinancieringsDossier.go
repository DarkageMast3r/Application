package models

import (
	r "Financiering/Repositories"
	"fmt"
)

type FinancieringsDossier struct {
	DossierID     int
	ClientID      int
	ZorgTechID    int
	AanvraagDatum string
	Budget        Budget
	Facturen      []Factuur
}

func GetAllDossiers() []FinancieringsDossier {
	var Dossiers []FinancieringsDossier
	innerJoins, remaining, err1, err2 := r.GetDossiers()
	if err1 == nil {
		for innerJoins.Next() {
			var Dossier FinancieringsDossier
			err := innerJoins.Scan(
			&Dossier.DossierID,
				&Dossier.ClientID,
				&Dossier.ZorgTechID,
				&Dossier.AanvraagDatum,
				&Dossier.Budget.ID,
				&Dossier.Budget.MaxBedrag,
				&Dossier.Budget.BeschikbaarBedrag,
				&Dossier.Budget.GebruiktBedrag,
				&Dossier.Budget.BudgetStatus,
			)
			if err != nil {
			fmt.Println("GetDossiers/Scan1: ", err)
				continue
			}
			Dossiers = append(Dossiers, Dossier)
		}
		innerJoins.Close()
	} else if err2 == nil {
		for remaining.Next() {
			var Dossier FinancieringsDossier
			err := remaining.Scan(
				&Dossier.DossierID,
				&Dossier.ClientID,
				&Dossier.ZorgTechID,
			)
			if err != nil {
			fmt.Println("GetDossiers/Scan2: ", err)
				continue
			}
			Dossiers = append(Dossiers, Dossier)
		}
		remaining.Close()
	} else {
		fmt.Println(err1)
		fmt.Println("")
		fmt.Println(err2)
	}	
	
	for i, val := range Dossiers {
		if val.Budget.BudgetStatus == "" {
			Dossiers[i].Budget.BudgetStatus = "Niet aangevraagd"
		}
	}

	return Dossiers
}

// Only works when it has a budget, perhaps add functionality for when there isnt a budget
func GetDossierbyID(ID int) FinancieringsDossier {
	var Dossier FinancieringsDossier
	Result, err := r.GetDossierbyID(ID, true)
	if err != nil {
		fmt.Println("Result: ", err)
	}
	nextable := Result.Next()
	defer Result.Close()
	if nextable == true {
		Result.Scan(
		&Dossier.DossierID,
		&Dossier.ClientID,
		&Dossier.ZorgTechID,
		&Dossier.AanvraagDatum,
		&Dossier.Budget.ID,
		&Dossier.Budget.MaxBedrag,
		&Dossier.Budget.BeschikbaarBedrag,
		&Dossier.Budget.GebruiktBedrag,
		&Dossier.Budget.BudgetStatus,
	)
	} else {
		Result, err = r.GetDossierbyID(ID, false)	
		nextable = Result.Next()
		if nextable == true {
			err = Result.Scan(
				&Dossier.DossierID,
				&Dossier.ClientID,
				&Dossier.ZorgTechID,
			)
			if err != nil {
				fmt.Println("GetDossierbyID/Scan: ", err)
			}
		} else {
			fmt.Println("Unable to input data")
		}
	}
	return Dossier
}

func GetClientBudget(clientID int) Budget {
	var Budget Budget
	Result, err := r.GetBudgetbyClientID(clientID)
	if err != nil {
		fmt.Println("GetClientBudget: ", err)
		return Budget
	}
	nextable := Result.Next()
	defer Result.Close()
	if nextable == true {
		Result.Scan(
		&Budget.ID,
		&Budget.MaxBedrag,
		&Budget.BeschikbaarBedrag,
		&Budget.GebruiktBedrag,
		&Budget.BudgetStatus,
	)
	}
	return Budget
}

// wip
// func VerwerkGoedkeuring(f *m.FinancieringsDossier, Goedgekeurd bool) {
// 	// if Goedgekeurd {
// 	// 	f.Budget.BudgetStatus = "Goedgekeurd"
// 	// } else {
// 	// 	f.Budget.BudgetStatus = "Afgewezen"
// 	// }
// }

// // wip
// func ReserveerBudget(f *m.FinancieringsDossier) {
// 	// f.Budget.BudgetStatus = "Gereserveerd"
// }

// //wip
// func VerwerkFactuur(f *m.FinancieringsDossier, factuur m.Factuur) {
// 	// err := UpdateBudget(&f.Budget, factuur.Bedrag)
// 	// log.Println(err)
// 	// not sure what to do with this error :KumiThink:
// }