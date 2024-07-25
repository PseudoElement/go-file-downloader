package types_module

type CommonReqBody struct {
	RowsCount int    `json:"rows_count"`
	DocType   string `json:"doc_type"`
	DocName   string `json:"doc_name"`
}

type DownloadTextReqBody struct {
	CommonReqBody
	ColumnsData []TextColumnInfo `json:"columns_data"`
}

type DownloadSqlReqBody struct {
	CommonReqBody
	ColumnsData     []SqlColumnInfo `json:"columns_data"`
	NeedCreateTable bool            `json:"need_create_table"`
}

type TextColumnInfo struct {
	Name string `json:"name"`
	/* 'BOOL' or 'NUMBER' or 'STRING' */
	Type string `json:"type"`
	// from 0 to 100%
	NullValuesPercent int `json:"null_values_percent"`
	Min               int `json:"min"`
	Max               int `json:"max"`
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
