package content_creators

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	ccr_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/constants"
	random_value_factory "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/factories"
	ccr_models "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators/models"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
	"github.com/pseudoelement/go-file-downloader/src/utils/logger"
	sql_utils "github.com/pseudoelement/go-file-downloader/src/utils/sql-utils"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type SqlContentCreator struct {
	mu     sync.Mutex
	logger logger.Logger
}

func NewSqlContentCreator(logger logger.Logger) *SqlContentCreator {
	return &SqlContentCreator{
		logger: logger,
	}
}

func (srv *SqlContentCreator) CreateFileContent(body interface{}) (string, error) {
	// srv.logger.AddLog("SQL_CreateFileContent", "Start!")
	sqlBody, ok := body.(types_module.DownloadSqlReqBody)
	if !ok {
		return "", fmt.Errorf("[SqlContentCreator] Invalid body type")
	}

	contentBuffer := new(bytes.Buffer)

	if sqlBody.NeedCreateTable {
		createTableQuery, err := srv.addTableCreationQuery(sqlBody)
		if err != nil {
			return "", err
		}
		contentBuffer.WriteString(createTableQuery)
	}

	incrementFns := make(map[string]func() int)
	for _, column := range sqlBody.ColumnsData {
		if column.Type == sql_constants.AUTO_INCREMENT {
			incrementFns[column.Name] = custom_utils.AutoIncrement(1)
		}
	}

	for i := 0; i < sqlBody.RowsCount; i++ {
		row, err := srv.addInsertRowQuery(sqlBody.ColumnsData, sqlBody.TableName, incrementFns)
		if err != nil {
			return "", err
		}

		rowWithParagraph := row + "\n\n"
		contentBuffer.WriteString(rowWithParagraph)
	}
	// srv.logger.ShowLogs("SQL_CreateFileContent")

	return contentBuffer.String(), nil
}

func (srv *SqlContentCreator) CreateFileContentAsync(body interface{}) (string, error) {
	// srv.logger.AddLog("SQL_CreateFileContentAsync", "Start")
	sqlBody, ok := body.(types_module.DownloadSqlReqBody)
	if !ok {
		return "", fmt.Errorf("[SqlContentCreator] Invalid body type")
	}

	contentBuffer := new(bytes.Buffer)

	if sqlBody.NeedCreateTable {
		createTableQuery, err := srv.addTableCreationQuery(sqlBody)
		if err != nil {
			return "", err
		}
		contentBuffer.WriteString(createTableQuery)
	}

	incrementFns := make(map[string]func() int)
	for _, column := range sqlBody.ColumnsData {
		if column.Type == sql_constants.AUTO_INCREMENT {
			incrementFns[column.Name] = custom_utils.AutoIncrement(1)
		}
	}

	errorChan := make(chan error)
	doneChan := make(chan bool)

	for i := 0; i < sqlBody.RowsCount; i++ {
		go func(index int) {
			srv.mu.Lock()
			defer srv.mu.Unlock()
			isLast := index == sqlBody.RowsCount-1

			row, err := srv.addInsertRowQuery(sqlBody.ColumnsData, sqlBody.TableName, incrementFns)
			if err != nil {
				errorChan <- err
			}

			rowWithParagraph := row + "\n\n"
			contentBuffer.WriteString(rowWithParagraph)

			if isLast {
				doneChan <- true
			}
		}(i)
	}

	select {
	case err := <-errorChan:
		return "", err
	case <-doneChan:
		// srv.logger.ShowLogs("SQL_CreateFileContentAsync")
		return contentBuffer.String(), nil
	}
}

func (srv *SqlContentCreator) addTableCreationQuery(body types_module.DownloadSqlReqBody) (string, error) {
	firstRow := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", body.TableName)
	var columns string
	for i, column := range body.ColumnsData {
		value := column.Name
		typeValue, ok := sql_constants.VALUE_TYPE_TO_SQL_TYPE[column.Type]
		if !ok {
			msg := fmt.Sprintf("[SqlContentCreator] %s - Not supported value type for sql!", column.Type)
			return "", fmt.Errorf(msg)
		}

		value += " " + typeValue

		if column.IsPrimaryKey {
			value += " PRIMARY KEY"
		}
		if column.NullValuesPercent == 0 {
			value += " NOT NULL"
		}
		if column.ForeignKeyData.RefColumnName != "" {
			value += fmt.Sprintf(" REFERENCES %s (%s)", column.ForeignKeyData.RefTableName, column.ForeignKeyData.RefColumnName)
		}

		if isLastColumn := i == len(body.ColumnsData)-1; !isLastColumn {
			value += ",\n"
		}

		columns += value
	}

	createTableQuery := firstRow + "(\n" + columns + ");\n\n"

	return createTableQuery, nil
}

func (srv *SqlContentCreator) addInsertRowQuery(columnsData []types_module.SqlColumnInfo, tableName string, incrementFns map[string]func() int) (string, error) {
	values := make([]string, 0, len(columnsData))
	columnNames := make([]string, 0, len(columnsData))

	for _, column := range columnsData {
		if isNullValue := column.NullValuesPercent > rand.Intn(100); isNullValue {
			continue
		}

		value, err := srv.createRandomValue(ccr_models.RandomValueCreatorParams{
			ValueType:   column.Type,
			Min:         column.Min,
			Max:         column.Max,
			IncrementFn: incrementFns[column.Name],
		})
		if err != nil {
			return "", err
		}

		values = append(values, value)
		columnNames = append(columnNames, column.Name)
	}

	valuesString := strings.Join(values, ", ")
	columnNamesString := strings.Join(columnNames, ", ")

	insertRowQuery := fmt.Sprintf(`INSERT INTO %s (%s)
VALUES (%s);`, tableName, columnNamesString, valuesString)

	return insertRowQuery, nil
}

func (srv *SqlContentCreator) createRandomValue(params ccr_models.RandomValueCreatorParams) (string, error) {
	var value string
	if params.ValueType == sql_constants.BOOL {
		value = string(custom_utils.CreateRandomByteForSql())
	} else {
		if val, err := random_value_factory.CommonRandomValueFactory(params); err != nil {
			return "", err
		} else {
			value = val
		}
	}

	if slice_utils_module.Contains(ccr_constants.TEXTLIKE_VALUE_TYPES, params.ValueType) {
		value = sql_utils.WrapStringInSingleQuotes(value)
	}

	return value, nil
}

var _ types_module.FileContentCreator = (*SqlContentCreator)(nil)
