package content_creators

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	random_value_factory "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/factories"
	services_models "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/models"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
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
	columnsWithValuesInCells := make([][]string, 0, len(body.ColumnsData))

	for i, columnData := range body.ColumnsData {
		isLastColumn := i == len(body.ColumnsData)-1
		capacity := body.RowsCount + 1
		cellsOfColumn := make([]string, 0, capacity)
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
	columnsWithValuesInCells := make([][]string, 0, len(body.ColumnsData))
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
	if params.ValueType == sql_constants.BOOL {
		value = strconv.FormatBool(custom_utils.CreateRandomBool())
	} else {
		if val, err := random_value_factory.CommonRandomValueFactory(params); err != nil {
			return "", err
		} else {
			value = val
		}
	}

	return value, nil
}

var _ types_module.FileContentCreator = (*TextContentCreator)(nil)
