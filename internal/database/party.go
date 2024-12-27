package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Inser tParty Data inserta registros únicos en la tabla PARTY si no existen.
func InsertPartyData(db *sql.Tx) error {
	// Query para insertar PARTY si no existe
	insertPartyQuery := `
	INSERT INTO PARTY (EMAIL, DATE_CREATED, PARTY_SEARCH_AS, CIIU_ID, NATIONALITY, PARTY_ACTIVITY_ID, PRACTICE_ID, 
	                   PARTY_CLASS_ID, COUNTRY_OF_BIRTH, NATIONALITY_DETAIL)
	SELECT DISTINCT COALESCE(EMAIL, 'migracion@bicevida.cl'),
	       NOW(),
	       CONCAT(APEPATERNO, ' ', APEMATERNO, ', ', NOMBRES),
	       1000, 1, 2, 1, 1000, 136, 136
	FROM temp_csv_asegurados t
	WHERE EMAIL IS NOT NULL
	ON DUPLICATE KEY UPDATE EMAIL=VALUES(EMAIL);
	`

	// Ejecutar la query
	_, err := db.Exec(insertPartyQuery)
	if err != nil {
		return fmt.Errorf("error insertando datos en PARTY: %v", err)
	}
	fmt.Println("Datos insertados en PARTY correctamente.")
	log.Println(insertPartyQuery)

	// Query para contar cuántos registros se insertaron
	countQuery := `SELECT COUNT(*) AS total_insertados FROM PARTY;`

	var totalInsertados int
	err = db.QueryRow(countQuery).Scan(&totalInsertados)
	if err != nil {
		return fmt.Errorf("error al contar registros en PARTY: %v", err)
	}
	fmt.Printf("Total de registros insertados en PARTY: %d\n", totalInsertados)
	return nil
}

//INSERT INTO PARTY (EMAIL, DATE_CREATED, PARTY_SEARCH_AS)
//SELECT DISTINCT COALESCE(EMAIL, 'migracion@bicevida.cl'), NOW(),
//CONCAT(APEPATERNO, ' ', APEMATERNO, ', ', NOMBRES)
//FROM temp_csv_asegurados
//WHERE EMAIL IS NOT NULL
//ON DUPLICATE KEY UPDATE EMAIL=VALUES(EMAIL);
