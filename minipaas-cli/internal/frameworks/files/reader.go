package files

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"

	"github.com/yunier-rojas/minipaas/minipaas-cli/internal/core/usecases"
)

type ReaderService struct{}

func NewReaderService() *ReaderService {
	return &ReaderService{}
}

func (r *ReaderService) ReadContent(file, name string) (string, []byte, error) {
	if file == "" && name == "" {
		return "", nil, errors.New("when no file is provided, --name is mandatory")
	}

	var content []byte
	var err error
	if file != "" {
		content, err = os.ReadFile(file)
	} else {
		content, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return "", nil, err
	}

	if name != "" {
		return name, content, nil
	}

	return filepath.Base(file), content, nil
}

var _ usecases.InputPort = (*ReaderService)(nil)
