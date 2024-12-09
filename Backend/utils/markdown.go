package utils

import (
    "bytes"

    "github.com/yuin/goldmark"
)

func ConvertMarkdownToHTML(markdown string) (string, error) {
    var buf bytes.Buffer
    if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
        return "", err
    }
    return buf.String(), nil
}
