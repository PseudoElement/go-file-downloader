package content_creators

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"

	types_module "github.com/pseudoelement/go-file-downloader/src/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
)

type TextContentCreator struct{}

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
	columnsWithFullCells := srv.createCellsWithSpacesAndParagraphs(columnsWithValuesInCells)
	fileContent := srv.concatAllCellsInTable(columnsWithFullCells, rowsCountWithHeader, len(textBody.ColumnsData))

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
	var err error
	columnsWithValuesInCells := custom_utils.Map(body.ColumnsData, func(columnData types_module.TextColumnInfo, ind int) []string {
		cellsOfColumn := []string{}
		incrementFn := custom_utils.AutoIncrement(1)
		for j := range body.RowsCount {
			if j == 0 {
				cellsOfColumn = append(cellsOfColumn, columnData.Name)
			}

			value, e := custom_utils.CreateRandomValueConvertedToString(custom_utils.RandomValueCreatorParams{
				ValueType:   columnData.Type,
				Min:         columnData.Min,
				Max:         columnData.Max,
				IncrementFn: incrementFn,
				IsSqlValue:  false,
			})
			err = e

			if isNullValue := columnData.NullValuesPercent > rand.Intn(100); isNullValue {
				value = "null"
			}

			cellsOfColumn = append(cellsOfColumn, value)
		}

		return cellsOfColumn
	})

	if err != nil {
		return nil, err
	}

	return columnsWithValuesInCells, nil
}

func (srv *TextContentCreator) createCellsWithSpacesAndParagraphs(columnsWithValuesInCells [][]string) [][]string {
	columnsWithFullCells := custom_utils.Map(columnsWithValuesInCells, func(columnWithValues []string, ind int) []string {
		columnWithSpaces := []string{}

		isLastColumn := ind == len(columnsWithValuesInCells)-1
		if !isLastColumn {
			valuesLengths := custom_utils.Map(columnWithValues, func(cellValue string, ind int) int {
				return len(cellValue)
			})
			// most broad word in column + 2 space between this column and the next
			columnWidth := slices.Max(valuesLengths) + 2
			for _, cellValue := range columnWithValues {
				valueLength := len(cellValue)
				spaces := strings.Repeat(" ", columnWidth-valueLength)
				cellValueWithSpace := cellValue + spaces

				columnWithSpaces = append(columnWithSpaces, cellValueWithSpace)
			}
		} else {
			for _, cellValue := range columnWithValues {
				cellValueWithParagraph := cellValue + "\n"

				columnWithSpaces = append(columnWithSpaces, cellValueWithParagraph)

			}
		}

		return columnWithSpaces
	})

	return columnsWithFullCells
}

var _ types_module.FileContentCreator = (*TextContentCreator)(nil)
