package database

import (
	"database/sql"
	"fmt"
)

func InsertIdentificationWithContext(tx *sql.Tx) error {
	query := `
	INSERT INTO IDENTIFICATION (IDENTIFICATION, IDENTIFICATION_TYPE_ID)
	SELECT DISTINCT RUT, 1
	FROM tiisa.asegurados
	WHERE RUT IS NOT NULL AND RUT != ''
	ON DUPLICATE KEY UPDATE IDENTIFICATION=VALUES(IDENTIFICATION);
	`

	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("error insertando en IDENTIFICATION desde asegurados: %v", err)
	}
	fmt.Println("Datos insertados en IDENTIFICATION correctamente desde asegurados.")
	return nil
}

func AssociatePartyIdentificationWithContext(tx *sql.Tx, context MigrationContext) error {
	query := `
	INSERT IGNORE INTO PARTY_IDENTIFICATION (PARTY_ID, IDENTIFICATION_ID)
	SELECT p.PARTY_ID, i.IDENTIFICATION_ID
	FROM PARTY p
	JOIN tiisa.asegurados t 
		ON p.PARTY_SEARCH_AS = CONCAT_WS(', ', t.APEPATERNO, t.APEMATERNO, t.NOMBRES)
		AND t.RAMO = ? AND t.NPOLIZA = ?
	JOIN IDENTIFICATION i 
		ON i.IDENTIFICATION = t.RUT
	WHERE t.RUT IS NOT NULL AND t.RUT != ''
	ON DUPLICATE KEY UPDATE IDENTIFICATION_ID = i.IDENTIFICATION_ID;
	`

	_, err := tx.Exec(query, context.Ramo, context.Npoliza)
	if err != nil {
		return fmt.Errorf("error asociando PARTY_IDENTIFICATION para póliza %s: %v", context.Npoliza, err)
	}
	fmt.Printf("PARTY_IDENTIFICATION asociada correctamente para RAMO: %s, NPOLIZA: %s.\n", context.Ramo, context.Npoliza)
	return nil
}
