package database

import (
	"database/sql"
	"fmt"
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
func InsertSingleRequestForContract(db *sql.Tx, contractID int64) (int64, error) {
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
        ?, 1, NULL, NULL, NULL, NULL, NULL, 
        '000-1111111111', NULL, 13000, 
        'MIGRACION TII', 'MIGRACION TII', 
        NOW(), NULL, NULL, 23869
    );
    `

	res, err := db.Exec(query, contractID)
	if err != nil {
		return 0, fmt.Errorf("error insertando REQUEST para CONTRACT_ID %d: %v", contractID, err)
	}

	requestID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo REQUEST_ID: %v", err)
	}

	return requestID, nil
}
