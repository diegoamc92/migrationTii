package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertSingleRequestForContract(db *sql.Tx, ctx *MigrationContext) (int64, error) {
	query := `
    INSERT INTO REQUEST (
        CONTRACT_ID,
        PENDING_INSPECTION,
        INSURER_ID,
        POLICY_ID,
        SECTION_ID,
        SUB_SECTION_ID,
        ENDORSEMENT_ID,
        CERTIFICATION_NUMBER,
        USER_ID,
        REQUEST_STATUS_ID,
        OBSERVATIONS,
        COMMENTS,
        CREATED_DATE,
        DUE_DATE,
        ACCOUNT_ID,
        AGENT_PARTY_ID
    )
    VALUES (
        ?,                        -- CONTRACT_ID
        1,                        -- PENDING_INSPECTION
        1020,                     -- INSURER_ID
        CONCAT(?, '-', CAST(? AS UNSIGNED)), -- POLICY_ID basado en RAMO y NPOLORI
        101,                      -- SECTION_ID
        3000,                     -- SUB_SECTION_ID
        0,                        -- ENDORSEMENT_ID
        '000-1111111111',         -- CERTIFICATION_NUMBER
        NULL,                     -- USER_ID
        13000,                    -- REQUEST_STATUS_ID
        'MIGRACION TII',          -- OBSERVATIONS
        'MIGRACION TII',          -- COMMENTS
        NOW(),                    -- CREATED_DATE
        NULL,                     -- DUE_DATE
        NULL,                     -- ACCOUNT_ID
        23869                     -- AGENT_PARTY_ID
    );
    `

	res, err := db.Exec(query, ctx.ContractID, ctx.Ramo, ctx.Npolori)
	if err != nil {
		return 0, fmt.Errorf("error insertando REQUEST para CONTRACT_ID %d: %v", ctx.ContractID, err)
	}

	requestID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo REQUEST_ID: %v", err)
	}

	log.Printf("REQUEST insertado correctamente para CONTRACT_ID %d con REQUEST_ID %d", ctx.ContractID, requestID)
	return requestID, nil
}
