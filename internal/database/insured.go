package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertInsuredData(db *sql.Tx) error {
	// Query para insertar INSURED si no existe
	insertInsuredQuery := `
	INSERT IGNORE INTO INSURED (PARTY_ID, CLIENT_CODE)
	SELECT PARTY_ID, NULL
	FROM PARTY;
	`

	_, err := db.Exec(insertInsuredQuery)
	if err != nil {
		return fmt.Errorf("error insertando en INSURED: %v", err)
	}

	fmt.Println("Datos insertados correctamente en INSURED.")
	log.Println(insertInsuredQuery)
	return nil
}
