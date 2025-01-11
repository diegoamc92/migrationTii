package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertPersonInsured(db *sql.Tx, ctx *MigrationContext) error {
	query := `
    INSERT INTO PERSON_INSURED (
        ITEM_ID,
        CONTRACT_ID,
        PARTY_ID,
        PERSON_INSURED_TYPE_ID
    )
    SELECT
        ROW_NUMBER() OVER (PARTITION BY t.RAMO, t.NPOLIZA) AS ITEM_ID,
        ? AS CONTRACT_ID, -- ID del contrato asociado
        p.PARTY_ID,
        CASE
            WHEN t.RUT = ? THEN 1 -- Titular
            ELSE 2 -- Carga
        END AS PERSON_INSURED_TYPE_ID
    FROM temp_csv_asegurados t
    JOIN PARTY p ON p.EMAIL = t.EMAIL
    WHERE t.RAMO = ? AND t.NPOLIZA = ? 
    ON DUPLICATE KEY UPDATE
        PERSON_INSURED_TYPE_ID = VALUES(PERSON_INSURED_TYPE_ID);
    `

	_, err := db.Exec(query, ctx.ContractID, ctx.Npolori, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando en PERSON_INSURED para Ramo: %s, Npoliza: %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	log.Printf("Datos insertados correctamente en PERSON_INSURED para Ramo: %s, Npoliza: %s", ctx.Ramo, ctx.Npoliza)
	return nil
}
