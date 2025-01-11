package database

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateTempIssuanceDates(db *sql.Tx, ctx MigrationContext) error {
	query := `
	CREATE TEMPORARY TABLE temp_issuance_dates AS
	SELECT
		NPOLORI,
		MIN(FINIVIG) AS CONTRACT_ISSUANCE_DATE
	FROM temp_csv_polizas
	WHERE CODESTADO = '03' -- Solo "EN VIGOR"
	  AND RAMO = ?          -- Filtrar por el ramo actual del contexto
	  AND NPOLORI = ?       -- Filtrar por la póliza original del contexto
	  AND NPOLIZA LIKE '%00' -- Identificar NPOLIZA raíz
	GROUP BY NPOLORI;
	`

    _, err := db.Exec(query, ctx.Ramo, ctx.Npolori)
	if err != nil {
        return fmt.Errorf("error creando temp_issuance_dates para RAMO %s, NPOLORI %s: %v", ctx.Ramo, ctx.Npolori, err)
	}

	fmt.Printf("Tabla temporal temp_issuance_dates creada correctamente para RAMO %s, NPOLORI %s.\n", ctx.Ramo, ctx.Npolori)
	log.Println(query)
	return nil
}

func InsertContractHeaderForPoliza(db *sql.Tx, ctx MigrationContext) (int64, error) {
	// Validar que los datos clave no sean vacíos
	if ctx.Ramo == "" || ctx.Npoliza == "" {
		return 0, fmt.Errorf("RAMO o NPOLIZA no pueden estar vacíos")
	}

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
    WHERE t.RAMO = ? AND t.NPOLIZA = ?
    GROUP BY t.NPOLIZA;
    `

	res, err := db.Exec(query, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return 0, fmt.Errorf("error insertando CONTRACT_HEADER para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	contractID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo CONTRACT_ID para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	if contractID == 0 {
		return 0, fmt.Errorf("no se generó un nuevo CONTRACT_ID para RAMO %s, NPOLIZA %s", ctx.Ramo, ctx.Npoliza)
	}

	log.Printf("CONTRACT_HEADER insertado correctamente para RAMO %s, NPOLIZA %s con CONTRACT_ID %d", ctx.Ramo, ctx.Npoliza, contractID)

	return contractID, nil
}
