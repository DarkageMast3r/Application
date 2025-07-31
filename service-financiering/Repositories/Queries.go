package repositories

import (
	m "Financiering/Models"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func GetDossiers() []m.FinancieringsDossier {
	var Dossiers []m.FinancieringsDossier
	db := Database_Get()
	innerJoins, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID;")
	if err != nil {
		fmt.Println("Innerjoin: ", err)
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
		fmt.Println("GetDossiers/Scan1: ", err)
			continue
		}
		Dossiers = append(Dossiers, Dossier)
	}
	innerJoins.Close()

	remaining, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID FROM financieringsdossier WHERE BudgetID is null;")
	if err != nil {
		fmt.Println("Remaining: ", err)
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
func GetDossierbyID(ID int) m.FinancieringsDossier {
	var Dossier m.FinancieringsDossier
	db := Database_Get()
	innerJoin, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID WHERE DossierID = ?;", ID)
	if err != nil {
		fmt.Println("Innerjoin: ", err)
		return Dossier
	}

	innerJoin.Next()
	defer innerJoin.Close()
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
		remaining, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID FROM financieringsdossier WHERE BudgetID is null AND DossierID = ?;", ID)
		defer remaining.Close()
		if err != nil {
			fmt.Println("GetDossierbyID/Remaining: ", err)
			return Dossier
		}
		remaining.Next()
		defer remaining.Close()
		err = remaining.Scan(
			&Dossier.DossierID,
			&Dossier.ClientID,
			&Dossier.ZorgTechID,
		)
		if err != nil {
			log.Println("GetDossierbyID/Scan: ", err)
		}
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

// func GetBudgetByID()

func NewBudget(MaxBedrag float64, BeschikbaarBedrag float64, GebruiktBedrag float64, BudgetStatus string) (int, error) {
	db := Database_Get()
	val, err := db.Exec("INSERT INTO budget(MaxBedrag, BeschikbaarBedrag, GebruiktBedrag, BudgetStatus) VALUES(?,?,?,?);", MaxBedrag, BeschikbaarBedrag, GebruiktBedrag, BudgetStatus) // create new budget
	var lastid int
	if err == nil {
		value, err := val.LastInsertId()
		if err != nil {
			log.Println("lastinsertid: ", err)
			return 0, err
		}
		lastid = int(value)
	}
	return lastid, err
}

func ConnectDossier(BudgetID int, DossierID int) error {
	db := Database_Get()
	_, err := db.Exec("UPDATE financieringsdossier SET budgetID = ?, Aanvraagdatum = CURRENT_DATE() WHERE DossierID = ?", BudgetID, DossierID)
	return err
}

// shouldn't even be called if there is no budget
func ProcessPayment(Gebruikt float64, Beschikbaar float64, Status string, ID int) error {
	db := Database_Get()
	_, err := db.Query("UPDATE budget SET GebruiktBedrag = ?, BeschikbaarBedrag = ?, BudgetStatus = ? WHERE ID = ?;", Gebruikt, Beschikbaar, Status, ID)
	return err
}
