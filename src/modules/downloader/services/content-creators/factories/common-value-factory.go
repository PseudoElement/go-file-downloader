package random_value_factory

import (
	"fmt"
	"strconv"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	ccr_models "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/models"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
)

func CommonRandomValueFactory(params ccr_models.RandomValueCreatorParams) (string, error) {
	var value string
	switch params.ValueType {
	case sql_constants.NUMBER:
		value = strconv.Itoa(int(custom_utils.CreateRandomNumber(params.Min, params.Max)))
	case sql_constants.CAR:
		value = custom_utils.CreateRandomCarName(params.Min, params.Max)
	case sql_constants.FIRST_NAME:
		value = custom_utils.CreateRandomFirstName(params.Min, params.Max)
	case sql_constants.LAST_NAME:
		value = custom_utils.CreateRandomLastName(params.Min, params.Max)
	case sql_constants.COUNTRY:
		value = custom_utils.CreateRandomCountryName(params.Min, params.Max)
	case sql_constants.WORK:
		value = custom_utils.CreateRandomWork(params.Min, params.Max)
	case sql_constants.STRING:
		value = custom_utils.CreateRandomWord(params.Min, params.Max, false)
	case sql_constants.AUTO_INCREMENT:
		if params.IncrementFn == nil {
			return "", fmt.Errorf("[TextContentCreator] params.incrementFn can't be empty!")
		}
		value = strconv.Itoa(params.IncrementFn())
	case sql_constants.DATE:
		value = custom_utils.CreateRandomDate(params.Min, params.Max)
	default:
		value = "unknown"
	}

	return value, nil
}
