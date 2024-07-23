package downloader_module

import (
	"os"
	"slices"
	"strconv"
	"strings"

	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
)

type DownloaderService struct{}

func NewDownloaderService() *DownloaderService {
	srv := &DownloaderService{}
	return srv
}

func (srv *DownloaderService) CreateTxtFileWithContentSync(body DownloadTextReqBody) (*os.File, error) {
	rowsCountWithHeader := body.RowsCount + 1
	columnsWithValuesInCells := srv.createCellsForColumns(body)
	columnsWithFullCells := srv.createCellsWithSpacesAndParagraphs(columnsWithValuesInCells)
	fileContent := srv.createFileContent(columnsWithFullCells, rowsCountWithHeader, len(body.ColumnsData))

	file, err := custom_utils.CreateFile(body.DocName, body.DocType, fileContent)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (srv *DownloaderService) createFileContent(columnsWithFullCells [][]string, rowsCountWithHeader int, columnsCount int) string {
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

func (srv *DownloaderService) createCellsWithSpacesAndParagraphs(columnsWithValuesInCells [][]string) [][]string {
	columnsWithFullCells := [][]string{}

	for ind, columnWithValues := range columnsWithValuesInCells {
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

		columnsWithFullCells = append(columnsWithFullCells, columnWithSpaces)
	}

	return columnsWithFullCells
}

func (srv *DownloaderService) createCellsForColumns(body DownloadTextReqBody) [][]string {
	columnsWithValuesInCells := [][]string{}

	for _, columnData := range body.ColumnsData {
		cellsOfColumn := []string{}
		for j := range body.RowsCount {
			if j == 0 {
				cellsOfColumn = append(cellsOfColumn, columnData.Name)
			}

			var value string
			switch columnData.Type {
			case BOOL:
				value = strconv.FormatBool(custom_utils.CreateRandomBool())
			case NUMBER:
				value = strconv.Itoa(custom_utils.CreateRandomNumber(columnData.Min, columnData.Max))
			case STRING:
				value = custom_utils.CreateRandomWord(columnData.Min, columnData.Max, false)
			default:
				value = "empty"
			}
			cellsOfColumn = append(cellsOfColumn, value)
		}

		columnsWithValuesInCells = append(columnsWithValuesInCells, cellsOfColumn)
	}

	return columnsWithValuesInCells
}
