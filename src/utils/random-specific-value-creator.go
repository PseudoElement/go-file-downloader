package custom_utils

import (
	"math/rand"

	common_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/common"
)

func CreateRandomFirstName(min int, max int) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.FIRST_NAMES_ARRAY))
	firstName := common_constants.FIRST_NAMES_ARRAY[randomIndex]
	if len(firstName) > max || len(firstName) < min {
		return CreateRandomFirstName(min, max)
	} else {
		return firstName
	}
}

func CreateRandomLastName(min int, max int) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.LAST_NAMES_ARRAY))
	lastName := common_constants.LAST_NAMES_ARRAY[randomIndex]
	if len(lastName) > max || len(lastName) < min {
		return CreateRandomLastName(min, max)
	} else {
		return lastName
	}
}

func CreateRandomCarName(min int, max int) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.CARS_ARRAY))
	carName := common_constants.CARS_ARRAY[randomIndex]
	if len(carName) > max || len(carName) < min {
		return CreateRandomCarName(min, max)
	} else {
		return carName
	}
}

func CreateRandomCountryName(min int, max int) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.COUNTRIES_ARRAY))
	countryName := common_constants.COUNTRIES_ARRAY[randomIndex]
	if len(countryName) > max || len(countryName) < min {
		return CreateRandomCountryName(min, max)
	} else {
		return countryName
	}
}

func CreateRandomWork(min int, max int) string {
	min, max = setDefaultMinMaxIfZeros(min, max)
	randomIndex := rand.Intn(len(common_constants.WORKS_ARRAY))
	work := common_constants.WORKS_ARRAY[randomIndex]
	if len(work) > max || len(work) < min {
		return CreateRandomWork(min, max)
	} else {
		return work
	}
}

func setDefaultMinMaxIfZeros(min int, max int) (int, int) {
	if min == 0 {
		min = 5
	}
	if max == 0 {
		max = 20
	}
	return min, max
}
