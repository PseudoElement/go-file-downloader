package types_module

import "os"

type FileContentCreator interface {
	CreateFileContent(args ...interface{}) (*os.File, error)
}
