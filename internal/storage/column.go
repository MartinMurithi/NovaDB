package storage

type Column struct {
	Column   string
	ColumnType   ColumnType
	IsPrimaryKey bool
	IsUnique     bool
}

