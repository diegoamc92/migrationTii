package database

import (
	"database/sql"
	"fmt"
	"log"
)

//func InsertIntoPolicy(db *sql.Tx, contractID int64) error {
//	query := `
//    INSERT INTO POLICY (
//        INSURER_PARTY_ID,
//        POLICY_ID,
//        SECTION_ID,
//        SUB_SECTION_ID,
//        ENDORSEMENT_ID,
//        ENDORSEMENT_TYPE_ID,
//        POLICY_STATUS_ID,
//        CONTRACT_ID,
//        POLICY_ISSUANCE_DATE,
//        ENDORSEMENT_DATE,
//        POLICY_ENTRANCE_DATE,
//        POLICY_PHISICAL_DATE_DELIVERY,
//        POLICY_PHISICAL_DATE_RECEPTION,
//        POLICY_ENDORSEMENT_DATE_FROM,
//        POLICY_ENDORSEMENT_DATE_TO,
//        DATA_SOURCE,
//        POLICY_ELECTRONIC_ACCEPTED,
//        RENEWED_BY,
//        RENEWED_NUMBER,
//        POLICY_AGREEMENT_NUMBER,
//        GRACE_PERIOD,
//        DATE_MODIFIED,
//        POLICY_AFFINITY_GROUP_ID,
//        UNITED_PREMIUM,
//        AGENT_PARTY_ID
//    )
//    SELECT
//        1020 AS INSURER_PARTY_ID,                      -- Aseguradora fija
//        CONCAT(t.RAMO, '-', CAST(t.NPOLORI AS UNSIGNED)) AS POLICY_ID,  -- Número de póliza original
//        101 AS SECTION_ID,                            -- Sección fija
//        3000 AS SUB_SECTION_ID,                       -- Sub-sección fija
//        0 AS ENDORSEMENT_ID,                          -- Endoso inicial
//        3000 AS ENDORSEMENT_TYPE_ID,                  -- Tipo de endoso: Nueva póliza
//        1000 AS POLICY_STATUS_ID,                     -- Estado: Vigente
//        ? AS CONTRACT_ID,                             -- ID del contrato asociado
//        MIN(t.FINIVIG) AS POLICY_ISSUANCE_DATE,       -- Fecha de emisión original
//        MIN(t.FINIVIG) AS ENDORSEMENT_DATE,           -- Fecha de endoso inicial
//        MIN(t.FINIVIG) AS POLICY_ENTRANCE_DATE,       -- Fecha de entrada
//        NULL AS POLICY_PHISICAL_DATE_DELIVERY,        -- Mantener NULL
//        NULL AS POLICY_PHISICAL_DATE_RECEPTION,       -- Mantener NULL
//        MIN(t.FINIVIG) AS POLICY_ENDORSEMENT_DATE_FROM, -- Fecha de inicio de vigencia
//        MAX(t.FTERVIG) AS POLICY_ENDORSEMENT_DATE_TO,   -- Fecha de término de vigencia
//        NULL AS DATA_SOURCE,                          -- Fuente de datos
//        0 AS POLICY_ELECTRONIC_ACCEPTED,              -- No aceptado electrónicamente
//        NULL AS RENEWED_BY,                           -- No renovada por nadie
//        NULL AS RENEWED_NUMBER,                       -- No tiene número de renovación
//        NULL AS POLICY_AGREEMENT_NUMBER,              -- Sin número de acuerdo
//        NULL AS GRACE_PERIOD,                         -- Sin período de gracia
//        NOW() AS DATE_MODIFIED,                       -- Fecha de modificación actual
//        NULL AS POLICY_AFFINITY_GROUP_ID,             -- Sin grupo de afinidad
//        NULL AS UNITED_PREMIUM,                       -- Prima unificada nula
//        23869 AS AGENT_PARTY_ID                       -- ID del agente fijo
//    FROM temp_csv_polizas t
//    JOIN CONTRACT_HEADER ch ON ch.CONTRACT_ID = ?
//    WHERE t.CODESTADO = '03'
//    GROUP BY t.RAMO, t.NPOLORI;
//    `
//
//	_, err := db.Exec(query, contractID)
//	if err != nil {
//		return fmt.Errorf("error insertando en POLICY para CONTRACT_ID %d: %v", contractID, err)
//	}
//
//	log.Printf("Datos insertados correctamente en POLICY para CONTRACT_ID: %d", contractID)
//	return nil
//}

func InsertIntoPolicy(db *sql.Tx, contractID int64, ramo string, nPolOri string) error {
	query := `
    INSERT INTO POLICY (
        INSURER_PARTY_ID,
        POLICY_ID,
        SECTION_ID,
        SUB_SECTION_ID,
        ENDORSEMENT_ID,
        ENDORSEMENT_TYPE_ID,
        POLICY_STATUS_ID,
        CONTRACT_ID,
        POLICY_ISSUANCE_DATE,
        ENDORSEMENT_DATE,
        POLICY_ENTRANCE_DATE,
        POLICY_PHISICAL_DATE_DELIVERY,
        POLICY_PHISICAL_DATE_RECEPTION,
        POLICY_ENDORSEMENT_DATE_FROM,
        POLICY_ENDORSEMENT_DATE_TO,
        DATA_SOURCE,
        POLICY_ELECTRONIC_ACCEPTED,
        RENEWED_BY,
        RENEWED_NUMBER,
        POLICY_AGREEMENT_NUMBER,
        GRACE_PERIOD,
        DATE_MODIFIED,
        POLICY_AFFINITY_GROUP_ID,
        UNITED_PREMIUM,
        AGENT_PARTY_ID
    )
    SELECT
        1020 AS INSURER_PARTY_ID,                              -- Aseguradora fija
        CONCAT(?, '-', CAST(? AS UNSIGNED)) AS POLICY_ID, -- POLICY_ID basado en NPOLORI
        101 AS SECTION_ID,                                     -- Sección fija
        3000 AS SUB_SECTION_ID,                                -- Sub-sección fija
        0 AS ENDORSEMENT_ID,                                   -- Endoso inicial
        3000 AS ENDORSEMENT_TYPE_ID,                           -- Tipo de endoso: Nueva póliza
        1000 AS POLICY_STATUS_ID,                              -- Estado: Vigente
        ? AS CONTRACT_ID,                                      -- ID del contrato asociado
        MIN(FINIVIG) AS POLICY_ISSUANCE_DATE,                  -- Fecha de emisión original
        MIN(FINIVIG) AS ENDORSEMENT_DATE,                      -- Fecha de endoso inicial
        MIN(FINIVIG) AS POLICY_ENTRANCE_DATE,                  -- Fecha de entrada
        NULL AS POLICY_PHISICAL_DATE_DELIVERY,                 -- Mantener NULL
        NULL AS POLICY_PHISICAL_DATE_RECEPTION,                -- Mantener NULL
        MIN(FINIVIG) AS POLICY_ENDORSEMENT_DATE_FROM,          -- Fecha de inicio de vigencia
        MIN(FTERVIG) AS POLICY_ENDORSEMENT_DATE_TO,            -- Fecha de término de vigencia
        NULL AS DATA_SOURCE,                                   -- Fuente de datos
        0 AS POLICY_ELECTRONIC_ACCEPTED,                       -- No aceptado electrónicamente
        NULL AS RENEWED_BY,                                    -- No renovada por nadie
        NULL AS RENEWED_NUMBER,                                -- No tiene número de renovación
        NULL AS POLICY_AGREEMENT_NUMBER,                       -- Sin número de acuerdo
        NULL AS GRACE_PERIOD,                                  -- Sin período de gracia
        NOW() AS DATE_MODIFIED,                                -- Fecha de modificación actual
        NULL AS POLICY_AFFINITY_GROUP_ID,                      -- Sin grupo de afinidad
        NULL AS UNITED_PREMIUM,                                -- Prima unificada nula
        23869 AS AGENT_PARTY_ID                                -- ID del agente fijo
    FROM temp_csv_polizas
    WHERE RAMO = ? AND NPOLORI = ?                            -- Tomar los datos del NPOLORI (original)
    GROUP BY RAMO, NPOLORI;
    `

	_, err := db.Exec(query, ramo, nPolOri, contractID, ramo, nPolOri)
	if err != nil {
		return fmt.Errorf("error insertando en POLICY para CONTRACT_ID %d: %v", contractID, err)
	}

	log.Printf("Datos insertados correctamente en POLICY para CONTRACT_ID: %d, RAMO: %s, NPOLORI: %s", contractID, ramo, nPolOri)
	return nil
}
