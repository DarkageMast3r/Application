package repositories

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func GetDossiers() (*sql.Rows, *sql.Rows, error, error) {
	db := Database_Get()
	innerJoins, err1 := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID;")

	remaining, err2 := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID FROM financieringsdossier WHERE BudgetID is null;")

	
	return innerJoins, remaining, err1, err2
}

func GetDossierbyID(ID int, BudgetPresent bool) (*sql.Rows, error) {
	db := Database_Get()
	if BudgetPresent == true {
	innerJoin, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID WHERE DossierID = ?;", ID)
	return innerJoin, err
	} else {
		remaining, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID FROM financieringsdossier WHERE BudgetID is null and DossierID = ?;", ID)
		return remaining, err
	}
}


func InsertDossier(clientID int, zorgTechID int) (sql.Result, error) {
	db := Database_Get()
	res, err := db.Exec("INSERT INTO financieringsdossier(ClientID, ZorgTechID) VALUES(?,?)", clientID, zorgTechID)
	return res, err
}

// shouldn't even be called if there is no budget
func ProcessPayment(Gebruikt float64, Beschikbaar float64, Status string, ID int) error {
	db := Database_Get()
	_, err := db.Query("UPDATE budget SET GebruiktBedrag = ?, BeschikbaarBedrag = ?, BudgetStatus = ? WHERE ID = ?;", Gebruikt, Beschikbaar, Status, ID)
	return err
}

func InsertBudget(MaxBedrag float64, BeschikbaarBedrag float64, GebruiktBedrag float64, BudgetStatus string) (sql.Result, error) {
	db := Database_Get()
	res, err := db.Exec("INSERT INTO budget(MaxBedrag, BeschikbaarBedrag, GebruiktBedrag, BudgetStatus) VALUES(?,?,?,?);", MaxBedrag, BeschikbaarBedrag, GebruiktBedrag, BudgetStatus) // create new budget
	return res, err
}

func RemoveDossier(dossierID int) error {
	db := Database_Get()
	_, err := db.Query("DELETE FROM financieringsdossier WHERE DossierID = ?", dossierID)
	return err
}

func ConnectDossier(BudgetID int, DossierID int) error {
	db := Database_Get()
	_, err := db.Exec("UPDATE financieringsdossier SET budgetID = ?, Aanvraagdatum = CURRENT_DATE() WHERE DossierID = ?", BudgetID, DossierID)
	return err
}