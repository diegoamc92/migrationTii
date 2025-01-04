package database

import (
	"database/sql"
	"fmt"
	"log"
)

//func TempContractId(db *sql.Tx) error {
//	query := `
//	CREATE TEMPORARY TABLE temp_contract_mapping AS
//		SELECT
//    	t.NPOLIZA AS NPOLIZA,
//    c.CONTRACT_ID AS CONTRACT_ID
//		FROM temp_csv_polizas t
//			JOIN CONTRACT_HEADER c
//    			ON t.NPOLIZA = c.CONTRACT_ID
//			WHERE t.CODESTADO = '03';
//`
//	_, err := db.Exec(query)
//	if err != nil {
//		return fmt.Errorf("error creando temp_contract_mapping: %v", err)
//	}
//
//	fmt.Println("Tabla temporal temp_contract_mapping creada correctamente.")
//	log.Println(query)
//	return nil
//}

// Insert Request inserta datos en la tabla REQUEST.
//
//	func InsertSingleRequestForContract(db *sql.Tx, contractID int64) (int64, error) {
//		query := `
//	   INSERT INTO REQUEST (
//	       CONTRACT_ID,
//	       PENDING_INSPECTION,
//	       INSURER_ID,
//	       POLICY_ID,
//	       SECTION_ID,
//	       SUB_SECTION_ID,
//	       ENDORSEMENT_ID,
//	       CERTIFICATION_NUMBER,
//	       USER_ID,
//	       REQUEST_STATUS_ID,
//	       OBSERVATIONS,
//	       COMMENTS,
//	       CREATED_DATE,
//	       DUE_DATE,
//	       ACCOUNT_ID,
//	       AGENT_PARTY_ID
//	   )
//	   VALUES (
//	       ?, 1, NULL, NULL, NULL, NULL, NULL,
//	       '000-1111111111', NULL, 13000,
//	       'MIGRACION TII', 'MIGRACION TII',
//	       NOW(), NULL, NULL, 23869
//	   );
//	   `
//
//		res, err := db.Exec(query, contractID)
//		if err != nil {
//			return 0, fmt.Errorf("error insertando REQUEST para CONTRACT_ID %d: %v", contractID, err)
//		}
//
//		requestID, err := res.LastInsertId()
//		if err != nil {
//			return 0, fmt.Errorf("error obteniendo REQUEST_ID: %v", err)
//		}
//
//		return requestID, nil
//	}
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

	res, err := db.Exec(query, ctx.ContractID, ctx.Ramo, ctx.NpolOri)
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
