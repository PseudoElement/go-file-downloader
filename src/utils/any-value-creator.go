package custom_utils

import (
	"fmt"
	"strconv"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/constants/sql"
)

type RandomValueCreatorParams struct {
	/* from sql_constants BOOL/STRING/NUMBER/AUTO_INCREMENT */
	ValueType   string
	Min         int
	Max         int
	IncrementFn func() int
	IsSqlValue  bool
}

/*
* @param valueType from sql_constants BOOL/STRING/NUMBER/AUTO_INCREMENT
 */
func CreateRandomValueConvertedToString(params RandomValueCreatorParams) (string, error) {
	var value string
	switch params.ValueType {
	case sql_constants.BOOL:
		if params.IsSqlValue {
			value = string(CreateRandomByteForSql())
		} else {
			value = strconv.FormatBool(CreateRandomBool())
		}
	case sql_constants.NUMBER:
		value = strconv.Itoa(CreateRandomNumber(params.Min, params.Max))
	case sql_constants.STRING:
		if params.IsSqlValue {
			value = CreateRandowWordForSqlTable(params.Min, params.Max, false)
		} else {
			value = CreateRandomWord(params.Min, params.Max, false)
		}
	case sql_constants.AUTO_INCREMENT:
		if params.IncrementFn == nil {
			return "", fmt.Errorf("[CreateRandomValueConvertedToString] params.incrementFn can't be empty!")
		}
		value = strconv.Itoa(params.IncrementFn())
	default:
		value = "unknown"
	}

	return value, nil
}
