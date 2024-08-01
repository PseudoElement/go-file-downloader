package types_module

type FileContentCreator interface {
	CreateFileContent(body interface{}) (string, error)
}

type RandomValueFactory interface {
	CreateRandomBool() bool
	CreateRandomNumber(min int, max int) int
	CreateRandomWord(minLength int, maxLength int, startUpperCase bool) string
}
