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
func InsertRequestCoverageValue(db *sql.Tx, requestID int64) error {
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
	VALUES (?, 1020, 101, 3000, 1, 25000, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 0.46409736, 0.38999736, 0.0741);
	`

	_, err := db.Exec(query, requestID)
	if err != nil {
		return fmt.Errorf("error insertando en REQUEST_COVERAGE_VALUE para REQUEST_ID %d: %v", requestID, err)
	}

	log.Printf("Datos insertados en REQUEST_COVERAGE_VALUE para REQUEST_ID: %d", requestID)
	return nil
}
