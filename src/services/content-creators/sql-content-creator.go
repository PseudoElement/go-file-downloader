package content_creators

import (
	"fmt"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/constants/sql"
	types_module "github.com/pseudoelement/go-file-downloader/src/types"
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

	return sqlFileContent, nil
}

func (srv *SqlContentCreator) addTableCreationQuery(body types_module.DownloadSqlReqBody) (string, error) {
	firstRow := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", body.DocName)
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
		value += ",\n"
		columns += value
	}

	createTableQuery := firstRow + "(\n" + columns + ");\n\n"

	return createTableQuery, nil
}

var _ types_module.FileContentCreator = (*SqlContentCreator)(nil)
