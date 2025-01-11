package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertAdherent(db *sql.Tx, ctx *MigrationContext) error {
	query := `
    INSERT INTO ADHERENT (
        ITEM_ID,
        CONTRACT_ID,
        PARTY_ID,
        PARTY_RELATION_TYPE_ID,
        INCORPORATION_DATE,
        INSURED_AMOUNT,
        ADHERENT_STATUS,
        VALIDITY_FROM
    )
    SELECT
        ROW_NUMBER() OVER (PARTITION BY t.NPOLIZA) AS ITEM_ID, -- Número incremental para ITEM_ID
        ? AS CONTRACT_ID,                                     -- CONTRACT_ID de la póliza
        p.PARTY_ID,                                           -- PARTY_ID de la carga
        CASE
            WHEN t.DESDEPEND = 'CONYUGE' THEN 1000
            WHEN t.DESDEPEND = 'HIJO (A)' THEN 2000
            ELSE NULL
        END AS PARTY_RELATION_TYPE_ID,
        NOW() AS INCORPORATION_DATE,                         -- Fecha actual
        CAST(REPLACE(c.SUMAASG, ',', '') AS UNSIGNED) AS INSURED_AMOUNT, -- Monto asegurado
        1 AS ADHERENT_STATUS,                                -- Estado activo
        t.FINIVIG AS VALIDITY_FROM                           -- Fecha de inicio de validez
    FROM temp_csv_asegurados t
    JOIN PARTY p ON p.RUT = TRIM(LEADING '0' FROM t.RUT)     -- Asociar con PARTY por RUT
    JOIN temp_csv_coberturas c ON t.RAMO = c.RAMO AND t.NPOLIZA = c.NPOLIZA -- Relación con coberturas
    WHERE t.RAMO = ? AND t.NPOLIZA = ? AND t.RUT NOT IN (
        SELECT DISTINCT RUT FROM PARTY WHERE PARTY_ID IN (
            SELECT HOLDER_PARTY_ID FROM CONTRACT_HEADER WHERE CONTRACT_ID = ?
        )
    );
    `

	_, err := db.Exec(query, ctx.ContractID, ctx.Ramo, ctx.Npoliza, ctx.ContractID)
	if err != nil {
		return fmt.Errorf("error insertando adherentes para CONTRACT_ID %d: %v", ctx.ContractID, err)
	}

	log.Printf("Adherentes insertados correctamente para CONTRACT_ID: %d", ctx.ContractID)
	return nil
}
