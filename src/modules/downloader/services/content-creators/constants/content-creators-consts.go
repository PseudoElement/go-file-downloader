package ccr_constants

import sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"

var TEXTLIKE_VALUE_TYPES = []string{
	sql_constants.STRING,
	sql_constants.CAR,
	sql_constants.COUNTRY,
	sql_constants.FIRST_NAME,
	sql_constants.LAST_NAME,
	sql_constants.WORK,
	sql_constants.DATE,
	sql_constants.BOOL,
}

var VALUE_LENGTHS = map[string]int{
	sql_constants.DATE:   12,
	sql_constants.NUMBER: 8,
}
