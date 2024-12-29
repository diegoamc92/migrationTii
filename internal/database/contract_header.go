package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Create TempIssuanceDates crea una tabla temporal con las fechas de emisión.
func CreateTempIssuanceDates(db *sql.Tx) error {
	query := `
	CREATE TEMPORARY TABLE temp_issuance_dates AS
	SELECT
		NPOLORI,
		MIN(FINIVIG) AS CONTRACT_ISSUANCE_DATE
	FROM temp_csv_polizas
	WHERE CODESTADO = '03' -- Solo "EN VIGOR"
	  AND NPOLIZA LIKE '%00' -- Identificar NPOLIZA raíz
	GROUP BY NPOLORI;
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creando temp_issuance_dates: %v", err)
	}

	fmt.Println("Tabla temporal temp_issuance_dates creada correctamente.")
	log.Println(query)
	return nil
}

// Insert ContractHeader inserta datos en CONTRACT_HEADER usando la tabla temporal temp_issuance_dates.
func InsertContractHeader(db *sql.Tx) ([]int64, error) {
	query := `
	INSERT INTO CONTRACT_HEADER (
		AGENCY_ID,
		INSURER_ID,
		PAYMENT_PLAN_ID,
		PAYMENT_TERM_ID,
		INSURED_PARTY_ID,
		HOLDER_PARTY_ID,
		SECTION_ID,
		SUB_SECTION,
		COVERAGE_PLAN_ID,
		CONTRACT_FROM,
		CONTRACT_TO,
		CURRENCY_ID,
		CONTRACT_ISSUANCE_DATE
	)
	SELECT
		1 AS AGENCY_ID,
		1020 AS INSURER_ID,
		CASE t.IDPERIODPAGO
			WHEN '004' THEN 1000
			WHEN '003' THEN 2000
			WHEN '002' THEN 3000
			WHEN '001' THEN 4000
			ELSE 5000
		END AS PAYMENT_PLAN_ID,
		pt.PAYMENT_TERM_ID,
		p.PARTY_ID AS INSURED_PARTY_ID,
		p.PARTY_ID AS HOLDER_PARTY_ID,
		101 AS SECTION_ID,
		3000 AS SUB_SECTION,
		7 AS COVERAGE_PLAN_ID,
		MIN(t.FINIVIG) AS CONTRACT_FROM,
		MAX(t.FTERVIG) AS CONTRACT_TO,
		4000 AS CURRENCY_ID,
		i.CONTRACT_ISSUANCE_DATE
	FROM temp_csv_polizas t
	JOIN temp_csv_asegurados a ON t.RAMO = a.RAMO AND t.NPOLIZA = a.NPOLIZA
	JOIN PARTY p ON p.EMAIL = a.EMAIL
	JOIN PAYMENT_TERM pt ON pt.PARTY_ID = p.PARTY_ID AND pt.ACCOUNT_NBR = t.NROCONDCOBRO
	LEFT JOIN temp_issuance_dates i ON i.NPOLORI = t.NPOLIZA
	WHERE t.CODESTADO = '03'
	GROUP BY t.NPOLORI;
	`

	res, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("error insertando en CONTRACT_HEADER: %v", err)
	}

	// Recuperar los IDs generados
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	selectQuery := `
	SELECT CONTRACT_ID
	FROM CONTRACT_HEADER
	ORDER BY CONTRACT_ID DESC
	LIMIT ?;
	`

	rows, err := db.Query(selectQuery, rowsAffected)
	if err != nil {
		return nil, fmt.Errorf("error recuperando IDs generados: %v", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error leyendo ID de CONTRACT_HEADER: %v", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
