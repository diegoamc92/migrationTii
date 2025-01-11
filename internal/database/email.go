package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertEmail(db *sql.Tx, ctx MigrationContext) error {
	// Insertar directamente en EMAIL, limitado al contexto si aplica
	insertEmailQuery := `
	INSERT INTO EMAIL (EMAIL, EMAIL_TYPE_ID, EMAIL_DEFAULT)
	SELECT DISTINCT EMAIL, 1, NULL
	FROM temp_csv_asegurados
	WHERE EMAIL IS NOT NULL
	AND RAMO = ? AND NPOLIZA = ?;  -- Filtrar por contexto
	`

	// Asociar el nuevo EMAIL_ID al PARTY_ID en PARTY_EMAIL, limitado al contexto si aplica
	insertPartyEmailQuery := `
	INSERT IGNORE INTO PARTY_EMAIL (EMAIL_ID, PARTY_ID)
	SELECT e.EMAIL_ID, p.PARTY_ID
	FROM EMAIL e
	JOIN temp_csv_asegurados t ON e.EMAIL = t.EMAIL
	JOIN PARTY p ON p.PARTY_ID = (
	    SELECT PARTY_ID
	    FROM PARTY
	    WHERE EMAIL = t.EMAIL
	)
	WHERE t.EMAIL IS NOT NULL
	AND t.RAMO = ? AND t.NPOLIZA = ?;  -- Filtrar por contexto
	`

	log.Printf("Iniciando inserción de emails para RAMO: %s, NPOLIZA: %s", ctx.Ramo, ctx.Npoliza)

	// Ejecutar la query para insertar emails
	if _, err := db.Exec(insertEmailQuery, ctx.Ramo, ctx.Npoliza); err != nil {
		return fmt.Errorf("error insertando en EMAIL para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	log.Printf("Emails insertados correctamente para RAMO: %s, NPOLIZA: %s.", ctx.Ramo, ctx.Npoliza)
	log.Println(insertEmailQuery)

	// Ejecutar la query para asociar emails con PARTY
	if _, err := db.Exec(insertPartyEmailQuery, ctx.Ramo, ctx.Npoliza); err != nil {
		return fmt.Errorf("error asociando PARTY_EMAIL para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	log.Printf("PARTY_EMAIL asociado correctamente para RAMO: %s, NPOLIZA: %s.", ctx.Ramo, ctx.Npoliza)
	log.Println(insertPartyEmailQuery)
	return nil
}
