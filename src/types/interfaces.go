package types_module

type FileContentCreator interface {
	CreateFileContent(body interface{}) (string, error)
}
