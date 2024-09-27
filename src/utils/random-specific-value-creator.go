package custom_utils

import (
	"math/rand"

	common_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/common"
)

func CreateRandomFirstName(min int64, max int64) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.FIRST_NAMES_ARRAY))
	firstName := common_constants.FIRST_NAMES_ARRAY[randomIndex]
	if len(firstName) > int(max) || len(firstName) < int(min) {
		return CreateRandomFirstName(min, max)
	} else {
		return firstName
	}
}

func CreateRandomLastName(min int64, max int64) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.LAST_NAMES_ARRAY))
	lastName := common_constants.LAST_NAMES_ARRAY[randomIndex]
	if len(lastName) > int(max) || len(lastName) < int(min) {
		return CreateRandomLastName(min, max)
	} else {
		return lastName
	}
}

func CreateRandomCarName(min int64, max int64) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.CARS_ARRAY))
	carName := common_constants.CARS_ARRAY[randomIndex]
	if len(carName) > int(max) || len(carName) < int(min) {
		return CreateRandomCarName(min, max)
	} else {
		return carName
	}
}

func CreateRandomCountryName(min int64, max int64) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.COUNTRIES_ARRAY))
	countryName := common_constants.COUNTRIES_ARRAY[randomIndex]
	if len(countryName) > int(max) || len(countryName) < int(min) {
		return CreateRandomCountryName(min, max)
	} else {
		return countryName
	}
}

func CreateRandomWork(min int64, max int64) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.WORKS_ARRAY))
	work := common_constants.WORKS_ARRAY[randomIndex]
	if len(work) > int(max) || len(work) < int(min) {
		return CreateRandomWork(min, max)
	} else {
		return work
	}
}

func setDefaultMinMaxIfZeros(min int64, max int64) (int64, int64) {
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 20
	}
	return min, max
}
