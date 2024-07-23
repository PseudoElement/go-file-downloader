package downloader_module

type DownloadTextReqBody struct {
	ColumnsData []TextColumnInfo `json:"columns_data"`
	RowsCount   int              `json:"rows_count"`
	DocType     string           `json:"doc_type"`
	DocName     string           `json:"doc_name"`
}

type DownloadSqlReqBody struct {
	ColumnsData     []SqlColumnInfo `json:"columns_data"`
	RowsCount       int             `json:"rows_count"`
	DocType         string          `json:"doc_type"`
	NeedCreateTable bool            `json:"need_create_table"`
}

type TextColumnInfo struct {
	Name string `json:"name"`
	/* 'bool' or 'number' or 'string' */
	Type              string `json:"type"`
	NullValuesPercent int    `json:"null_values_percent"`
	Min               int    `json:"min"`
	Max               int    `json:"max"`
}

type SqlColumnInfo struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	NullValuesPercent int    `json:"null_values_percent"`
	IsPrimaryKey      bool   `json:"is_primary_key"`
	ForeignKeyData    struct {
		RefTableName  string `json:"reference_table_name"`
		RefColumnName string `json:"reference_column_name"`
	} `json:"foreign_key_data"`
}
