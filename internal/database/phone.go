package database

import (
	"database/sql"
	"fmt"
	"log"
)


func InsertPhone(db *sql.Tx, ctx MigrationContext) error {
	// Insertar directamente en PHONE, limitado al contexto
	insertPhoneQuery := `
	INSERT INTO PHONE (PHONE_TYPE_ID, COUNTRY_CODE, AREA_CODE, PHONE_NUMBER, PHONE_DEFAULT)
	SELECT DISTINCT 1, 56, 9, SUBSTRING_INDEX(TELEFONO, '-', -1), NULL
	FROM temp_csv_asegurados
	WHERE TELEFONO IS NOT NULL
	AND RAMO = ? AND NPOLIZA = ?; -- Filtrar por contexto
	`

	// Asociar los nuevos PHONE_ID al PARTY_ID en PARTY_PHONE, limitado al contexto
	insertPartyPhoneQuery := `
	INSERT IGNORE INTO PARTY_PHONE (PHONE_ID, PARTY_ID)
	SELECT ph.PHONE_ID, p.PARTY_ID
	FROM PHONE ph
	JOIN temp_csv_asegurados t ON ph.PHONE_NUMBER = SUBSTRING_INDEX(t.TELEFONO, '-', -1)
	JOIN PARTY p ON p.EMAIL = t.EMAIL
	WHERE t.TELEFONO IS NOT NULL
	AND t.RAMO = ? AND t.NPOLIZA = ?; -- Filtrar por contexto
	`

	log.Printf("Iniciando inserción de teléfonos para RAMO: %s, NPOLIZA: %s", ctx.Ramo, ctx.Npoliza)

	// Ejecutar la query para insertar teléfonos
	if _, err := db.Exec(insertPhoneQuery, ctx.Ramo, ctx.Npoliza); err != nil {
		return fmt.Errorf("error insertando en PHONE para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	log.Printf("Teléfonos insertados correctamente para RAMO: %s, NPOLIZA: %s.", ctx.Ramo, ctx.Npoliza)
	log.Println(insertPhoneQuery)

	// Ejecutar la query para asociar teléfonos con PARTY
	if _, err := db.Exec(insertPartyPhoneQuery, ctx.Ramo, ctx.Npoliza); err != nil {
		return fmt.Errorf("error asociando PARTY_PHONE para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	log.Printf("PARTY_PHONE asociado correctamente para RAMO: %s, NPOLIZA: %s.", ctx.Ramo, ctx.Npoliza)
	log.Println(insertPartyPhoneQuery)
	return nil
}
