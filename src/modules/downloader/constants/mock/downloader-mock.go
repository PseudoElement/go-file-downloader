package mock_constants

import types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"

var ROWS_COUNT = 40_000

var MOCK_SQL_COLUMNS_DATA = []types_module.SqlColumnInfo{
	types_module.SqlColumnInfo{
		Name:              "id",
		Type:              "AUTO_INCREMENT",
		IsPrimaryKey:      true,
		NullValuesPercent: 0,
	},
	types_module.SqlColumnInfo{
		Name:              "first_name",
		Type:              "FIRST_NAME",
		NullValuesPercent: 0,
		Min:               5,
		Max:               15,
	},
	types_module.SqlColumnInfo{
		Name:              "last_name",
		Type:              "LAST_NAME",
		NullValuesPercent: 0,
		Min:               5,
		Max:               15,
	},
	types_module.SqlColumnInfo{
		Name: "is_married",
		Type: "BOOL",
	},
	types_module.SqlColumnInfo{
		Name:              "region",
		Type:              "COUNTRY",
		NullValuesPercent: 0,
		Min:               5,
		Max:               30,
	},
	types_module.SqlColumnInfo{
		Name:              "car",
		Type:              "CAR",
		NullValuesPercent: 50,
		Min:               5,
		Max:               30,
	},
	types_module.SqlColumnInfo{
		Name: "work_id",
		Type: "STRING",
		ForeignKeyData: types_module.ForeignKeyData{
			RefTableName:  "works",
			RefColumnName: "work_id",
		},
		NullValuesPercent: 0,
		Min:               5,
		Max:               10,
	},
}

var MOCK_TEXT_COLUMNS_DATA = []types_module.TextColumnInfo{
	types_module.TextColumnInfo{
		Name:              "Id",
		Type:              "AUTO_INCREMENT",
		NullValuesPercent: 0,
	},
	types_module.TextColumnInfo{
		Name:              "Name",
		Type:              "FIRST_NAME",
		NullValuesPercent: 0,
		Min:               5,
		Max:               10,
	},
	types_module.TextColumnInfo{
		Name:              "Surname",
		Type:              "LAST_NAME",
		NullValuesPercent: 0,
		Min:               10,
		Max:               15,
	},
	types_module.TextColumnInfo{
		Name: "IsMarried",
		Type: "BOOL",
	},
	types_module.TextColumnInfo{
		Name:              "Country",
		Type:              "COUNTRY",
		NullValuesPercent: 0,
		Min:               5,
		Max:               30,
	},
	types_module.TextColumnInfo{
		Name:              "Child",
		Type:              "STRING",
		NullValuesPercent: 50,
		Min:               20,
		Max:               30,
	},
	types_module.TextColumnInfo{
		Name:              "WorkTitle",
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               5,
		Max:               10,
	},
}
