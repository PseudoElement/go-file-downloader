package sql_constants

import value_types "github.com/pseudoelement/go-file-downloader/src/constants/value-types"

var VALUE_TYPE_TO_SQL_TYPE = map[string]string{
	value_types.AUTO_INCREMENT: "SERIAL",
	value_types.BOOL:           "BIT",
	value_types.NUMBER:         "INT",
	value_types.STRING:         "TEXT",
}
