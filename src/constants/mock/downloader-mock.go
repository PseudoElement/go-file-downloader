package mock_constants

import types_module "github.com/pseudoelement/go-file-downloader/src/types"

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
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               5,
		Max:               10,
	},
	types_module.SqlColumnInfo{
		Name:              "last_name",
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               10,
		Max:               15,
	},
	types_module.SqlColumnInfo{
		Name: "is_married",
		Type: "BOOL",
	},
	types_module.SqlColumnInfo{
		Name:              "region",
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               17,
		Max:               30,
	},
	types_module.SqlColumnInfo{
		Name:              "child",
		Type:              "STRING",
		NullValuesPercent: 90,
		Min:               20,
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
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               5,
		Max:               10,
	},
	types_module.TextColumnInfo{
		Name:              "Surname",
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               10,
		Max:               15,
	},
	types_module.TextColumnInfo{
		Name: "IsMarried",
		Type: "BOOL",
	},
	types_module.TextColumnInfo{
		Name:              "Region",
		Type:              "STRING",
		NullValuesPercent: 0,
		Min:               17,
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
