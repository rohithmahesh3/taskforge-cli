package output

import (
	"os"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Formatter struct {
	Format  string
	NoColor bool
	Wide    bool
}

func NewFormatter(format string, noColor bool) *Formatter {
	if noColor {
		color.NoColor = true
	}

	return &Formatter{
		Format:  NormalizeFormat(format),
		NoColor: noColor,
	}
}

func (f *Formatter) Print(data interface{}) error {
	format := NormalizeFormat(f.Format)
	if err := ValidateFormat(format); err != nil {
		return err
	}

	return f.printYAML(data)
}

func (f *Formatter) printYAML(data interface{}) error {
	normalized, err := normalizeStructuredOutput(data)
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(os.Stdout)
	defer func() {
		_ = encoder.Close()
	}()
	return encoder.Encode(normalized)
}

func Success(msg string) {
	color.Green("✓ %s", msg)
}

func Error(msg string) {
	color.Red("✗ %s", msg)
}

func Warning(msg string) {
	color.Yellow("⚠ %s", msg)
}

func Info(msg string) {
	color.Cyan("ℹ %s", msg)
}
