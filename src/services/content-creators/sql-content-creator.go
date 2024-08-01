package content_creators

import (
	"fmt"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/constants/sql"
	types_module "github.com/pseudoelement/go-file-downloader/src/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
)

type SqlContentCreator struct{}

func NewSqlContentCreator() *SqlContentCreator {
	return &SqlContentCreator{}
}

func (srv *SqlContentCreator) CreateFileContent(body interface{}) (string, error) {
	sqlBody, ok := body.(types_module.DownloadSqlReqBody)
	if !ok {
		return "", fmt.Errorf("[SqlContentCreator] Invalid body type")
	}
	var sqlFileContent string

	if sqlBody.NeedCreateTable {
		createTableQuery, err := srv.addTableCreationQuery(sqlBody)
		if err != nil {
			return "", err
		}
		sqlFileContent += createTableQuery
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
		sqlFileContent += row + "\n\n"
	}

	return sqlFileContent, nil
}

func (srv *SqlContentCreator) addTableCreationQuery(body types_module.DownloadSqlReqBody) (string, error) {
	firstRow := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", body.TableName)
	var columns string
	for _, column := range body.ColumnsData {
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

		value += ",\n"
		columns += value
	}

	createTableQuery := firstRow + "(\n" + columns + ");\n\n"

	return createTableQuery, nil
}

func (srv *SqlContentCreator) addInsertRowQuery(columnsData []types_module.SqlColumnInfo, tableName string, incrementFns map[string]func() int) (string, error) {
	var values string
	var columnNames string
	for i, column := range columnsData {
		value, err := custom_utils.CreateRandomValueConvertedToString(custom_utils.RandomValueCreatorParams{
			ValueType:   column.Type,
			Min:         column.Min,
			Max:         column.Max,
			IncrementFn: incrementFns[column.Name],
			IsSqlValue:  true,
		})
		if err != nil {
			return "", err
		}

		if i == len(columnsData)-1 {
			values += value
			columnNames += column.Name
		} else {
			values += value + ", "
			columnNames += column.Name + ", "
		}
	}

	insertRowQuery := fmt.Sprintf(`INSERT INTO %s (%s)
VALUES (%s);`, tableName, columnNames, values)

	return insertRowQuery, nil
}

var _ types_module.FileContentCreator = (*SqlContentCreator)(nil)
