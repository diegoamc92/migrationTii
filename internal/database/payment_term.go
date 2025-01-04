package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Insert PaymentTerm inserta datos en PAYMENT_TERM asociados a PARTY_ID.
//
//	func InsertPaymentTerm(db *sql.Tx) error {
//		query := `
//		INSERT INTO PAYMENT_TERM (
//	   PARTY_ID, PAYMENT_TYPE_ID, CREDIT_CARD_ID, CURRENCY_ID, ACCOUNT_NBR,
//	   INTER_ACCOUNT_NBR, EXPIRATION, PAYMENT_KEY, REQUIRES_RECEIPT, BANK_ID,
//	   BANK_BRANCH_ID, PERIOD_ID, FIRSTNAME, LASTNAME, SECOND_LASTNAME, RUT,
//	   ACCOUNT_TYPE_ID, COMMUNE_ID, EMPLOYEER_ADDRESS, EMPLOYEER_CONTACT_PHONE,
//	   EMPLOYEER_CONTACT_NAME, ORGANIZATION_ACTIVITY, FIRST_PRIME_DPP_PAYMENT
//
// )
// SELECT
//
//	p.PARTY_ID,                            -- Relacionar con el PARTY_ID correcto
//	CASE t.DESCTPCONDCOBRO
//	    WHEN 'CARGO A CUENTA' THEN 2000    -- PAC
//	    WHEN 'CARGO A TARJETA' THEN 3000  -- PAT
//	    WHEN 'COBRO DIRECTO' THEN 1000    -- DIRECTO
//	    ELSE NULL
//	END AS PAYMENT_TYPE_ID,
//	CASE
//	    WHEN t.DESCTPCONDCOBRO = 'CARGO A TARJETA' THEN 1 -- Simula ID de tarjeta
//	    ELSE NULL
//	END AS CREDIT_CARD_ID,
//	3000 AS CURRENCY_ID,                   -- CLP
//	t.NROCONDCOBRO AS ACCOUNT_NBR,         -- Número de cuenta o tarjeta
//	NULL AS INTER_ACCOUNT_NBR,
//	CASE
//	    WHEN t.DESCTPCONDCOBRO = 'CARGO A TARJETA' THEN '2028-10-01 00:00:00'
//	    ELSE NULL
//	END AS EXPIRATION,
//	NULL AS PAYMENT_KEY,
//	NULL AS REQUIRES_RECEIPT,
//	CASE t.IDCONDCOBRO
//	    WHEN 'CHILE' THEN 1
//	    WHEN 'INTER' THEN 2
//	    WHEN 'ESTADO' THEN 3
//	    WHEN 'SCOTIA' THEN 4
//	    WHEN 'CREDIT' THEN 5
//	    WHEN 'CORP' THEN 6
//	    WHEN 'BICE' THEN 7
//	    WHEN 'SANTAN' THEN 8
//	    WHEN 'ITAU' THEN 9
//	    WHEN 'SECUR' THEN 10
//	    WHEN 'FALABE' THEN 11
//	    WHEN 'BBVA' THEN 13
//	    ELSE NULL
//	END AS BANK_ID,
//	1 AS BANK_BRANCH_ID,                   -- Valor por defecto
//	CASE t.IDPERIODPAGO
//	    WHEN '004' THEN 1000               -- MENSUAL
//	    WHEN '005' THEN 2000               -- TRIMESTRAL
//	    WHEN '006' THEN 3000               -- SEMESTRAL
//	    WHEN '007' THEN 4000               -- ANUAL
//	    WHEN '008' THEN 5000               -- PRIMA ÚNICA
//	    ELSE NULL
//	END AS PERIOD_ID,
//	NULL AS FIRSTNAME,
//	NULL AS LASTNAME,
//	NULL AS SECOND_LASTNAME,
//	NULL AS RUT,
//	NULL AS ACCOUNT_TYPE_ID,
//	NULL AS COMMUNE_ID,
//	NULL AS EMPLOYEER_ADDRESS,
//	NULL AS EMPLOYEER_CONTACT_PHONE,
//	NULL AS EMPLOYEER_CONTACT_NAME,
//	NULL AS ORGANIZATION_ACTIVITY,
//	0 AS FIRST_PRIME_DPP_PAYMENT
//
// FROM temp_csv_polizas t
// JOIN temp_csv_asegurados a
//
//	ON t.RAMO = a.RAMO AND t.NPOLIZA = a.NPOLIZA -- Relación entre las tablas temporales
//
// JOIN PARTY p
//
//	ON p.EMAIL = a.EMAIL                        -- Relación directa con el PARTY único
//
// WHERE t.CODESTADO = '03'
// GROUP BY p.PARTY_ID;
// `
//
//		_, err := db.Exec(query)
//		if err != nil {
//			return fmt.Errorf("error insertando en PAYMENT_TERM: %v", err)
//		}
//
//		fmt.Println("Datos insertados correctamente en PAYMENT_TERM.")
//		log.Println()
//		return nil
//	}
func InsertPaymentTerm(db *sql.Tx, ctx MigrationContext) error {
	query := `
	INSERT INTO PAYMENT_TERM (
		PARTY_ID, PAYMENT_TYPE_ID, CREDIT_CARD_ID, CURRENCY_ID, ACCOUNT_NBR,
		INTER_ACCOUNT_NBR, EXPIRATION, PAYMENT_KEY, REQUIRES_RECEIPT, BANK_ID,
		BANK_BRANCH_ID, PERIOD_ID, FIRSTNAME, LASTNAME, SECOND_LASTNAME, RUT,
		ACCOUNT_TYPE_ID, COMMUNE_ID, EMPLOYEER_ADDRESS, EMPLOYEER_CONTACT_PHONE,
		EMPLOYEER_CONTACT_NAME, ORGANIZATION_ACTIVITY, FIRST_PRIME_DPP_PAYMENT
	)
	SELECT
		p.PARTY_ID,                            -- Relacionar con el PARTY_ID correcto
		CASE t.DESCTPCONDCOBRO
			WHEN 'CARGO A CUENTA' THEN 2000    -- PAC
			WHEN 'CARGO A TARJETA' THEN 3000  -- PAT
			WHEN 'COBRO DIRECTO' THEN 1000    -- DIRECTO
			ELSE NULL
		END AS PAYMENT_TYPE_ID,
		CASE
			WHEN t.DESCTPCONDCOBRO = 'CARGO A TARJETA' THEN 1 -- Simula ID de tarjeta
			ELSE NULL
		END AS CREDIT_CARD_ID,
		3000 AS CURRENCY_ID,                   -- CLP
		t.NROCONDCOBRO AS ACCOUNT_NBR,         -- Número de cuenta o tarjeta
		NULL AS INTER_ACCOUNT_NBR,
		CASE
			WHEN t.DESCTPCONDCOBRO = 'CARGO A TARJETA' THEN '2028-10-01 00:00:00'
			ELSE NULL
		END AS EXPIRATION,
		NULL AS PAYMENT_KEY,
		NULL AS REQUIRES_RECEIPT,
		CASE t.IDCONDCOBRO
			WHEN 'CHILE' THEN 1
			WHEN 'INTER' THEN 2
			WHEN 'ESTADO' THEN 3
			WHEN 'SCOTIA' THEN 4
			WHEN 'CREDIT' THEN 5
			WHEN 'CORP' THEN 6
			WHEN 'BICE' THEN 7
			WHEN 'SANTAN' THEN 8
			WHEN 'ITAU' THEN 9
			WHEN 'SECUR' THEN 10
			WHEN 'FALABE' THEN 11
			WHEN 'BBVA' THEN 13
			ELSE NULL
		END AS BANK_ID,
		1 AS BANK_BRANCH_ID,                   -- Valor por defecto
		CASE t.IDPERIODPAGO
			WHEN '004' THEN 1000               -- MENSUAL
			WHEN '005' THEN 2000               -- TRIMESTRAL
			WHEN '006' THEN 3000               -- SEMESTRAL
			WHEN '007' THEN 4000               -- ANUAL
			WHEN '008' THEN 5000               -- PRIMA ÚNICA
			ELSE NULL
		END AS PERIOD_ID,
		NULL AS FIRSTNAME,
		NULL AS LASTNAME,
		NULL AS SECOND_LASTNAME,
		NULL AS RUT,
		NULL AS ACCOUNT_TYPE_ID,
		NULL AS COMMUNE_ID,
		NULL AS EMPLOYEER_ADDRESS,
		NULL AS EMPLOYEER_CONTACT_PHONE,
		NULL AS EMPLOYEER_CONTACT_NAME,
		NULL AS ORGANIZATION_ACTIVITY,
		0 AS FIRST_PRIME_DPP_PAYMENT
	FROM temp_csv_polizas t
	JOIN temp_csv_asegurados a 
		ON t.RAMO = a.RAMO AND t.NPOLIZA = a.NPOLIZA -- Relación entre las tablas temporales
	JOIN PARTY p 
		ON p.EMAIL = a.EMAIL                        -- Relación directa con el PARTY único
	WHERE t.CODESTADO = '03'
		AND t.RAMO = ? AND t.NPOLIZA = ?            -- Filtro por contexto
	GROUP BY p.PARTY_ID;
	`

	_, err := db.Exec(query, ctx.Ramo, ctx.Npoliza)
	if err != nil {
		return fmt.Errorf("error insertando en PAYMENT_TERM para RAMO %s, NPOLIZA %s: %v", ctx.Ramo, ctx.Npoliza, err)
	}

	fmt.Printf("Datos insertados correctamente en PAYMENT_TERM para RAMO %s, NPOLIZA %s.\n", ctx.Ramo, ctx.Npoliza)
	log.Println(query)
	return nil
}
