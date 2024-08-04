package sql_utils

func WrapStringInSingleQuotes(str string) string {
	return "'" + str + "'"
}
