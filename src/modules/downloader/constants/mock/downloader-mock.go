package mock_constants

import types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"

var ROWS_COUNT = 100_000

var MOCK_SQL_COLUMNS_DATA = []types_module.SqlColumnInfo{
	types_module.SqlColumnInfo{
		Name:         "id",
		Type:         "AUTO_INCREMENT",
		IsPrimaryKey: true,
	},
	types_module.SqlColumnInfo{
		Name:              "first_name",
		Type:              "FIRST_NAME",
		NullValuesPercent: 0,
	},
	types_module.SqlColumnInfo{
		Name:              "last_name",
		Type:              "LAST_NAME",
		NullValuesPercent: 0,
	},
	types_module.SqlColumnInfo{
		Name: "is_married",
		Type: "BOOL",
	},
	types_module.SqlColumnInfo{
		Name:              "rand_num",
		Type:              "NUMBER",
		NullValuesPercent: 30,
		Min:               -1_000_000,
		Max:               10_000_000,
	},
	types_module.SqlColumnInfo{
		Name:              "rand_str",
		Type:              "STRING",
		NullValuesPercent: 10,
		Min:               1,
		Max:               30,
	},
	types_module.SqlColumnInfo{
		Name:              "rand_country",
		Type:              "COUNTRY",
		NullValuesPercent: 0,
	},
	types_module.SqlColumnInfo{
		Name:              "rand_car",
		Type:              "CAR",
		NullValuesPercent: 10,
	},
	types_module.SqlColumnInfo{
		Name: "birth_date",
		Type: "DATE",
		ForeignKeyData: types_module.ForeignKeyData{
			RefTableName:  "birthdays",
			RefColumnName: "birth_id",
		},
		NullValuesPercent: 10,
	},
}

var MOCK_TEXT_COLUMNS_DATA = []types_module.TextColumnInfo{
	types_module.TextColumnInfo{
		Name: "Id",
		Type: "AUTO_INCREMENT",
	},
	types_module.TextColumnInfo{
		Name:              "first_name",
		Type:              "FIRST_NAME",
		NullValuesPercent: 0,
	},
	types_module.TextColumnInfo{
		Name:              "last_name",
		Type:              "LAST_NAME",
		NullValuesPercent: 0,
	},
	types_module.TextColumnInfo{
		Name: "rand_bool",
		Type: "BOOL",
	},
	types_module.TextColumnInfo{
		Name:              "rand_country",
		Type:              "COUNTRY",
		NullValuesPercent: 0,
		Min:               5,
		Max:               20,
	},
	types_module.TextColumnInfo{
		Name:              "rand_str",
		Type:              "STRING",
		NullValuesPercent: 50,
		Min:               10,
		Max:               20,
	},
	types_module.TextColumnInfo{
		Name:              "rand_num",
		Type:              "NUMBER",
		NullValuesPercent: 30,
		Min:               -1_000_000,
		Max:               10_000_000,
	},
	types_module.TextColumnInfo{
		Name:              "date",
		Type:              "DATE",
		NullValuesPercent: 50,
	},
}
