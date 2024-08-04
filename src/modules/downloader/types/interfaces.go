package types_module

type FileContentCreator interface {
	CreateFileContent(body interface{}) (string, error)
	CreateFileContentAsync(body interface{}) (string, error)
}
