package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration Context representa los datos necesarios para migrar una póliza
type MigrationContext struct {
	Npoliza        string
	Npolori        string
	Ramo           string
	RUT            string
	Nombres        string
	ApeMaterno     string
	ApePaterno     string
	Parentesco     string
	CoberturaTipo  string
	MontoCobertura float64
	Email          string
	Telefono       string
	Direccion      string
	Ciudad         string
	Region         string
	Pais           string
	Premium        float64
	Deductible     float64
	ContractID     int64
	RequestID      int64
	AdherentID     int64
}

// ProcessPoliza maneja la lógica de inserciones para una póliza
func ProcessPoliza(tx *sql.Tx, context MigrationContext) error {
	// Insertar PARTY
	if err := InsertPartyWithContext(tx, context); err != nil {
		return fmt.Errorf("error en PARTY para póliza %s: %v", context.Npoliza, err)
	}

	// Insertar IDENTIFICATION
	if err := InsertIdentificationWithContext(tx); err != nil {
		return fmt.Errorf("error en IDENTIFICATION para póliza %s: %v", context.Npoliza, err)
	}

	// Asociar IDENTIFICATION con PARTY
	if err := AssociatePartyIdentificationWithContext(tx, context); err != nil {
		return fmt.Errorf("error asociando PARTY_IDENTIFICATION: %v", err)
	}

	// Insertar inInsertInsuredData
	if err := InsertInsuredData(tx); err != nil {
		return fmt.Errorf("error insertando INSURED_DATA: %v", err)
	}

	// InsertPolicyHolderData
	if err := InsertPolicyHolderData(tx); err != nil {
		return fmt.Errorf("error insertando POLICY_HOLDER: %v", err)
	}

	// Insertar EMAIL, PHONE, ADDRESS
	if err := InsertEmail(tx, context); err != nil {
		return fmt.Errorf("error insertando EMAIL: %v", err)
	}

	if err := InsertPhone(tx, context); err != nil {
		return fmt.Errorf("error insertando PHONE: %v", err)
	}

	if err := InsertAddress(tx, context); err != nil {
		return fmt.Errorf("error insertando ADDRESS: %v", err)
	}

	if err := AssociatePartyAddress(tx, context); err != nil {
		return fmt.Errorf("error asociando PARTY_ADDRESS: %v", err)
	}

	// Insertar PERSON
	if err := InsertPersonData(tx, context); err != nil {
		return fmt.Errorf("error insertando PERSON: %v", err)
	}

	// Insertar PAYMENT_TERM
	if err := InsertPaymentTerm(tx, context); err != nil {
		return fmt.Errorf("error insertando PAYMENT_TERM: %v", err)
	}

	if err := CreateTempIssuanceDates(tx, context); err != nil {
		return fmt.Errorf("error creando tabla temporal: %v", err)
	}

	// Insertar ADHERENT
	if err := InsertAdherent(tx, &context); err != nil {
		return fmt.Errorf("error insertando ADHERENT: %v", err)
	}

	// Insertar PARTY_RELATION
	if err := InsertPartyRelation(tx, &context); err != nil {
		return fmt.Errorf("error insertando PARTY_RELATION: %v", err)
	}

	// Insertar PERSON_INSURED
	if err := InsertPersonInsured(tx, &context); err != nil {
		return fmt.Errorf("error insertando PERSON_INSURED: %v", err)
	}

	log.Println("Recuperando lista de pólizas desde tiisa.polizas...")
	polizas, err := GetPolizasFromDynamicTable(tx)
	if err != nil {
		return fmt.Errorf("error recuperando pólizas: %v", err)
	}

	for _, poliza := range polizas {
		ctx := MigrationContext{
			Ramo:    poliza["RAMO"],
			Npoliza: poliza["NPOLIZA"],
			Npolori: poliza["NPOLORI"],
		}

		log.Printf("Procesando póliza: RAMO %s, NPOLIZA %s, NPOLORI %s", ctx.Ramo, ctx.Npoliza, ctx.Npolori)

		// Insertar en CONTRACT_HEADER
		contractID, err := InsertContractHeaderForPoliza(tx, ctx)
		if err != nil {
			log.Printf("Error insertando CONTRACT_HEADER para póliza %s-%s: %v", ctx.Ramo, ctx.Npoliza, err)
			continue
		}
		ctx.ContractID = contractID

		// Insertar en POLICY
		if err := InsertIntoPolicy(tx, &ctx); err != nil {
			log.Printf("Error insertando POLICY para CONTRACT_ID %d: %v", ctx.ContractID, err)
			continue
		}

		// Insertar en POLICY_COVERAGE_VALUE
		if err := InsertPolicyCoverageValue(tx, &ctx); err != nil {
			log.Printf("Error insertando POLICY_COVERAGE_VALUE para CONTRACT_ID %d: %v", ctx.ContractID, err)
			continue
		}

		// Insertar en POLICY_PARAMETER
		if err := InsertPolicyParameter(tx, &ctx); err != nil {
			log.Printf("Error insertando en POLICY_PARAMETER: %v", err)
			continue
		}

		// Insertar en POLICY_ECONOMICS
		if err := InsertPolicyEconomics(tx, &ctx); err != nil {
			log.Printf("Error insertando en POLICY_ECONOMICS: %v", err)
			continue
		}

		// Insertar REQUEST
		requestID, err := InsertSingleRequestForContract(tx, &ctx)
		if err != nil {
			log.Printf("Error insertando REQUEST para CONTRACT_ID %d: %v", ctx.ContractID, err)
			continue
		}
		ctx.RequestID = requestID

		// Insertar REQUEST_COVERAGE_VALUE
		if err := InsertRequestCoverageValue(tx, &ctx); err != nil {
			log.Printf("Error insertando REQUEST_COVERAGE_VALUE para REQUEST_ID %d: %v", ctx.RequestID, err)
		}

		// Insertar REQUEST_ECONOMICS
		if err := InsertRequestEconomics(tx, &ctx); err != nil {
			log.Printf("Error insertando REQUEST_ECONOMICS para REQUEST_ID %d: %v", ctx.RequestID, err)
		}

		// Insertar REQUEST_PARAMETER
		if err := InsertRequestParameter(tx, &ctx); err != nil {
			log.Printf("Error insertando REQUEST_PARAMETER para REQUEST_ID %d: %v", ctx.RequestID, err)
		}
	}

	return nil
}

// GetPolizasFromDynamicTable obtiene las pólizas desde tiisa.polizas
func GetPolizasFromDynamicTable(tx *sql.Tx) ([]map[string]string, error) {
	query := "SELECT RAMO, NPOLIZA, NPOLORI FROM tiisa.polizas"
	rows, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query para obtener pólizas: %v", err)
	}
	defer rows.Close()

	var polizas []map[string]string
	for rows.Next() {
		var ramo, npoliza, npolori string
		if err := rows.Scan(&ramo, &npoliza, &npolori); err != nil {
			return nil, fmt.Errorf("error escaneando fila: %v", err)
		}
		poliza := map[string]string{
			"RAMO":    ramo,
			"NPOLIZA": npoliza,
			"NPOLORI": npolori,
		}
		polizas = append(polizas, poliza)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return polizas, nil
}
