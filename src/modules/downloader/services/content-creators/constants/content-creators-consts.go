package content_creator_constants

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
