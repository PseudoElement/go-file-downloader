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
	srv.logger.AddLog("CreateFileContent", "Start!")
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

	srv.logger.AddLog("CreateFileContent", "End!")
	srv.logger.ShowLogs("CreateFileContent")

	return fileContent, nil
}

func (srv *TextContentCreator) CreateFileContentAsync(body interface{}) (string, error) {
	srv.logger.AddLog("CreateFileContentAsync", "Start!")
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

	srv.logger.AddLog("CreateFileContentAsync", "End!")
	srv.logger.ShowLogs("CreateFileContentAsync")

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
			if rowNumber == 0 {
				headerBuffer := bytes.NewBufferString(columnData.Name)
				srv.addSpaceOrParagraphToValue(headerBuffer, isLastColumn, columnData)
				cellsOfColumn = append(cellsOfColumn, headerBuffer.String())
			}

			value, err := srv.createRandomValue(ccr_models.RandomValueCreatorParams{
				ValueType:   columnData.Type,
				Min:         columnData.Min,
				Max:         columnData.Max,
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

			srv.addSpaceOrParagraphToValue(valueBuffer, isLastColumn, columnData)
			cellsOfColumn = append(cellsOfColumn, valueBuffer.String())
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
					headerBuffer := bytes.NewBufferString(columnData.Name)
					srv.addSpaceOrParagraphToValue(headerBuffer, isLastColumn, columnData)
					cellsOfColumn[rowNumber] = headerBuffer.String()
				}

				value, err := srv.createRandomValue(ccr_models.RandomValueCreatorParams{
					ValueType:   columnData.Type,
					Min:         columnData.Min,
					Max:         columnData.Max,
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

				srv.addSpaceOrParagraphToValue(valueBuffer, isLastColumn, columnData)

				cellsOfColumn[rowNumber+1] = valueBuffer.String()
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
		max := math.Max(float64(len(columnData.Name)), float64(columnData.Max))
		columnWidth := math.Max(max, 5)
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

var _ types_module.FileContentCreator = (*TextContentCreator)(nil)
