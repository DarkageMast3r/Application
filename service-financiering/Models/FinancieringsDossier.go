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
	if err1 != nil {
		fmt.Println("Innerjoin: ", err1)
		return Dossiers
	} else if err2 != nil {
		fmt.Println("Innerjoin: ", err2)
		return Dossiers
	}
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
	Result, err, budgetPresent := r.GetDossierbyID(ID)
	if err != nil {
		fmt.Println("Result: ", err)
	}

	if budgetPresent == true {
	Result.Next()
	defer Result.Close()
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
		Result.Next()
		defer Result.Close()
		err = Result.Scan(
			&Dossier.DossierID,
			&Dossier.ClientID,
			&Dossier.ZorgTechID,
		)
		if err != nil {
			fmt.Println("GetDossierbyID/Scan: ", err)
		}
	}
	return Dossier
}