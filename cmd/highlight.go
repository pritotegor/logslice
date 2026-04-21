package cmd

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// LevelColors maps log level strings to ANSI color codes.
var LevelColors = map[string]string{
	"error":   ColorRed,
	"err":     ColorRed,
	"warn":    ColorYellow,
	"warning": ColorYellow,
	"info":    ColorGreen,
	"debug":   ColorCyan,
}

// HighlightLevel wraps a log level value with its corresponding ANSI color.
// If the level is unrecognized or color is disabled, the value is returned as-is.
func HighlightLevel(level string, colorEnabled bool) string {
	if !colorEnabled {
		return level
	}
	key := strings.ToLower(strings.TrimSpace(level))
	if color, ok := LevelColors[key]; ok {
		return fmt.Sprintf("%s%s%s%s", ColorBold, color, level, ColorReset)
	}
	return level
}

// HighlightField wraps a field key with cyan color for emphasis.
func HighlightField(key string, colorEnabled bool) string {
	if !colorEnabled {
		return key
	}
	return fmt.Sprintf("%s%s%s", ColorCyan, key, ColorReset)
}

// HighlightMatch wraps a matched substring within a string with bold+yellow color.
func HighlightMatch(text, match string, colorEnabled bool) string {
	if !colorEnabled || match == "" {
		return text
	}
	highlighted := fmt.Sprintf("%s%s%s%s", ColorBold, ColorYellow, match, ColorReset)
	return strings.ReplaceAll(text, match, highlighted)
}
