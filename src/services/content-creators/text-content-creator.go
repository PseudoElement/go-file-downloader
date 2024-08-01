package content_creators

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/constants/sql"
	services_models "github.com/pseudoelement/go-file-downloader/src/services/models"
	types_module "github.com/pseudoelement/go-file-downloader/src/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
)

type TextContentCreator struct {
	mu sync.Mutex
}

func NewTextContentCreator() *TextContentCreator {
	return &TextContentCreator{}
}

func (srv *TextContentCreator) CreateFileContent(body interface{}) (string, error) {
	textBody, ok := body.(types_module.DownloadTextReqBody)
	if !ok {
		return "", fmt.Errorf("[TextContentCreator] invalid body type")
	}

	rowsCountWithHeader := textBody.RowsCount + 1
	columnsWithValuesInCells, err := srv.createCellsForColumns(textBody)
	if err != nil {
		return "", err
	}
	fileContent := srv.concatAllCellsInTable(columnsWithValuesInCells, rowsCountWithHeader, len(textBody.ColumnsData))

	return fileContent, nil
}

func (srv *TextContentCreator) CreateFileContentAsync(body interface{}) (string, error) {
	textBody, ok := body.(types_module.DownloadTextReqBody)
	if !ok {
		return "", fmt.Errorf("[TextContentCreator] invalid body type")
	}

	rowsCountWithHeader := textBody.RowsCount + 1
	columnsWithValuesInCells, err := srv.createCellsForColumnsAsync(textBody)
	if err != nil {
		return "", err
	}
	fileContent := srv.concatAllCellsInTable(columnsWithValuesInCells, rowsCountWithHeader, len(textBody.ColumnsData))

	return fileContent, nil
}

func (srv *TextContentCreator) concatAllCellsInTable(columnsWithFullCells [][]string, rowsCountWithHeader int, columnsCount int) string {
	var fileContent string
	for rowNumber := range rowsCountWithHeader {
		var row string
		for columnNumber := 0; columnNumber < columnsCount; columnNumber++ {
			row += columnsWithFullCells[columnNumber][rowNumber]
		}
		fileContent += row
	}

	return fileContent
}

func (srv *TextContentCreator) createCellsForColumns(body types_module.DownloadTextReqBody) ([][]string, error) {
	columnsWithValuesInCells := [][]string{}

	for i, columnData := range body.ColumnsData {
		isLastColumn := i == len(body.ColumnsData)-1
		cellsOfColumn := []string{}
		incrementFn := custom_utils.AutoIncrement(1)

		for rowNumber := range body.RowsCount {
			if rowNumber == 0 {
				header := srv.addSpaceOrParagraphToValue(columnData.Name, isLastColumn, columnData)
				cellsOfColumn = append(cellsOfColumn, header)
			}

			value, err := srv.createRandomValue(services_models.RandomValueCreatorParams{
				ValueType:   columnData.Type,
				Min:         columnData.Min,
				Max:         columnData.Max,
				IncrementFn: incrementFn,
			})
			if err != nil {
				return nil, err
			}

			if isNullValue := columnData.NullValuesPercent > rand.Intn(100); isNullValue {
				value = "null"
			}

			value = srv.addSpaceOrParagraphToValue(value, isLastColumn, columnData)
			cellsOfColumn = append(cellsOfColumn, value)
		}

		columnsWithValuesInCells = append(columnsWithValuesInCells, cellsOfColumn)
	}

	return columnsWithValuesInCells, nil
}

func (srv *TextContentCreator) createCellsForColumnsAsync(body types_module.DownloadTextReqBody) ([][]string, error) {
	columnsWithValuesInCells := make([][]string, len(body.ColumnsData))
	errorsChan := make(chan error, len(body.ColumnsData)*body.RowsCount)
	var wg sync.WaitGroup

	for i, column := range body.ColumnsData {
		wg.Add(1)
		go func(ind int, columnData types_module.TextColumnInfo) {
			defer wg.Done()
			isLastColumn := ind == len(body.ColumnsData)-1
			cellsOfColumn := make([]string, body.RowsCount+1)
			incrementFn := custom_utils.AutoIncrement(1)

			for rowNumber := range body.RowsCount {
				if rowNumber == 0 {
					header := srv.addSpaceOrParagraphToValue(columnData.Name, isLastColumn, columnData)
					// NECESSARY ASSIGN VALUE BY INDEX USING GO CONCURENCY
					cellsOfColumn[rowNumber] = header
				}

				value, err := srv.createRandomValue(services_models.RandomValueCreatorParams{
					ValueType:   columnData.Type,
					Min:         columnData.Min,
					Max:         columnData.Max,
					IncrementFn: incrementFn,
				})
				if err != nil {
					errorsChan <- err
					continue
				}
				if isNullValue := columnData.NullValuesPercent > rand.Intn(100); isNullValue {
					value = "null"
				}

				value = srv.addSpaceOrParagraphToValue(value, isLastColumn, columnData)
				cellsOfColumn[rowNumber+1] = value
			}

			columnsWithValuesInCells[ind] = cellsOfColumn
		}(i, column)

	}
	wg.Wait()
	close(errorsChan)

	for err := range errorsChan {
		if err != nil {
			return nil, err
		}
	}

	return columnsWithValuesInCells, nil
}

func (srv *TextContentCreator) addSpaceOrParagraphToValue(value string, isLastColumn bool, columnData types_module.TextColumnInfo) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
			fmt.Println("Recovered value ===> ", value)
		}
	}()

	isHeader := value == columnData.Name

	if isHeader && isLastColumn {
		value += "\n\n"
	} else if isLastColumn {
		value += "\n"
	} else {
		max := math.Max(float64(len(columnData.Name)), float64(columnData.Max))
		columnWidth := math.Max(max, 5)
		spaces := strings.Repeat(" ", int(columnWidth)-len(value)+2)
		value += spaces
	}
	return value
}

func (srv *TextContentCreator) createRandomValue(params services_models.RandomValueCreatorParams) (string, error) {
	var value string
	switch params.ValueType {
	case sql_constants.BOOL:
		value = strconv.FormatBool(custom_utils.CreateRandomBool())
	case sql_constants.NUMBER:
		value = strconv.Itoa(custom_utils.CreateRandomNumber(params.Min, params.Max))
	case sql_constants.STRING:
		value = custom_utils.CreateRandomWord(params.Min, params.Max, false)
	case sql_constants.AUTO_INCREMENT:
		if params.IncrementFn == nil {
			return "", fmt.Errorf("[TextContentCreator] params.incrementFn can't be empty!")
		}
		value = strconv.Itoa(params.IncrementFn())
	default:
		value = "unknown"
	}

	return value, nil
}

var _ types_module.FileContentCreator = (*TextContentCreator)(nil)
