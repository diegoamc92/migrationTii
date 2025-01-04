package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Insert PolicyParameter inserta los parámetros asociados a las pólizas en la tabla POLICY_PARAMETER.
//
//	func InsertPolicyParameter(db *sql.Tx, ramo string, nPolOri string) error {
//		query := `
//		INSERT INTO POLICY_PARAMETER (
//	       INSURER_PARTY_ID,
//	       POLICY_ID,
//	       SECTION_ID,
//	       SUB_SECTION_ID,
//	       ENDORSEMENT_ID,
//	       POLICY_PARAMETER_KEY,
//	       POLICY_PARAMETER_DESC,
//	       POLICY_PARAMETER_VALUE
//	   )
//	   SELECT
//	       1020 AS INSURER_PARTY_ID,                               -- ID del asegurador
//	       CONCAT(?, '-', CAST(? AS UNSIGNED)) AS POLICY_ID,       -- POLICY_ID basado en RAMO y NPOLORI
//	       101 AS SECTION_ID,                                      -- Sección fija
//	       3000 AS SUB_SECTION_ID,                                 -- Sub-sección fija
//	       0 AS ENDORSEMENT_ID,                                    -- Endoso inicial
//	       k.PARAMETER_KEY,                                        -- Clave de parámetro
//	       k.PARAMETER_DESC,                                       -- Descripción del parámetro
//	       CASE
//	           WHEN k.PARAMETER_KEY = 'COVERAGE_KEY' THEN CAST(SUBSTRING(c.CLAVCOB, 2) AS UNSIGNED)  -- Cortar primer cero de CLAVCOB
//	           WHEN k.PARAMETER_KEY = 'CONTRACT_TI' THEN t.REQUEST                                   -- Valor del campo REQUEST del CSV para CONTRACT_TI
//	           ELSE k.PARAMETER_VALUE                                                                -- Valores originales para las demás claves
//	       END AS POLICY_PARAMETER_VALUE
//	   FROM (
//	       -- Subconsulta para los valores del ejemplo
//	       SELECT 'BELONGS_TO_BLACK_LIST' AS PARAMETER_KEY, 'BELONGS_TO_BLACK_LIST' AS PARAMETER_DESC, 'false' AS PARAMETER_VALUE
//	       UNION ALL SELECT 'BLACK_LIST_RESPONSE_CODE', 'BLACK_LIST_RESPONSE_CODE', '0'
//	       UNION ALL SELECT 'COVERAGE_KEY', 'COVERAGE_KEY', ''
//	       UNION ALL SELECT 'CONTRACT_TI', 'CONTRACT_TI', ''
//	       UNION ALL SELECT 'CUMULUS_VALUE', 'CUMULUS_VALUE', '0.0'
//	       UNION ALL SELECT 'DEPENDENTS', 'DEPENDENTS', '0'
//	       UNION ALL SELECT 'HAS_ANOTHER_PENDING_REQUEST', 'HAS_ANOTHER_PENDING_REQUEST', 'false'
//	       UNION ALL SELECT 'HAS_EXTRA_PREMIUM_POLICY', 'HAS_EXTRA_PREMIUM_POLICY', 'false'
//	       UNION ALL SELECT 'HAS_FINANCIAL_RISK', 'HAS_FINANCIAL_RISK', 'false'
//	       UNION ALL SELECT 'HAS_REJECTED_REQUEST', 'HAS_REJECTED_REQUEST', 'false'
//	       UNION ALL SELECT 'HAS_RISKY_ACTIVITY', 'HAS_RISKY_ACTIVITY', 'false'
//	       UNION ALL SELECT 'HISTORIC_RISK', 'HISTORIC_RISK', 'false'
//	       UNION ALL SELECT 'IMC_VALUE', 'IMC_VALUE', '23.437499999999996'
//	       UNION ALL SELECT 'INSURED_AND_HOLDER_RELATIONSHIP_IS_OTHER', 'INSURED_AND_HOLDER_RELATIONSHIP_IS_OTHER', 'MISMO'
//	       UNION ALL SELECT 'INTEGRATED_WITH_SAM', 'INTEGRATED_WITH_SAM', 'true'
//	       UNION ALL SELECT 'INTEGRATION_TII', 'INTEGRATION_TII', 'true'
//	       UNION ALL SELECT 'IS_FOREIGN_PERSON', 'IS_FOREIGN_PERSON', 'false'
//	       UNION ALL SELECT 'REQUIRE_MEDICAL_PROTOCOL', 'REQUIRE_MEDICAL_PROTOCOL', 'false'
//	       UNION ALL SELECT 'RESOURCE_DATA', 'RESOURCE_DATA', 'CXCVZCXVZXCVZXCVZCXVZXCV'
//	       UNION ALL SELECT 'TASK_CREATED_DATE', 'TASK_CREATED_DATE', '1735328530313'
//	       UNION ALL SELECT 'VALID_IMC', 'VALID_IMC', 'true'
//	       UNION ALL SELECT 'VALID_QUESTIONNAIRE', 'VALID_QUESTIONNAIRE', 'true'
//	       UNION ALL SELECT 'VALID_REINSURANCE_AMOUNT_VALIDATION', 'VALID_REINSURANCE_AMOUNT_VALIDATION', 'true'
//	   ) AS k
//	   WHERE NOT EXISTS (
//	       SELECT 1
//	       FROM POLICY_PARAMETER pp
//	       WHERE pp.INSURER_PARTY_ID = 1020
//	         AND pp.POLICY_ID = CONCAT(?, '-', CAST(? AS UNSIGNED))
//	         AND pp.POLICY_PARAMETER_KEY = k.PARAMETER_KEY
//	   )
//	   ON DUPLICATE KEY UPDATE
//	       POLICY_PARAMETER_DESC = VALUES(POLICY_PARAMETER_DESC),
//	       POLICY_PARAMETER_VALUE = VALUES(POLICY_PARAMETER_VALUE);
//		`
//
//		_, err := db.Exec(query, ramo, nPolOri, ramo, nPolOri)
//		if err != nil {
//			return fmt.Errorf("error insertando en POLICY_PARAMETER para POLICY_ID %s-%s: %v", ramo, nPolOri, err)
//		}
//
//		log.Printf("Datos insertados (o sobrescritos) correctamente en POLICY_PARAMETER para POLICY_ID: %s-%s", ramo, nPolOri)
//		return nil
//	}
func InsertPolicyParameter(db *sql.Tx, ctx *MigrationContext) error {
	query := `
	INSERT INTO POLICY_PARAMETER (
        INSURER_PARTY_ID,
        POLICY_ID,
        SECTION_ID,
        SUB_SECTION_ID,
        ENDORSEMENT_ID,
        POLICY_PARAMETER_KEY,
        POLICY_PARAMETER_DESC,
        POLICY_PARAMETER_VALUE
    )
    SELECT
        1020 AS INSURER_PARTY_ID,                               -- ID del asegurador
        CONCAT(?, '-', CAST(? AS UNSIGNED)) AS POLICY_ID,       -- POLICY_ID basado en RAMO y NPOLORI
        101 AS SECTION_ID,                                      -- Sección fija
        3000 AS SUB_SECTION_ID,                                 -- Sub-sección fija
        0 AS ENDORSEMENT_ID,                                    -- Endoso inicial
        k.PARAMETER_KEY,                                        -- Clave de parámetro
        k.PARAMETER_DESC,                                       -- Descripción del parámetro
        CASE
            WHEN k.PARAMETER_KEY = 'COVERAGE_KEY' THEN CAST(SUBSTRING(c.CLAVCOB, 2) AS UNSIGNED)  -- Cortar primer cero de CLAVCOB
            WHEN k.PARAMETER_KEY = 'CONTRACT_TI' THEN t.REQUEST                                   -- Valor del campo REQUEST del CSV para CONTRACT_TI
            ELSE k.PARAMETER_VALUE                                                                -- Valores originales para las demás claves
        END AS POLICY_PARAMETER_VALUE
    FROM (
        -- Subconsulta para los valores del ejemplo
        SELECT 'BELONGS_TO_BLACK_LIST' AS PARAMETER_KEY, 'BELONGS_TO_BLACK_LIST' AS PARAMETER_DESC, 'false' AS PARAMETER_VALUE
        UNION ALL SELECT 'BLACK_LIST_RESPONSE_CODE', 'BLACK_LIST_RESPONSE_CODE', '0'
        UNION ALL SELECT 'COVERAGE_KEY', 'COVERAGE_KEY', ''
        UNION ALL SELECT 'CONTRACT_TI', 'CONTRACT_TI', ''
        UNION ALL SELECT 'CUMULUS_VALUE', 'CUMULUS_VALUE', '0.0'
        UNION ALL SELECT 'DEPENDENTS', 'DEPENDENTS', '0'
        UNION ALL SELECT 'HAS_ANOTHER_PENDING_REQUEST', 'HAS_ANOTHER_PENDING_REQUEST', 'false'
        UNION ALL SELECT 'HAS_EXTRA_PREMIUM_POLICY', 'HAS_EXTRA_PREMIUM_POLICY', 'false'
        UNION ALL SELECT 'HAS_FINANCIAL_RISK', 'HAS_FINANCIAL_RISK', 'false'
        UNION ALL SELECT 'HAS_REJECTED_REQUEST', 'HAS_REJECTED_REQUEST', 'false'
        UNION ALL SELECT 'HAS_RISKY_ACTIVITY', 'HAS_RISKY_ACTIVITY', 'false'
        UNION ALL SELECT 'HISTORIC_RISK', 'HISTORIC_RISK', 'false'
        UNION ALL SELECT 'IMC_VALUE', 'IMC_VALUE', '23.437499999999996'
        UNION ALL SELECT 'INSURED_AND_HOLDER_RELATIONSHIP_IS_OTHER', 'INSURED_AND_HOLDER_RELATIONSHIP_IS_OTHER', 'MISMO'
        UNION ALL SELECT 'INTEGRATED_WITH_SAM', 'INTEGRATED_WITH_SAM', 'true'
        UNION ALL SELECT 'INTEGRATION_TII', 'INTEGRATION_TII', 'true'
        UNION ALL SELECT 'IS_FOREIGN_PERSON', 'IS_FOREIGN_PERSON', 'false'
        UNION ALL SELECT 'REQUIRE_MEDICAL_PROTOCOL', 'REQUIRE_MEDICAL_PROTOCOL', 'false'
        UNION ALL SELECT 'RESOURCE_DATA', 'RESOURCE_DATA', 'CXCVZCXVZXCVZXCVZCXVZXCV'
        UNION ALL SELECT 'TASK_CREATED_DATE', 'TASK_CREATED_DATE', '1735328530313'
        UNION ALL SELECT 'VALID_IMC', 'VALID_IMC', 'true'
        UNION ALL SELECT 'VALID_QUESTIONNAIRE', 'VALID_QUESTIONNAIRE', 'true'
        UNION ALL SELECT 'VALID_REINSURANCE_AMOUNT_VALIDATION', 'VALID_REINSURANCE_AMOUNT_VALIDATION', 'true'
    ) AS k
    LEFT JOIN temp_csv_polizas t ON t.RAMO = ? AND t.NPOLIZA = ?
    LEFT JOIN temp_csv_coberturas c ON t.RAMO = c.RAMO AND t.NPOLIZA = c.NPOLIZA
    WHERE t.CODESTADO = '03' -- Solo pólizas activas
    ON DUPLICATE KEY UPDATE 
        POLICY_PARAMETER_DESC = VALUES(POLICY_PARAMETER_DESC),
        POLICY_PARAMETER_VALUE = VALUES(POLICY_PARAMETER_VALUE);
	`

	_, err := db.Exec(query, ctx.Ramo, ctx.NpolOri, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando en POLICY_PARAMETER para POLICY_ID %s-%s: %v", ctx.Ramo, ctx.NpolOri, err)
	}

	log.Printf("Datos insertados (o sobrescritos) correctamente en POLICY_PARAMETER para POLICY_ID: %s-%s", ctx.Ramo, ctx.NpolOri)
	return nil
}
