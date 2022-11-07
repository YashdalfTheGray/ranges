package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

func generateLogLine(uri, method, protocol, remote string) string {
	return fmt.Sprintf("[%s] \"%s %s %s\" %s", time.Now().UTC(), method, uri, protocol, remote)
}

func getHtmlForRange(selectedRange *RangeDetails) (string, error) {
	tmpl, parseErr := template.New("range").Parse(htmlTemplate)
	if parseErr != nil {
		return "", parseErr
	}
	buf := bytes.NewBuffer(make([]byte, 0))

	execErr := tmpl.Execute(buf, selectedRange)
	if execErr != nil {
		return "", execErr
	}

	return buf.String(), nil
}
