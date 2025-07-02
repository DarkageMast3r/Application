package repositories

import (
	m "Financiering/Models"
	"fmt"
)

func ReadTable() []m.FinancieringsDossier {
	var Dossiers []m.FinancieringsDossier
	db := Database_Get()
	val, err := db.Query("SELECT financieringsdossier.DossierID, financieringsdossier.ClientID, financieringsdossier.ZorgTechID, financieringsdossier.AanvraagDatum, budget.ID, budget.MaxBedrag, budget.BeschikbaarBedrag, budget.GebruiktBedrag, budget.BudgetStatus FROM financieringsdossier INNER JOIN budget on financieringsdossier.BudgetID=budget.ID;")
	if err != nil {
		fmt.Println("ReadTable:", err)
		return Dossiers
	}
	for val.Next() {
		var Dossier m.FinancieringsDossier
		err := val.Scan(
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
			fmt.Println("ReadTable: ", err)
			continue
		}
		Dossiers = append(Dossiers, Dossier)
	}
	return Dossiers
}
