package xtextract

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/grokify/awsgo/xtextract/ocrutil"
	"github.com/grokify/mogo/pointer"
)

type Blocks []*textract.Block

/*
func (b Blocks) Lines(blockTypes []string, textTypes []string) []string {
	var lines []string
	for _, bi := range b {
		s := strings.TrimSpace(pointer.ToString(bi.Text))
		if s != "" {
			lines = append(lines, s)
		}
	}
	return lines
}
*/

func (b Blocks) LinesByBlockText() map[string][]string {
	m := map[string][]string{}
	for _, bi := range b {
		key := strings.Join(
			[]string{
				strings.TrimSpace(pointer.ToString(bi.BlockType)),
				strings.TrimSpace(pointer.ToString(bi.TextType)),
			},
			"__")
		s := strings.TrimSpace(pointer.ToString(bi.Text))
		m[key] = append(m[key], s)
	}
	return m
}

func (b Blocks) TextResults() ocrutil.TextResults {
	tr := ocrutil.TextResults{
		Lines: b.Lines(),
	}
	m := b.LinesByBlockText()
	if w, ok := m["WORD__PRINTED"]; ok {
		tr.WordsPrinted = w
	}
	return tr
}

func (b Blocks) Lines() []string {
	lines := []string{}
	for _, bi := range b {
		if strings.TrimSpace(pointer.ToString(bi.BlockType)) != BlockTypeLine {
			continue
		}
		s := strings.TrimSpace(pointer.ToString(bi.Text))
		lines = append(lines, s)
	}
	return lines
}

/*
type TextResults struct {
	OCRService   string
	OCRDateTime  time.Time
	Lines        []string
	WordsPrinted []string
}

*/
