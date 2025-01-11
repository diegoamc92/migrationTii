package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertPartyWithContext(db *sql.Tx, context MigrationContext) error {
	query := `
	INSERT INTO PARTY (
		EMAIL,
		DATE_CREATED,
		PARTY_SEARCH_AS,
		CIIU_ID,
		NATIONALITY,
		PARTY_ACTIVITY_ID,
		PRACTICE_ID,
		PARTY_CLASS_ID,
		COUNTRY_OF_BIRTH,
		NATIONALITY_DETAIL
	)
	SELECT DISTINCT
		COALESCE(EMAIL, 'migracion@bicevida.cl'),
		NOW(),
		CONCAT(APEPATERNO, ' ', APEMATERNO, ', ', NOMBRES),
		1000, 1, 2, 1, 1000, 136, 136
	FROM temp_csv_asegurados
	WHERE RAMO = ? AND NPOLIZA = ?;
	`

	_, err := db.Exec(query, context.Ramo, context.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando PARTY para póliza %s: %v", context.Npoliza, err)
	}

	log.Printf("PARTY creado para RAMO: %s, NPOLIZA: %s", context.Ramo, context.Npoliza)
	return nil
}
