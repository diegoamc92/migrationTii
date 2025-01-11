package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertPartyRelation(db *sql.Tx, ctx *MigrationContext) error {
	query := `
    INSERT INTO PARTY_RELATION (
        PARTY_RELATION_TYPE_ID,
        PARTY_ID,
        PARTY_ID_RELATED
    )
    SELECT
        CASE t.DESDEPEND
            WHEN 'CONYUGE' THEN 1000
            WHEN 'HIJO (A)' THEN 2000
            ELSE 3000 -- Otro tipo de relación
        END AS PARTY_RELATION_TYPE_ID,
        (SELECT p.PARTY_ID FROM PARTY p WHERE p.EMAIL = t_tit.EMAIL) AS PARTY_ID, -- Titular
        p_adherente.PARTY_ID AS PARTY_ID_RELATED -- Adherente
    FROM temp_csv_asegurados t_adherente
    JOIN temp_csv_asegurados t_tit ON t_tit.RAMO = t_adherente.RAMO AND t_tit.NPOLIZA = t_adherente.NPOLIZA AND t_tit.RUT = ?
    JOIN PARTY p_adherente ON p_adherente.EMAIL = t_adherente.EMAIL
    WHERE t_adherente.RAMO = ? AND t_adherente.NPOLIZA = ? AND t_adherente.RUT != ?
    ON DUPLICATE KEY UPDATE
        PARTY_RELATION_TYPE_ID = VALUES(PARTY_RELATION_TYPE_ID);
    `

	_, err := db.Exec(query, ctx.NpolOri, ctx.Ramo, ctx.Npoliza, ctx.NpolOri)
	if err != nil {
		return fmt.Errorf("error insertando en PARTY_RELATION para Ramo: %s, Npoliza: %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	log.Printf("Relaciones insertadas en PARTY_RELATION para Ramo: %s, Npoliza: %s", ctx.Ramo, ctx.Npoliza)
	return nil
}
