package main

import (
	"database/sql"
	"fmt"
	"log"
	"migrationTii/config"
	"migrationTii/internal/data_loader"
	"migrationTii/internal/database"
	"migrationTii/pkg/report"
)

//	func main() {
//		// Crear reporte
//		r := report.NewReport("pkg/report/reporte/execution_report.log")
//
//		// 1. Cargar configuración
//		cfg, err := config.LoadConfig()
//		if err != nil {
//			r.MarkFailure(err)
//			r.Save()
//			log.Fatalf("Error cargando configuración: %v", err)
//		}
//		r.Add("Configuración cargada correctamente.")
//
//		// 2. Crear conexión a la base de datos
//		db, err := database.CreateConnection(cfg)
//		if err != nil {
//			r.MarkFailure(err)
//			r.Save()
//			log.Fatalf("Error conectando a la DB: %v", err)
//		}
//		defer db.Close()
//		r.Add("Conexión a la base de datos establecida.")
//
//		// Configurar nivel de aislamiento
//		_, err = db.Exec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
//		if err != nil {
//			log.Fatalf("Error configurando nivel de aislamiento: %v", err)
//		}
//
//		// 3. Iniciar transacción
//		start := time.Now()
//		tx, err := db.Begin()
//		if err != nil {
//			r.MarkFailure(err)
//			r.Save()
//			log.Fatalf("Error iniciando transacción: %v", err)
//		}
//		r.Add("Transacción iniciada.")
//
//		defer func() {
//			if err != nil {
//				r.MarkFailure(err)
//				if rollbackErr := tx.Rollback(); rollbackErr != nil {
//					log.Fatalf("Error al hacer rollback: %v", rollbackErr)
//				}
//				r.Add("Rollback realizado exitosamente.")
//				r.Save()
//				log.Fatalf("Error detectado: %v", err)
//			}
//		}()
//
//		// 4. Crear tablas temporales
//		blockStart := time.Now()
//
//		// Verificar existencia de archivos
//		if !fileExists("pkg/utils/data/MIGSA_ASEGURADOS_CSV.csv") {
//			log.Fatalf("Archivo no encontrado: MIGSA_ASEGURADOS_CSV.csv")
//		} else {
//			log.Println("Archivo encontrado: MIGSA_ASEGURADOS_CSV.csv")
//		}
//
//		if !fileExists("pkg/utils/data/MIGSA_POLIZAS-18-12_CSV.csv") {
//			log.Fatalf("Archivo no encontrado: MIGSA_POLIZAS-18-12_CSV.csv")
//		} else {
//			log.Println("Archivo encontrado: MIGSA_POLIZAS-18-12_CSV.csv")
//		}
//
//		if !fileExists("pkg/utils/data/MIGSA_COBERTURAS_CSV.csv") {
//			log.Fatalf("Archivo no encontrado: MIGSA_COBERTURAS_CSV.csv")
//		} else {
//			log.Println("Archivo encontrado: MIGSA_COBERTURAS_CSV.csv")
//		}
//
//		if err := database.CreateTempTable(tx); err != nil {
//			return
//		}
//
//		blockStart = time.Now()
//		if err := database.CreateCleanedTempTable(tx); err != nil {
//			return
//		}
//
//		aseguradosData, err := data_loader.CleanAseguradosData("pkg/utils/data/ASEGURADOS.csv")
//		if err != nil {
//			log.Fatalf("Error procesando CSV de asegurados: %v", err)
//		}
//
//		if err := database.LoadAseguradosData(tx, aseguradosData); err != nil {
//			log.Fatalf("Error insertando datos de asegurados en temp_csv_asegurados: %v", err)
//		}
//
//		r.Add(fmt.Sprintf("Datos de asegurados cargados en %v.", time.Since(blockStart)))
//
//		blockStart = time.Now()
//		polizasData, err := data_loader.CleanDataPolizas("pkg/utils/data/POLIZAS.csv")
//		if err != nil {
//			return
//		}
//		if err := database.LoadPolizasData(tx, polizasData); err != nil {
//			return
//		}
//		r.Add(fmt.Sprintf("Datos de pólizas cargados en %v.", time.Since(blockStart)))
//
//		coberturasData, err := data_loader.CleanAndProcessCoverageData("pkg/utils/data/COBERTURAS.csv")
//		if err != nil {
//			log.Fatalf("Error procesando CSV de coberturas: %v", err)
//		}
//
//		if err := database.LoadCoberturasData(tx, coberturasData); err != nil {
//			log.Fatalf("Error insertando datos de asegurados en temp_csv_coberturas: %v", err)
//		}
//
//		r.Add(fmt.Sprintf("Datos de coberturas cargados en %v.", time.Since(blockStart)))
//type MigrationContext struct {
//	Ramo      string // Código de ramo
//	Npoliza   string // Número de póliza actual
//	NpolOri   string // Número de póliza original (si aplica)
//	Asegurado string // Nombre del asegurado (si aplica)
//}

func main() {
	// Crear reporte
	r := report.NewReport("pkg/report/reporte/execution_report.log")

	// Configuración y conexión a la base de datos
	cfg, err := config.LoadConfig()
	if err != nil {
		r.MarkFailure(err)
		r.Save()
		log.Fatalf("Error cargando configuración: %v", err)
	}
	db, err := database.CreateConnection(cfg)
	if err != nil {
		r.MarkFailure(err)
		r.Save()
		log.Fatalf("Error conectando a la DB: %v", err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		r.MarkFailure(err)
		r.Save()
		log.Fatalf("Error iniciando transacción: %v", err)
	}

	if err := database.CreateTempTable(tx); err != nil {
		return
	}

	if err := database.CreateCleanedTempTable(tx); err != nil {
		return
	}

	// Cargar datos de temporales
	aseguradosData, err := data_loader.CleanAseguradosData("pkg/utils/data/ASEGURADOS.csv")
	if err != nil {
		log.Fatalf("Error procesando CSV de asegurados: %v", err)
	}
	if err := database.LoadAseguradosData(tx, aseguradosData); err != nil {
		log.Fatalf("Error insertando datos de asegurados en temp_csv_asegurados: %v", err)
	}

	polizasData, err := data_loader.CleanDataPolizas("pkg/utils/data/POLIZAS_CSV.csv")
	if err != nil {
		log.Fatalf("Error procesando CSV de pólizas: %v", err)
	}
	if err := database.LoadPolizasData(tx, polizasData); err != nil {
		log.Fatalf("Error insertando datos de pólizas en temp_csv_polizas: %v", err)
	}

	coberturasData, err := data_loader.CleanAndProcessCoverageData("pkg/utils/data/COBERTURAS.csv")
	if err != nil {
		log.Fatalf("Error procesando CSV de coberturas: %v", err)
	}
	if err := database.LoadCoberturasData(tx, coberturasData); err != nil {
		log.Fatalf("Error insertando datos de coberturas en temp_csv_coberturas: %v", err)
	}

	// Crear y limpiar RUT
	if err := database.CreateTempCleanedRUT(tx); err != nil {
		log.Fatalf("Error creando tabla temporal RUT: %v", err)
	}

	// Obtener pólizas y procesarlas
	polizas, err := database.GetPolizasFromTempTable(tx)
	if err != nil {
		log.Fatalf("Error recuperando pólizas: %v", err)
	}

	for _, poliza := range polizas {
		context := database.MigrationContext{
			Ramo:    poliza["RAMO"],
			Npoliza: poliza["NPOLIZA"],
			NpolOri: poliza["NPOLORI"],
		}

		// Procesar cada póliza
		if err := processPoliza(tx, context); err != nil {
			log.Printf("Error procesando póliza %s: %v", context.Npoliza, err)
			continue
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		log.Fatalf("Error confirmando la transacción: %v", err)
	}

}

func processPoliza(tx *sql.Tx, context database.MigrationContext) error {
	// Insertar PARTY
	if err := database.InsertPartyWithContext(tx, context); err != nil {
		return fmt.Errorf("error en PARTY para póliza %s: %v", context.Npoliza, err)
	}

	// Insertar IDENTIFICATION
	if err := database.InsertIdentificationWithContext(tx); err != nil {
		return fmt.Errorf("error en IDENTIFICATION para póliza %s: %v", context.Npoliza, err)
	}

	// Asociar IDENTIFICATION con PARTY
	if err := database.AssociatePartyIdentificationWithContext(tx, context); err != nil {
		log.Fatalf("Error asociando PARTY_IDENTIFICATION: %v", err)
	}

	if err := database.InsertEmail(tx, context); err != nil {

		log.Fatalf("Error insertando EMAIL: %v", err)
	}

	if err := database.InsertPhone(tx, context); err != nil {
		log.Fatalf("Error insertando PHONE: %v", err)
	}

	if err := database.InsertAddress(tx, context); err != nil {
		log.Fatalf("Error insertando ADDRESS: %v", err)
	}

	if err := database.AssociatePartyAddress(tx, context); err != nil {
		log.Fatalf("Error insertando PARTY_ADDRESS: %v", err)
	}

	if err := database.InsertPersonData(tx, context); err != nil {
		log.Fatalf("Error insertando PERSON: %v", err)
	}

	if err := database.InsertPaymentTerm(tx, context); err != nil {
		log.Fatalf("Error insertando PAYMENT_TERM: %v", err)
	}

	if err := database.CreateTempIssuanceDates(tx, context); err != nil {
		log.Fatalf("Error creando tabla temporal: %v", err)
	}

	log.Println("Recuperando lista de pólizas desde temp_csv_polizas...")
	polizas, err := database.GetPolizasFromTempTable(tx)
	if err != nil {
		log.Fatalf("Error recuperando pólizas: %v", err)
	}

	for _, poliza := range polizas {
		// Crear contexto para la póliza actual
		ctx := database.MigrationContext{
			Ramo:    poliza["RAMO"],
			Npoliza: poliza["NPOLIZA"],
			NpolOri: poliza["NPOLORI"],
		}

		log.Printf("Procesando póliza: RAMO %s, NPOLIZA %s, NPOLORI %s", ctx.Ramo, ctx.Npoliza, ctx.NpolOri)

		// Insertar titular en PERSON_INSURED
		if err := database.InsertPersonInsured(tx, &context); err != nil { // `true` indica titular
        	log.Fatalf("Error insertando titular en PERSON_INSURED: %v", err)
    	}

    	// Procesar adherentes
    	if err := database.InsertAdherents(tx, &context); err != nil {
        	log.Fatalf("Error insertando adherentes: %v", err)
    	}

		// Insertar en CONTRACT_HEADER
		contractID, err := database.InsertContractHeaderForPoliza(tx, ctx)
		if err != nil {
			log.Printf("Error insertando CONTRACT_HEADER para póliza %s-%s: %v", ctx.Ramo, ctx.Npoliza, err)
			continue
		}
		ctx.ContractID = contractID // Actualizar el contexto con el ContractID

		
		// Crear PARTY_RELATION para relaciones titular-adherentes
		if err := database.InsertPartyRelation(tx, &context); err != nil {
    		log.Fatalf("Error insertando PARTY_RELATION: %v", err)
		}

		// Insertar titular en PERSON_INSURED
		if err := database.InsertPersonInsured(tx, &context); err != nil { // `true` indica titular
        	log.Fatalf("Error insertando titular en PERSON_INSURED: %v", err)
    	}

    	// Procesar adherentes
    	if err := database.InsertAdherents(tx, &context); err != nil {
        	log.Fatalf("Error insertando adherentes: %v", err)
    	}

		// Insertar en POLICY
		if err := database.InsertIntoPolicy(tx, &ctx); err != nil {
			log.Printf("Error insertando POLICY para CONTRACT_ID %d: %v", ctx.ContractID, err)
			continue
		}

		log.Printf("Insertando en POLICY_COVERAGE_VALUE...")
		if err := database.InsertPolicyCoverageValue(tx, &ctx); err != nil {
			log.Printf("Error insertando POLICY_COVERAGE_VALUE para CONTRACT_ID %d: %v", ctx.ContractID, err)
			continue
		}
		log.Printf("Datos insertados en POLICY_COVERAGE_VALUE correctamente.")

		log.Printf("Insertando en POLICY_PARAMETER...")
		if err := database.InsertPolicyParameter(tx, &ctx); err != nil {
			log.Printf("Error insertando en POLICY_PARAMETER: %v", err)
			continue
		}
		log.Printf("Datos insertados en POLICY_PARAMETER correctamente.")

		log.Printf("Insertando en POLICY_ECONOMICS...")
		if err := database.InsertPolicyEconomics(tx, &ctx); err != nil {
			log.Printf("Error insertando en POLICY_ECONOMICS: %v", err)
			continue
		}
		log.Printf("Datos insertados en POLICY_ECONOMICS correctamente.")

		// Insertar en REQUEST
		requestID, err := database.InsertSingleRequestForContract(tx, &ctx)
		if err != nil {
			log.Printf("Error insertando REQUEST para CONTRACT_ID %d: %v", ctx.ContractID, err)
			continue
		}
		ctx.RequestID = requestID // Actualizar el contexto con el RequestID

		// Insertar en REQUEST_COVERAGE_VALUE
		if err := database.InsertRequestCoverageValue(tx, &ctx); err != nil {
			log.Printf("Error insertando REQUEST_COVERAGE_VALUE para REQUEST_ID %d: %v", ctx.RequestID, err)
		}

		// Insertar en REQUEST_ECONOMICS
		if err := database.InsertRequestEconomics(tx, &ctx); err != nil {
			log.Printf("Error insertando REQUEST_ECONOMICS para REQUEST_ID %d: %v", ctx.RequestID, err)
		}

		// Insertar en REQUEST_PARAMETER
		if err := database.InsertRequestParameter(tx, &ctx); err != nil {
			log.Printf("Error insertando REQUEST_PARAMETER para REQUEST_ID %d: %v", ctx.RequestID, err)
		}

		//if err := database.InsertBillingStatement(tx, &ctx); err != nil {
		//	log.Fatalf("Error insertando en BILLING_STATEMENT: %v", err)
		//}
	}

	return nil
}

//
//	// 6. Procesar datos de asegurados
//	log.Println("Procesando datos de asegurados...")
//	blockStart = time.Now()
//	if err := database.InsertPartyData(tx); err != nil {
//		//handleError(database.InsertPartyData(tx), "Error insertando PARTY", r)
//		log.Fatalf("Error insertando PARTY: %v", err)
//	}
//	if err := database.CreateTempCleanedRUT(tx); err != nil {
//		return
//	}
//
//	if err := database.InsertInsuredData(tx); err != nil {
//		log.Fatalf("Error insertando INSURED: %v", err)
//	}
//
//	if err := database.InsertPolicyHolderData(tx); err != nil {
//		log.Fatalf("Error insertando POLICY_HOLDER: %v", err)
//	}
//
//	if err := database.InsertIdentification(tx); err != nil {
//		//handleError(database.InsertIdentification(tx), "Error insertando IDENTIFICATION", r)
//		log.Fatalf("Error insertando IDENTIFICATION: %v", err)
//	}
//	if err := database.AssociatePartyIdentification(tx); err != nil {
//		log.Fatalf("Error insertando PARTY_IDENTIFICATION: %v", err)
//	}
//	log.Println("Iniciando inserción en EMAIL...")
//	if err := database.InsertEmail(tx); err != nil {
//		//handleError(database.InsertEmail(tx), "Error insertando EMAIL", r)
//		log.Fatalf("Error insertando EMAIL: %v", err)
//	}
//	log.Println("EMAIL insertado correctamente.")
//
//	log.Println("Iniciando inserción en PHONE...")
//	if err := database.InsertPhone(tx); err != nil {
//		log.Fatalf("Error insertando PHONE: %v", err)
//	}
//	log.Println("PHONE insertado correctamente.")
//
//	log.Println("Insertando en ADDRESS...")
//	if err := database.InsertAddress(tx); err != nil {
//		log.Fatalf("Error insertando ADDRESS: %v", err)
//	}
//	log.Println("Datos insertados en ADDRESS correctamente.")
//
//	log.Println("Insertando en PARTY_ADDRESS...")
//	if err := database.AssociatePartyAddress(tx); err != nil {
//		log.Fatalf("Error insertando PARTY_ADDRESS: %v", err)
//	}
//	log.Println("Datos insertados en PARTY_ADDRESS correctamente.")
//
//	log.Println("Insertando en PERSON...")
//	if err := database.InsertPersonData(tx); err != nil {
//		log.Fatalf("Error insertando PERSON: %v", err)
//	}
//	log.Println("Datos insertados en PERSON correctamente.")
//
//	log.Println("Insertando en PAYMENT_TERM...")
//	if err := database.InsertPaymentTerm(tx); err != nil {
//		log.Fatalf("Error insertando PAYMENT_TERM: %v", err)
//	}
//	log.Println("Datos insertados en PAYMENT_TERM correctamente.")
//
//	if err := database.CreateTempIssuanceDates(tx); err != nil {
//		log.Fatalf("Error creando tabla temporal: %v", err)
//	}
// Crear lista de pólizas (RAMO, NPOLIZA)
//log.Println("Recuperando lista de pólizas desde temp_csv_polizas...")
//polizas, err := database.GetPolizasFromTempTable(tx)
//if err != nil {
//	log.Fatalf("Error recuperando pólizas: %v", err)
//}
//
//for _, poliza := range polizas {
//	ramo := poliza["RAMO"]
//	nPoliza := poliza["NPOLIZA"]
//	nPolOri := poliza["NPOLORI"]
//
//	// Insertar en CONTRACT_HEADER
//	contractID, err := database.InsertContractHeaderForPoliza(tx, ramo, nPoliza)
//	if err != nil {
//		log.Printf("Error insertando CONTRACT_HEADER para póliza %s-%s: %v", ramo, nPoliza, err)
//		continue
//	}
//
//	// Insertar en POLICY
//	if err := database.InsertIntoPolicy(tx, contractID, ramo, nPolOri); err != nil {
//		log.Printf("Error insertando POLICY para CONTRACT_ID %d: %v", contractID, err)
//		continue
//	}
//
//	log.Printf("Insertando en POLICY_COVERAGE_VALUE...")
//	if err := database.InsertPolicyCoverageValue(tx, ramo, nPolOri); err != nil {
//		log.Fatalf("Error insertando POLICY_COVERAGE_VALUE para CONTRACT_ID %d: %v", contractID, err)
//	}
//	log.Printf("Datos insertados en POLICY_COVERAGE_VALUE correctamente.")
//
//	log.Printf("Insertando en POLICY_PARAMETER...")
//	if err := database.InsertPolicyParameter(tx, ramo, nPolOri); err != nil {
//		log.Fatalf("Error insertando en POLICY_PARAMETER: %v", err)
//	}
//	log.Printf("Datos insertados en POLICY_PARAMETER correctamente.")
//
//	log.Printf("Insertando en POLICY_ECONOMICS...")
//	if err := database.InsertPolicyEconomics(tx, ramo, nPolOri); err != nil {
//		log.Fatalf("Error insertando en POLICY_ECONOMICS: %v", err)
//	}
//	log.Printf("Datos insertados en POLICY_ECONOMICS correctamente.")
//
//	//log.Printf("Insertando en BILLING_STATEMENT...")
//	//if err := database.InsertBillingStatement(tx); err != nil {
//	//	log.Fatalf("Error insertando en BILLING_STATEMENT: %v", err)
//	//}
//	//log.Printf("Datos insertados en BILLING_STATEMENT correctamente.")
//
//	// Insertar en REQUEST
//	requestID, err := database.InsertSingleRequestForContract(tx, contractID)
//	if err != nil {
//		log.Printf("Error insertando REQUEST para CONTRACT_ID %d: %v", contractID, err)
//		continue
//	}
//
//	// Insertar en REQUEST_COVERAGE_VALUE
//	if err := database.InsertRequestCoverageValue(tx, requestID); err != nil {
//		log.Printf("Error insertando REQUEST_COVERAGE_VALUE para REQUEST_ID %d: %v", requestID, err)
//	}
//
//	// Insertar en REQUEST_ECONOMICS
//	if err := database.InsertRequestEconomics(tx, requestID); err != nil {
//		log.Printf("Error insertando REQUEST_ECONOMICS para REQUEST_ID %d: %v", requestID, err)
//	}
//
//	// Insertar en REQUEST_PARAMETER
//	if err := database.InsertRequestParameter(tx, requestID); err != nil {
//		log.Printf("Error insertando REQUEST_PARAMETER para REQUEST_ID %d: %v", requestID, err)
//	}
//}
//
//	log.Println("Proceso completado correctamente.")
//
//	//log.Printf("Insertando en POLICY_COVERAGE_VALUE...")
//	//if err := database.InsertPolicyCoverageValue(tx); err != nil {
//	//	log.Fatalf("Error insertando en POLICY_COVERAGE_VALUE: %v", err)
//	//}
//	//log.Printf("Datos insertados en POLICY_COVERAGE_VALUE correctamente.")
//	//
//	//log.Printf("Insertando en POLICY_PARAMETER...")
//	//if err := database.InsertPolicyParameter(tx); err != nil {
//	//	log.Fatalf("Error insertando en POLICY_PARAMETER: %v", err)
//	//}
//	//log.Printf("Datos insertados en POLICY_PARAMETER correctamente.")
//	//
//	//log.Printf("Insertando en POLICY_ECONOMICS...")
//	//if err := database.InsertPolicyEconomics(tx); err != nil {
//	//	log.Fatalf("Error insertando en POLICY_ECONOMICS: %v", err)
//	//}
//	//log.Printf("Datos insertados en POLICY_ECONOMICS correctamente.")
//
//	//r.Add(fmt.Sprintf("Datos de POLICY procesados en %v.", time.Since(blockStart)))
//
//	// 10. Insertar en BILLING_STATEMENT
//	//log.Println("Procesando datos de BILLING_STATEMENT...")
//	//blockStart = time.Now()
//	//
//	//log.Printf("Insertando en BILLING_STATEMENT...")
//	//if err := database.InsertBillingStatement(tx); err != nil {
//	//	log.Fatalf("Error insertando en BILLING_STATEMENT: %v", err)
//	//}
//	//log.Printf("Datos insertados en BILLING_STATEMENT correctamente.")
//	//
//	//r.Add(fmt.Sprintf("Datos de BILLING_STATEMENT procesados en %v.", time.Since(blockStart)))
//
//	// 11. Confirmar transacción
//	log.Println("Intentando confirmar la transacción...")
//	if err := tx.Commit(); err != nil {
//		r.MarkFailure(err)
//		r.Save()
//		log.Fatalf("Error al confirmar la transacción: %v", err)
//	}
//	log.Println("Transacción confirmada correctamente.")
//
//	r.MarkSuccess()
//	r.Add(fmt.Sprintf("Transacción completada exitosamente en %v.", time.Since(start)))
//
//	// Guardar reporte final
//	if err := r.Save(); err != nil {
//		log.Fatalf("Error guardando el reporte: %v", err)
//	}
//
//	log.Println("Proceso completado correctamente.")
//}
//

//
//func handleError(err error, message string, r *report.Report) {
//	if err != nil {
//		r.MarkFailure(err)
//		log.Fatalf("%s: %v", message, err)
//	}
//}
