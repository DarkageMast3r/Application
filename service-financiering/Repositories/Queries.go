package repositories

import (
	m "Financiering/Models"
	"fmt"
)

func GetDossiers() []m.FinancieringsDossier {
	var Dossiers []m.FinancieringsDossier
	db := Database_Get()
	innerJoins, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID;")
	if err != nil {
		fmt.Println(err)
		return Dossiers
	}
	for innerJoins.Next() {
		var Dossier m.FinancieringsDossier
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
			fmt.Println(err)
			continue
		}
		Dossiers = append(Dossiers, Dossier)
	}

	remaining, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID FROM financieringsdossier WHERE BudgetID is null;")
	if err != nil {
		fmt.Println(err)
		return Dossiers
	}
	for remaining.Next() {
		var Dossier m.FinancieringsDossier
		err := remaining.Scan(
			&Dossier.DossierID,
			&Dossier.ClientID,
			&Dossier.ZorgTechID,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Dossiers = append(Dossiers, Dossier)
	}
	for i, val := range Dossiers {
		if val.Budget.BudgetStatus == "" {
			Dossiers[i].Budget.BudgetStatus = "Niet aangevraagd"
		}
	}
	return Dossiers
}

// Only works when it has a budget, perhaps add functionality for when there isnt a budget
func GetDossierbyID(ID int) m.FinancieringsDossier {
	var Dossier m.FinancieringsDossier
	db := Database_Get()
	innerJoin, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID WHERE financieringsdossier.DossierID = ?;", ID)
	if err != nil {
		fmt.Println(err)
		return Dossier
	}
	innerJoin.Next()
	err = innerJoin.Scan(
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
		fmt.Println(err)
		return Dossier
	}
	return Dossier
}

func InsertDossier(clientID int, zorgTechID int) error {
	db := Database_Get()
	_, err := db.Query("INSERT INTO financieringsdossier(ClientID, ZorgTechID) VALUES(?,?)", clientID, zorgTechID)
	return err
}

func RemoveDossier(dossierID int) error {
	db := Database_Get()
	_, err := db.Query("DELETE FROM financieringsdossier WHERE DossierID = ?", dossierID)
	return err
}
