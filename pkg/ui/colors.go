package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	// Color presets
	Green  = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	White  = color.New(color.FgWhite).SprintFunc()
	Gray   = color.New(color.FgHiBlack).SprintFunc()

	// Bold colors
	BoldGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()
	BoldYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	BoldRed    = color.New(color.FgRed, color.Bold).SprintFunc()
	BoldBlue   = color.New(color.FgBlue, color.Bold).SprintFunc()
)

// Success prints a success message
func Success(msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", Green("✓"), fmt.Sprintf(msg, args...))
}

// Warning prints a warning message
func Warning(msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", Yellow("⚠"), fmt.Sprintf(msg, args...))
}

// Error prints an error message
func Error(msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", Red("✗"), fmt.Sprintf(msg, args...))
}

// Info prints an info message
func Info(msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", Blue("ℹ"), fmt.Sprintf(msg, args...))
}

// Skipped prints a skipped message
func Skipped(msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", Gray("⊘"), fmt.Sprintf(msg, args...))
}

// Title prints a title
func Title(msg string, args ...interface{}) {
	fmt.Printf("\n%s %s\n\n", BoldBlue("==="), fmt.Sprintf(msg, args...))
}

// Step prints a step message
func Step(step, total int, msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", Cyan(fmt.Sprintf("[%d/%d]", step, total)), fmt.Sprintf(msg, args...))
}

// Println prints colored text
func Println(colorFn func(...interface{}) string, msg string, args ...interface{}) {
	fmt.Println(colorFn(fmt.Sprintf(msg, args...)))
}

// Printf prints formatted colored text
func Printf(colorFn func(...interface{}) string, msg string, args ...interface{}) {
	fmt.Print(colorFn(fmt.Sprintf(msg, args...)))
}
