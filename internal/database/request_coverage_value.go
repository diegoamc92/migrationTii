package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func GetRequestIDs(db *sql.Tx, contractIDs []int64) ([]int64, error) {
	query := `
	SELECT REQUEST_ID
	FROM REQUEST
	WHERE CONTRACT_ID IN (?` + strings.Repeat(",?", len(contractIDs)-1) + `)
	`
	args := make([]interface{}, len(contractIDs))
	for i, id := range contractIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error recuperando REQUEST_IDs: %v", err)
	}
	defer rows.Close()

	var requestIDs []int64
	for rows.Next() {
		var requestID int64
		if err := rows.Scan(&requestID); err != nil {
			return nil, fmt.Errorf("error leyendo REQUEST_ID: %v", err)
		}
		requestIDs = append(requestIDs, requestID)
	}

	return requestIDs, nil
}

// RequestCoverageValue inserta valores en la tabla REQUEST_COVERAGE_VALUE.
//func InsertRequestCoverageValue(db *sql.Tx, requestID int64) error {
//	query := `
//	INSERT INTO REQUEST_COVERAGE_VALUE (
//		REQUEST_ID,
//		IC_INSURER_PARTY_ID,
//		IC_SECTION_ID,
//		IC_SUB_SECTION_ID,
//		IC_COVERAGE_ID,
//		INSURED_AMOUNT,
//		DEDUCTIBLE,
//		INITIAL_COVERAGE_GRACE_PERIOD,
//		FUT_PAYR_COVERAGE_GRACE_PERI,
//		INIT_COVER_GROUP_GRACE_PERIOD,
//		FUT_PAYR_COVER_GR_GRACE_PERIOD,
//		MAXIMUM_PERIOD_OF_COVERAGE,
//		REINSURED_AMOUNT,
//		PREMIUM,
//		BASE_PREMIUM,
//		TAX_VALUE
//	)
//	VALUES (?, 1020, 101, 3000, 1, 25000, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0.46409736, 0.38999736, 0.0741);
//	`
//
//	_, err := db.Exec(query, requestID)
//	if err != nil {
//		return fmt.Errorf("error insertando en REQUEST_COVERAGE_VALUE para REQUEST_ID %d: %v", requestID, err)
//	}
//
//	log.Printf("Datos insertados en REQUEST_COVERAGE_VALUE para REQUEST_ID: %d", requestID)
//	return nil
//}

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
