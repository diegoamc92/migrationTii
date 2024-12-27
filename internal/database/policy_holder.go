package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertPolicyHolderData(db *sql.Tx) error {
	// Query para insertar POLICY_HOLDER si no existe
	InsertPolicyHolderQuery := `
	INSERT IGNORE INTO POLICY_HOLDER (PARTY_ID)
		SELECT PARTY_ID
	FROM PARTY;
	`

	_, err := db.Exec(InsertPolicyHolderQuery)
	if err != nil {
		return fmt.Errorf("error insertando en POLICY_HOLDER: %v", err)
	}

	fmt.Println("Datos insertados correctamente en POLICY_HOLDER.")
	log.Println(InsertPolicyHolderQuery)
	return nil
}
