package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Insert PolicyCoverageValue inserta datos en la tabla POLICY_COVERAGE_VALUE
func InsertPolicyCoverageValue(db *sql.Tx, ramo string, nPolOri string) error {
	query := `
	INSERT INTO POLICY_COVERAGE_VALUE (
        INSURER_PARTY_ID,
        POLICY_ID,
        SECTION_ID,
        SUB_SECTION_ID,
        ENDORSEMENT_ID,
        IC_INSURER_PARTY_ID,
        IC_SECTION_ID,
        IC_SUB_SECTION_ID,
        IC_COVERAGE_ID,
        INSURED_AMOUNT,
        DEDUCTIBLE,
        INITIAL_COVERAGE_GRACE_PERIOD,
        FUT_PAYROLL_COVER_GRACE_PERIOD,
        INIT_COVER_GROUP_GRACE_PERIOD,
        FUT_PAYROLL_COVER_G_GRACE_PERI,
        MAXIMUM_PERIOD_OF_COVERAGE,
        REINSURED_AMOUNT,
        SALARY_MULTIPLIER,
        PREMIUM,
        BASE_PREMIUM,
        TAX_VALUE
    )
    SELECT
        1020 AS INSURER_PARTY_ID,                            -- Aseguradora fija
        CONCAT(?, '-', CAST(? AS UNSIGNED)) AS POLICY_ID,    -- POLICY_ID basado en RAMO y NPOLORI
        101 AS SECTION_ID,                                   -- Sección fija
        3000 AS SUB_SECTION_ID,                              -- Sub-sección fija
        0 AS ENDORSEMENT_ID,                                 -- Endoso inicial
        1020 AS IC_INSURER_PARTY_ID,                         -- Mismo asegurador
        101 AS IC_SECTION_ID,                                -- Sección fija
        3000 AS IC_SUB_SECTION_ID,                           -- Sub-sección fija
        1 AS IC_COVERAGE_ID,                                 -- ID de cobertura fija
        CAST(REPLACE(c.SUMAASG, ',', '') AS UNSIGNED) AS INSURED_AMOUNT, -- Monto asegurado (extraído del CSV)
        NULL AS DEDUCTIBLE,                                  -- Deductible nulo
        NULL AS INITIAL_COVERAGE_GRACE_PERIOD,               -- Periodo de gracia inicial nulo
        NULL AS FUT_PAYROLL_COVER_GRACE_PERIOD,              -- Periodo de gracia futura nulo
        NULL AS INIT_COVER_GROUP_GRACE_PERIOD,               -- Grupo inicial nulo
        NULL AS FUT_PAYROLL_COVER_G_GRACE_PERI,              -- Grupo futuro nulo
        NULL AS MAXIMUM_PERIOD_OF_COVERAGE,                  -- Máximo periodo nulo
        NULL AS REINSURED_AMOUNT,                            -- Monto reasegurado fijo
        NULL AS SALARY_MULTIPLIER,                           -- Multiplicador de salario nulo
        0.46409736 AS PREMIUM, 
        0.38999736 AS BASE_PREMIUM, 
        0.0741 AS TAX_VALUE
    FROM temp_csv_polizas p
    JOIN temp_csv_coberturas c 
        ON p.RAMO = c.RAMO AND p.NPOLIZA = c.NPOLIZA         -- Relación entre tablas temporales
    WHERE p.NPOLIZA LIKE '%00'                              -- Filtro para pólizas originales
    ON DUPLICATE KEY UPDATE
        INSURED_AMOUNT = VALUES(INSURED_AMOUNT),
        PREMIUM = VALUES(PREMIUM),
        BASE_PREMIUM = VALUES(BASE_PREMIUM),
        TAX_VALUE = VALUES(TAX_VALUE);
	`

	_, err := db.Exec(query, ramo, nPolOri)
	if err != nil {
		return fmt.Errorf("error insertando en POLICY_COVERAGE_VALUE para POLICY_ID %s-%s: %v", ramo, nPolOri, err)
	}

	log.Printf("Datos insertados correctamente en POLICY_COVERAGE_VALUE para POLICY_ID: %s-%s", ramo, nPolOri)
	return nil
}
