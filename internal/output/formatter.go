package output

import (
	"fmt"
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
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", color.GreenString("✓ %s", msg))
}

func Error(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", color.RedString("✗ %s", msg))
}

func Warning(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", color.YellowString("⚠ %s", msg))
}

func Info(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", color.CyanString("ℹ %s", msg))
}
