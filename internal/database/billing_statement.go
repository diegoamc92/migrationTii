package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Insert BillingStatement inserta datos en la tabla BILLING_STATEMENT.
//
//	func InsertBillingStatement(db *sql.Tx) error {
//		query := `
//	   INSERT INTO BILLING_STATEMENT (
//	   BILLING_STATEMENT_ID,
//	   BILL_TYPE_ID,
//	   BILL_DATE_FROM,
//	   BILL_DATE_TO,
//	   PAYMENT_TERM_ID,
//	   BILL_AMOUNT,
//	   BILLING_STATEMENT_SEQ,
//	   INSURER_PARTY_ID,
//	   POLICY_ID,
//	   SECTION_ID,
//	   SUB_SECTION_ID,
//	   ENDORSEMENT_ID,
//	   OTHER_AMOUNT,
//	   BILLING_STATEMENT_MAX_SEQ,
//	   BILLING_CHANNEL_TYPE_ID,
//	   ITEM_ID,
//	   CONTRACT_ID,
//	   BILLING_NUMBER,
//	   BILLING_STATEMENT_STATUS_ID,
//	   SUB_TYPE_ID,
//	   CLIENT_CODE,
//	   COMMENTS
//
// )
// SELECT
//
//	NULL AS BILLING_STATEMENT_ID, -- Generado automáticamente
//	1000 AS BILL_TYPE_ID,         -- Tipo de factura fija
//	t.FINIVIG AS BILL_DATE_FROM,  -- Fecha desde (FINIVIG del CSV)
//	t.FTERVIG AS BILL_DATE_TO,    -- Fecha hasta (FTERVIG del CSV)
//	pt.PAYMENT_TERM_ID,           -- ID de la condición de pago
//	0.4641 AS BILL_AMOUNT,        -- Monto de la factura
//	NULL AS BILLING_STATEMENT_SEQ,
//	1020 AS INSURER_PARTY_ID,     -- Aseguradora fija
//	CONCAT(t.RAMO, '-', t.NPOLORI) AS POLICY_ID, -- ID de la póliza generado
//	101 AS SECTION_ID,            -- Sección fija
//	3000 AS SUB_SECTION_ID,       -- Sub-sección fija
//	0 AS ENDORSEMENT_ID,          -- Endoso inicial
//	NULL AS OTHER_AMOUNT,         -- Monto adicional
//	NULL AS BILLING_STATEMENT_MAX_SEQ,
//	NULL AS BILLING_CHANNEL_TYPE_ID,
//	NULL AS ITEM_ID,
//	c.CONTRACT_ID,                -- ID del contrato asociado
//	NULL AS BILLING_NUMBER,       -- Número de factura
//	5000 AS BILLING_STATEMENT_STATUS_ID, -- Estado de la factura
//	NULL AS SUB_TYPE_ID,
//	NULL AS CLIENT_CODE,
//	NULL AS COMMENTS
//
// FROM temp_csv_polizas t
//
//	JOIN PAYMENT_TERM pt ON pt.ACCOUNT_NBR = t.NROCONDCOBRO -- Asociar con Payment Term
//	JOIN CONTRACT_HEADER c ON c.CONTRACT_ID = pt.PAYMENT_TERM_ID -- Asociar con el contrato
//
// WHERE t.CODESTADO = '03';
//
//	   `
//
//		_, err := db.Exec(query)
//		if err != nil {
//			return fmt.Errorf("error insertando en BILLING_STATEMENT: %v", err)
//		}
//
//		fmt.Println("Datos insertados correctamente en BILLING_STATEMENT.")
//		log.Println(query)
//		return nil
//	}
func InsertBillingStatement(db *sql.Tx, ctx *MigrationContext) error {
	query := `
    INSERT INTO BILLING_STATEMENT (
        BILL_TYPE_ID,
        BILL_DATE_FROM,
        BILL_DATE_TO,
        PAYMENT_TERM_ID,
        BILL_AMOUNT,
        INSURER_PARTY_ID,
        POLICY_ID,
        SECTION_ID,
        SUB_SECTION_ID,
        ENDORSEMENT_ID,
        CONTRACT_ID,
        BILLING_STATEMENT_STATUS_ID
    )
    SELECT
        1000 AS BILL_TYPE_ID,         -- Tipo de factura fija
        t.FINIVIG AS BILL_DATE_FROM,  -- Fecha desde (FINIVIG del CSV)
        t.FTERVIG AS BILL_DATE_TO,    -- Fecha hasta (FTERVIG del CSV)
        ch.PAYMENT_TERM_ID,           -- PAYMENT_TERM_ID desde CONTRACT_HEADER
        0.4641 AS BILL_AMOUNT,        -- Monto de la factura
        1020 AS INSURER_PARTY_ID,     -- Aseguradora fija
        CONCAT(t.RAMO, '-', t.NPOLIZA) AS POLICY_ID, -- ID de la póliza generado
        101 AS SECTION_ID,            -- Sección fija
        3000 AS SUB_SECTION_ID,       -- Sub-sección fija
        0 AS ENDORSEMENT_ID,          -- Endoso inicial
        ch.CONTRACT_ID,               -- ID del contrato desde CONTRACT_HEADER
        5000 AS BILLING_STATEMENT_STATUS_ID -- Estado de la factura
    FROM temp_csv_polizas t
    JOIN CONTRACT_HEADER ch 
        ON ch.CONTRACT_ID = ?         -- Asociar directamente usando ContractID del contexto
    WHERE t.RAMO = ? AND t.NPOLIZA = ?; -- Filtrar por RAMO y NPOLIZA desde el contexto
    `

	_, err := db.Exec(query, ctx.ContractID, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando en BILLING_STATEMENT para CONTRACT_ID %d: %v", ctx.ContractID, err)
	}

	log.Printf("Datos insertados correctamente en BILLING_STATEMENT para CONTRACT_ID: %d, RAMO: %s, NPOLIZA: %s", ctx.ContractID, ctx.Ramo, ctx.Npoliza)
	return nil
}
