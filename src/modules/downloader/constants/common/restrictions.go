package common_constants

import sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"

type ColumnRestriction struct {
	MaximumMin int
	MinimalMax int
}

var RESTRICTIONS_BY_COLUMN_TYPE = map[string]ColumnRestriction{
	sql_constants.FIRST_NAME: ColumnRestriction{
		MaximumMin: 10,
		MinimalMax: 11,
	},
	sql_constants.LAST_NAME: ColumnRestriction{
		MaximumMin: 7,
	},
	sql_constants.CAR: ColumnRestriction{
		MaximumMin: 7,
	},
	sql_constants.COUNTRY: ColumnRestriction{
		MaximumMin: 10,
		MinimalMax: 11,
	},
}
