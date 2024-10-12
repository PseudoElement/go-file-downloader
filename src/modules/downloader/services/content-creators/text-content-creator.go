package content_creators

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	ccr_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/constants"
	random_value_factory "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/factories"
	ccr_models "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/models"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
	"github.com/pseudoelement/go-file-downloader/src/utils/logger"
)

type TextContentCreator struct {
	mu     sync.Mutex
	logger logger.Logger
}

func NewTextContentCreator(logger logger.Logger) *TextContentCreator {
	return &TextContentCreator{
		logger: logger,
	}
}

func (srv *TextContentCreator) CreateFileContent(body interface{}) (string, error) {
	srv.logger.AddLog("TXT_CreateFileContent", "Start!")
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

	srv.logger.ShowLogs("TXT_CreateFileContent")
	return fileContent, nil
}

func (srv *TextContentCreator) CreateFileContentAsync(body interface{}) (string, error) {
	srv.logger.AddLog("TXT_CreateFileContentAsync", "Start!")

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

	srv.logger.ShowLogs("TXT_CreateFileContentAsync")

	return fileContent, nil
}

func (srv *TextContentCreator) concatAllCellsInTable(columnsWithFullCells [][]string, rowsCountWithHeader int, columnsCount int) string {
	contentBuffer := new(bytes.Buffer)
	for rowNumber := range rowsCountWithHeader {
		rowBuffer := new(bytes.Buffer)
		for columnNumber := 0; columnNumber < columnsCount; columnNumber++ {
			rowBuffer.WriteString(columnsWithFullCells[columnNumber][rowNumber])
		}
		contentBuffer.WriteString(rowBuffer.String())
	}

	return contentBuffer.String()
}

func (srv *TextContentCreator) createCellsForColumns(body types_module.DownloadTextReqBody) ([][]string, error) {
	columnsWithValuesInCells := make([][]string, 0, len(body.ColumnsData))

	for i, columnData := range body.ColumnsData {
		isLastColumn := i == len(body.ColumnsData)-1
		capacity := body.RowsCount + 1
		cellsOfColumn := make([]string, 0, capacity)
		incrementFn := custom_utils.AutoIncrement(1)

		for rowNumber := range body.RowsCount {
			min, max := srv.getDefaultMinMaxParams(columnData)
			columnWithFixedMinMax := types_module.TextColumnInfo{
				Name:              columnData.Name,
				Type:              columnData.Type,
				NullValuesPercent: columnData.NullValuesPercent,
				Min:               min,
				Max:               max,
			}

			if rowNumber == 0 {
				headerBuffer := bytes.NewBufferString(columnData.Name)
				srv.addSpaceOrParagraphToValue(headerBuffer, isLastColumn, columnWithFixedMinMax)
				cellsOfColumn = append(cellsOfColumn, headerBuffer.String())
			}

			value, err := srv.createRandomValue(ccr_models.RandomValueCreatorParams{
				ValueType:   columnWithFixedMinMax.Type,
				Min:         columnWithFixedMinMax.Min,
				Max:         columnWithFixedMinMax.Max,
				IncrementFn: incrementFn,
			})
			valueBuffer := bytes.NewBufferString(value)

			if err != nil {
				return nil, err
			}

			if isNullValue := columnData.NullValuesPercent > rand.Intn(100); isNullValue {
				valueBuffer.Reset()
				valueBuffer.WriteString("null")
			}

			srv.addSpaceOrParagraphToValue(valueBuffer, isLastColumn, columnWithFixedMinMax)
			cellsOfColumn = append(cellsOfColumn, valueBuffer.String())
		}

		columnsWithValuesInCells = append(columnsWithValuesInCells, cellsOfColumn)
	}

	return columnsWithValuesInCells, nil
}

func (srv *TextContentCreator) createCellsForColumnsAsync(body types_module.DownloadTextReqBody) ([][]string, error) {
	columns := make([][]string, len(body.ColumnsData), len(body.ColumnsData))
	errorsChan := make(chan error, len(body.ColumnsData)*body.RowsCount)

	var wg sync.WaitGroup

	for i, column := range body.ColumnsData {
		wg.Add(1)
		go func(ind int, columnData types_module.TextColumnInfo) {
			defer wg.Done()

			isLastColumn := ind == len(body.ColumnsData)-1
			cellsOfColumn := make([]string, body.RowsCount+1, body.RowsCount+1)
			incrementFn := custom_utils.AutoIncrement(1)

			for rowNumber := range body.RowsCount {
				min, max := srv.getDefaultMinMaxParams(columnData)
				columnWithFixedMinMax := types_module.TextColumnInfo{
					Name:              columnData.Name,
					Type:              columnData.Type,
					NullValuesPercent: columnData.NullValuesPercent,
					Min:               min,
					Max:               max,
				}

				if rowNumber == 0 {
					headerBuffer := bytes.NewBufferString(columnData.Name)
					srv.addSpaceOrParagraphToValue(headerBuffer, isLastColumn, columnWithFixedMinMax)
					cellsOfColumn[rowNumber] = headerBuffer.String()
				}

				value, err := srv.createRandomValue(ccr_models.RandomValueCreatorParams{
					ValueType:   columnWithFixedMinMax.Type,
					Min:         columnWithFixedMinMax.Min,
					Max:         columnWithFixedMinMax.Max,
					IncrementFn: incrementFn,
				})
				valueBuffer := bytes.NewBufferString(value)

				if err != nil {
					errorsChan <- err
					continue
				}
				if isNullValue := columnData.NullValuesPercent > rand.Intn(100); isNullValue {
					valueBuffer.Reset()
					valueBuffer.WriteString("null")
				}

				srv.addSpaceOrParagraphToValue(valueBuffer, isLastColumn, columnWithFixedMinMax)

				cellsOfColumn[rowNumber+1] = valueBuffer.String()
			}

			columns[ind] = cellsOfColumn
		}(i, column)

	}
	wg.Wait()
	close(errorsChan)

	for err := range errorsChan {
		if err != nil {
			return nil, err
		}
	}

	return columns, nil
}

func (srv *TextContentCreator) addSpaceOrParagraphToValue(valueBuffer *bytes.Buffer, isLastColumn bool, columnData types_module.TextColumnInfo) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
			fmt.Println("Recovered value ===> ", valueBuffer.String())
		}
	}()

	isHeader := valueBuffer.String() == columnData.Name

	if isHeader && isLastColumn {
		valueBuffer.WriteString("\n\n")
	} else if isLastColumn {
		valueBuffer.WriteString("\n")
	} else {
		var columnWidth float64
		if columnData.Type == sql_constants.NUMBER || columnData.Type == sql_constants.DATE {
			// cause Min/Max can be over 1_000_000++
			columnWidth = math.Max(float64(len(valueBuffer.String())), float64(len(columnData.Name)))
			columnWidth = math.Max(float64(ccr_constants.VALUE_LENGTHS[columnData.Type]), columnWidth)
		} else {
			columnWidth = math.Max(float64(len(columnData.Name)), float64(columnData.Max))
		}
		valueLen := len(valueBuffer.String())
		spaces := strings.Repeat(" ", int(columnWidth)-valueLen+2)
		valueBuffer.WriteString(spaces)
	}
}

func (srv *TextContentCreator) createRandomValue(params ccr_models.RandomValueCreatorParams) (string, error) {
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

func (srv *TextContentCreator) getDefaultMinMaxParams(column types_module.TextColumnInfo) (int64, int64) {
	var max int64 = column.Max
	var min int64 = column.Min
	if column.Max == 0 {
		if column.Type == sql_constants.DATE {
			max = 1727500009835
		} else {
			max = 20
		}
	}
	if column.Min == 0 {
		if column.Type == sql_constants.DATE {
			min = 172000000983
		} else {
			min = 1
		}
	}
	return min, max
}

var _ types_module.FileContentCreator = (*TextContentCreator)(nil)
