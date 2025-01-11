package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertPolicyEconomics(db *sql.Tx, ctx *MigrationContext) error {
	query := `
    INSERT INTO POLICY_ECONOMICS (
        ECONOMIC_ITEM_ID,
        ECONOMIC_VALUE,
        INSURER_PARTY_ID,
        POLICY_ID,
        SECTION_ID,
        SUB_SECTION_ID,
        ENDORSEMENT_ID,
        ECONOMIC_VALUE_DATE,
        TAX_ID
    )
    SELECT
        e.ECONOMIC_ITEM_ID,
        CASE
            WHEN e.ECONOMIC_ITEM_ID = 9000 THEN CAST(c.SUMAASG AS DECIMAL(15, 4)) -- SUMA ASEGURADA
            WHEN e.ECONOMIC_ITEM_ID = 29000 THEN CAST(c.PMAANUAL AS DECIMAL(15, 4)) -- PRIMA ANUAL
            WHEN e.ECONOMIC_ITEM_ID = 8000 THEN CAST(c.IVAANUAL AS DECIMAL(15, 4)) -- IVA
            ELSE e.ECONOMIC_VALUE
        END AS ECONOMIC_VALUE,
        1020 AS INSURER_PARTY_ID,             -- ID fijo del asegurador
        CONCAT(?, '-', CAST(? AS UNSIGNED)) AS POLICY_ID, -- Generar POLICY_ID basado en RAMO y NPOLORI
        101 AS SECTION_ID,                    -- Sección fija
        3000 AS SUB_SECTION_ID,               -- Sub-sección fija
        0 AS ENDORSEMENT_ID,                  -- Endoso inicial
        NOW() AS ECONOMIC_VALUE_DATE,         -- Fecha actual
        NULL AS TAX_ID                        -- Sin impuesto específico
    FROM (
        -- Valores base para los ECONOMIC_ITEM_ID
        SELECT 1000 AS ECONOMIC_ITEM_ID, 0 AS ECONOMIC_VALUE
        UNION ALL SELECT 2000, 0.4641
        UNION ALL SELECT 3000, 0
        UNION ALL SELECT 4000, 0
        UNION ALL SELECT 5000, 0
        UNION ALL SELECT 6000, 0.4641
        UNION ALL SELECT 7000, 0
        UNION ALL SELECT 8000, 0.0741
        UNION ALL SELECT 9000, 30000       -- Este será ajustado
        UNION ALL SELECT 10000, 0
        UNION ALL SELECT 11000, 0
        UNION ALL SELECT 12000, 0
        UNION ALL SELECT 13000, 0
        UNION ALL SELECT 14000, 0
        UNION ALL SELECT 15000, 0
        UNION ALL SELECT 16000, 0
        UNION ALL SELECT 17000, 0
        UNION ALL SELECT 18000, 0
        UNION ALL SELECT 19000, 0
        UNION ALL SELECT 20000, 0
        UNION ALL SELECT 21000, 0.39
        UNION ALL SELECT 24000, 0
        UNION ALL SELECT 25000, 0
        UNION ALL SELECT 26000, 0
        UNION ALL SELECT 28000, 0
        UNION ALL SELECT 29000, 5.1056
    ) e
    JOIN temp_csv_coberturas c
        ON c.RAMO = ? AND c.NPOLIZA = ?
    ON DUPLICATE KEY UPDATE
        ECONOMIC_VALUE = VALUES(ECONOMIC_VALUE),
        ECONOMIC_VALUE_DATE = VALUES(ECONOMIC_VALUE_DATE);
	`

	_, err := db.Exec(query, ctx.Ramo, ctx.Npolori, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando en POLICY_ECONOMICS para POLICY_ID %s-%s: %v", ctx.Ramo, ctx.Npolori, err)
	}

	log.Printf("Datos insertados correctamente en POLICY_ECONOMICS para POLICY_ID: %s-%s", ctx.Ramo, ctx.Npolori)
	return nil
}
