package data_loader

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"migrationTii/internal/utils"
	"os"
	"strings"
)

// Process TXTFile procesa archivos TXT y crea tablas dinámicamente
func ProcessTXTFile(db *sql.DB, filePath string, tableName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error abriendo archivo TXT: %v", err)
	}
	defer file.Close()

	// Leer encabezados y contenido
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	var headers []string
	var rows [][]string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Split(line, ";")
		if lineNumber == 0 {
			headers = fields
		} else {
			rows = append(rows, fields)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error leyendo archivo TXT: %v", err)
	}

	// Crear tabla dinámica en tiisa
	createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS tiisa.%s (", tableName)
	for i, header := range headers {
		createTableQuery += fmt.Sprintf("%s VARCHAR(255)", header)
		if i < len(headers)-1 {
			createTableQuery += ", "
		}
	}
	createTableQuery += ") CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

	if _, err := db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("error creando tabla %s: %v", tableName, err)
	}

	// Insertar datos en la tabla creada
	insertQuery := fmt.Sprintf("INSERT INTO tiisa.%s (%s) VALUES (%s)",
		tableName,
		strings.Join(headers, ", "),
		strings.Repeat("?, ", len(headers)-1)+"?")

	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("error preparando query de inserción: %v", err)
	}
	defer stmt.Close()

	for _, row := range rows {
		values := make([]interface{}, len(row))
		for i, val := range row {
			cleanedValue := strings.TrimSpace(val)
			if headers[i] == "RUT" {
				cleanedValue = utils.CleanRUT(cleanedValue)
			}
			values[i] = cleanedValue
		}
		if _, err := stmt.Exec(values...); err != nil {
			log.Printf("Error insertando fila: %v", err)
		}
	}

	log.Printf("Archivo %s procesado y datos insertados en tabla %s", filePath, tableName)
	return nil
}
