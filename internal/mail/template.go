package mail

import (
	"io"
	"os"
	"strings"
)

type TypeTemplate string

const (
	FreeRequest TypeTemplate = "freeRequest.html"
)

func LoadTemplate(template TypeTemplate) (*string, error) {
	startPath, _ := os.Getwd()
	file, err := os.Open(startPath + "/internal/mail/template/" + string(template))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var builder strings.Builder
	if _, err := io.Copy(&builder, file); err != nil {
		return nil, err
	}

	strContent := builder.String()

	return &strContent, nil
}
