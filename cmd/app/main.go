package main

import (
	"log"
	"migrationTii/config"
	"migrationTii/internal/data_loader"
	"migrationTii/internal/database"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Crear conexión a la base de datos
	db, err := database.CreateConnection(cfg)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	// Cargar datos desde TXT al esquema tiisa
	txtFiles := map[string]string{
		"asegurados": "pkg/utils/data/ASEGURADOS.csv",
		"polizas":    "pkg/utils/data/POLIZAS.csv",
		"coberturas": "pkg/utils/data/COBERTURAS.csv",
	}

	for tableName, filePath := range txtFiles {
		if err := data_loader.ProcessTXTFile(db, filePath, tableName); err != nil {
			log.Fatalf("Error procesando archivo %s: %v", filePath, err)
		}
		log.Printf("Datos de %s procesados correctamente.", tableName)
	}

	// Procesar datos del esquema tiisa al esquema ISB
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error iniciando transacción: %v", err)
	}

	polizas, err := database.GetPolizasFromDynamicTable(tx)
	if err != nil {
		log.Fatalf("Error obteniendo pólizas y relaciones: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Fatalf("Transacción revertida debido a un error: %v", err)
		} else {
			tx.Commit()
			log.Println("Transacción completada exitosamente.")
		}
	}()

	for _, poliza := range polizas {
		migrationContext := database.MigrationContext{
			Npoliza: poliza["Npoliza"],
			// Add other fields as necessary
		}
		if err := database.ProcessPoliza(tx, migrationContext); err != nil {
			log.Printf("Error procesando póliza %s: %v", poliza["Npoliza"], err)
			continue
		}
		log.Printf("Póliza %s procesada correctamente.", poliza["Npoliza"])
	}

	log.Println("Migración completada exitosamente.")
}
