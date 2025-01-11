package database

import (
	"database/sql"
	"fmt"
	"log"
)


func InsertIntoPolicy(db *sql.Tx, ctx *MigrationContext) error {
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
        CONCAT(?, '-', CAST(? AS UNSIGNED)) AS POLICY_ID,      -- POLICY_ID basado en NPOLORI
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
    GROUP BY RAMO, NPOLORI
    ON DUPLICATE KEY UPDATE
        INSURER_PARTY_ID = VALUES(INSURER_PARTY_ID),
        SECTION_ID = VALUES(SECTION_ID),
        SUB_SECTION_ID = VALUES(SUB_SECTION_ID),
        ENDORSEMENT_ID = VALUES(ENDORSEMENT_ID),
        ENDORSEMENT_TYPE_ID = VALUES(ENDORSEMENT_TYPE_ID),
        POLICY_STATUS_ID = VALUES(POLICY_STATUS_ID),
        CONTRACT_ID = VALUES(CONTRACT_ID),
        POLICY_ISSUANCE_DATE = VALUES(POLICY_ISSUANCE_DATE),
        ENDORSEMENT_DATE = VALUES(ENDORSEMENT_DATE),
        POLICY_ENTRANCE_DATE = VALUES(POLICY_ENTRANCE_DATE),
        POLICY_ENDORSEMENT_DATE_FROM = VALUES(POLICY_ENDORSEMENT_DATE_FROM),
        POLICY_ENDORSEMENT_DATE_TO = VALUES(POLICY_ENDORSEMENT_DATE_TO),
        DATE_MODIFIED = NOW(),
        AGENT_PARTY_ID = VALUES(AGENT_PARTY_ID);
    `

	_, err := db.Exec(query, ctx.Ramo, ctx.Npolori, ctx.ContractID, ctx.Ramo, ctx.Npolori)
	if err != nil {
		return fmt.Errorf("error insertando en POLICY para CONTRACT_ID %d: %v", ctx.ContractID, err)
	}

	log.Printf("Datos insertados correctamente en POLICY para CONTRACT_ID: %d, RAMO: %s, NPOLORI: %s", ctx.ContractID, ctx.Ramo, ctx.Npolori)
	return nil
}
