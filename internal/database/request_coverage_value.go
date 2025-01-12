package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertRequestCoverageValue(db *sql.Tx, ctx *MigrationContext) error {
	query := `
	INSERT INTO REQUEST_COVERAGE_VALUE (
		REQUEST_ID,
		IC_INSURER_PARTY_ID,
		IC_SECTION_ID,
		IC_SUB_SECTION_ID,
		IC_COVERAGE_ID,
		INSURED_AMOUNT,
		DEDUCTIBLE,
		INITIAL_COVERAGE_GRACE_PERIOD,
		FUT_PAYR_COVERAGE_GRACE_PERI,
		INIT_COVER_GROUP_GRACE_PERIOD,
		FUT_PAYR_COVER_GR_GRACE_PERIOD,
		MAXIMUM_PERIOD_OF_COVERAGE,
		REINSURED_AMOUNT,
		PREMIUM,
		BASE_PREMIUM,
		TAX_VALUE
	)
	SELECT 
		?,                                -- REQUEST_ID
		1020 AS IC_INSURER_PARTY_ID,      -- ID del asegurador
		101 AS IC_SECTION_ID,             -- Sección fija
		3000 AS IC_SUB_SECTION_ID,        -- Sub-sección fija
		1 AS IC_COVERAGE_ID,              -- ID de cobertura fija
		CAST(REPLACE(c.SUMAASG, ',', '') AS DECIMAL(15, 4)) AS INSURED_AMOUNT, -- Monto asegurado (extraído del CSV)
		NULL AS DEDUCTIBLE,               -- Deductible nulo
		NULL AS INITIAL_COVERAGE_GRACE_PERIOD, -- Periodo de gracia inicial nulo
		NULL AS FUT_PAYR_COVERAGE_GRACE_PERI,  -- Periodo de gracia futura nulo
		NULL AS INIT_COVER_GROUP_GRACE_PERIOD, -- Grupo inicial nulo
		NULL AS FUT_PAYR_COVER_GR_GRACE_PERIOD, -- Grupo futuro nulo
		NULL AS MAXIMUM_PERIOD_OF_COVERAGE,    -- Máximo periodo nulo
		NULL AS REINSURED_AMOUNT,         -- Monto reasegurado fijo
		0.46409736 AS PREMIUM,            -- Prima fija
		0.38999736 AS BASE_PREMIUM,       -- Prima base fija
		0.0741 AS TAX_VALUE               -- Valor de impuesto fijo
	FROM temp_csv_coberturas c
	WHERE c.RAMO = ? AND c.NPOLIZA = ?; -- Filtrar por póliza específica
	`

	_, err := db.Exec(query, ctx.RequestID, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando en REQUEST_COVERAGE_VALUE para REQUEST_ID %d: %v", ctx.RequestID, err)
	}

	log.Printf("Datos insertados en REQUEST_COVERAGE_VALUE para REQUEST_ID: %d", ctx.RequestID)
	return nil
}
