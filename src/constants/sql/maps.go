package sql_constants

var VALUE_TYPE_TO_SQL_TYPE = map[string]string{
	AUTO_INCREMENT: "SERIAL",
	BOOL:           "BIT",
	NUMBER:         "INT",
	STRING:         "TEXT",
}
