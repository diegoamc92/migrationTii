package database

type MigrationContext struct {
	Ramo       string // Código de ramo
	Npoliza    string // Número de póliza actual
	NpolOri    string // Número de póliza original (si aplica)
	Asegurado  string // Nombre del asegurado (si aplica)
	ContractID int64
	RequestID  int64
}
