package output

import (
	"io"
	"os"
	"testing"

	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeStructuredOutputPassesThroughMarkdownFields(t *testing.T) {
	normalized, err := normalizeStructuredOutput(taskforge.Issue{
		Name:        "Probe",
		Description: "# Title\n\nHello **world**",
	})
	require.NoError(t, err)

	record, ok := normalized.(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, "Probe", record["name"])
	assert.Equal(t, "# Title\n\nHello **world**", record["description"])
	assert.NotContains(t, record, "description_html")
	assert.NotContains(t, record, "description_markdown")
}

func TestNormalizeStructuredOutputPassesThroughComments(t *testing.T) {
	normalized, err := normalizeStructuredOutput([]taskforge.Comment{
		{
			ID:      "1",
			Comment: "Hi ~~there~~",
		},
	})
	require.NoError(t, err)

	items, ok := normalized.([]interface{})
	require.True(t, ok)
	require.Len(t, items, 1)

	record, ok := items[0].(map[string]interface{})
	require.True(t, ok)

	assert.NotContains(t, record, "comment_html")
	assert.NotContains(t, record, "comment_stripped")
	assert.Equal(t, "Hi ~~there~~", record["comment"])
}

func TestFormatterPrintDefaultsToYAML(t *testing.T) {
	formatter := NewFormatter("", false)

	output, err := captureStdout(t, func() error {
		return formatter.Print(taskforge.Issue{
			Name:        "Probe",
			Description: "Hello **world**",
		})
	})
	require.NoError(t, err)

	assert.Contains(t, output, "name: Probe")
	assert.Contains(t, output, "description: Hello **world**")
	assert.NotContains(t, output, "description_html")
}

func TestFormatterPrintRejectsTable(t *testing.T) {
	formatter := NewFormatter("table", false)

	err := formatter.Print(taskforge.Issue{Name: "Probe"})
	require.Error(t, err)
	assert.EqualError(t, err, `invalid output format "table": supported format is yaml`)
}

func TestValidateFormat(t *testing.T) {
	require.NoError(t, ValidateFormat(""))
	require.NoError(t, ValidateFormat("yaml"))

	err := ValidateFormat("bogus")
	require.Error(t, err)
	assert.EqualError(t, err, `invalid output format "bogus": supported format is yaml`)
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer
	runErr := fn()
	require.NoError(t, writer.Close())
	os.Stdout = originalStdout

	out, readErr := io.ReadAll(reader)
	require.NoError(t, readErr)
	require.NoError(t, reader.Close())

	return string(out), runErr
}
