package storage

type Column struct {
	Name   string
	ColumnType   ColumnType
	IsPrimaryKey bool
	IsUnique     bool
}

