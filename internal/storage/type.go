package storage

type ColumnType string

const (
	IntType   ColumnType = "INT"   // IDs, counts, ages, numeric identifiers
	TextType  ColumnType = "TEXT"  // Names, emails, descriptions, strings
	BoolType  ColumnType = "BOOL"  // Flags, active/inactive status, yes/no fields
	DateType  ColumnType = "DATE"  // created_at, updated_at, or date fields
	FloatType ColumnType = "FLOAT" // Prices, amounts, balances
)
