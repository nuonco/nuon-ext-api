package tui

import "github.com/charmbracelet/lipgloss"

// Adaptive colors — internal light/dark values

var (
	lightPrimaryColor        = lipgloss.CompleteColor{TrueColor: "#8040BF", ANSI256: "99", ANSI: "5"}
	lightSecondaryColor      = lipgloss.CompleteColor{TrueColor: "#527FE8", ANSI256: "69", ANSI: "12"}
	lightAccentColor         = lipgloss.CompleteColor{TrueColor: "#D6B0FC", ANSI256: "183", ANSI: "13"}
	lightTextColor           = lipgloss.CompleteColor{TrueColor: "#000000", ANSI256: "0", ANSI: "0"}
	lightSubtleColor         = lipgloss.CompleteColor{TrueColor: "#C3C3C3", ANSI256: "250", ANSI: "7"}
	lightSuccessColor        = lipgloss.CompleteColor{TrueColor: "#439B92", ANSI256: "72", ANSI: "10"}
	lightWarningColor        = lipgloss.CompleteColor{TrueColor: "#FCA04A", ANSI256: "214", ANSI: "11"}
	lightErrorColor          = lipgloss.CompleteColor{TrueColor: "#991B1B", ANSI256: "88", ANSI: "1"}
	lightInfoColor           = lipgloss.CompleteColor{TrueColor: "#527FE8", ANSI256: "69", ANSI: "12"}
	lightBorderActiveColor   = lipgloss.CompleteColor{TrueColor: "#8040BF", ANSI256: "99", ANSI: "5"}
	lightBorderInactiveColor = lipgloss.CompleteColor{TrueColor: "#C3C3C3", ANSI256: "250", ANSI: "7"}

	darkPrimaryColor        = lipgloss.CompleteColor{TrueColor: "#D6B0FC", ANSI256: "183", ANSI: "13"}
	darkSecondaryColor      = lipgloss.CompleteColor{TrueColor: "#99B7FF", ANSI256: "111", ANSI: "12"}
	darkAccentColor         = lipgloss.CompleteColor{TrueColor: "#8040BF", ANSI256: "99", ANSI: "5"}
	darkTextColor           = lipgloss.CompleteColor{TrueColor: "#FFFFFF", ANSI256: "15", ANSI: "5"}
	darkSubtleColor         = lipgloss.CompleteColor{TrueColor: "#B9B9B9", ANSI256: "249", ANSI: "7"}
	darkSuccessColor        = lipgloss.CompleteColor{TrueColor: "#5BBFB5", ANSI256: "80", ANSI: "10"}
	darkWarningColor        = lipgloss.CompleteColor{TrueColor: "#FFBD7F", ANSI256: "223", ANSI: "11"}
	darkErrorColor          = lipgloss.CompleteColor{TrueColor: "#FF8383", ANSI256: "210", ANSI: "9"}
	darkInfoColor           = lipgloss.CompleteColor{TrueColor: "#527FE8", ANSI256: "69", ANSI: "12"}
	darkBorderActiveColor   = lipgloss.CompleteColor{TrueColor: "#8040BF", ANSI256: "99", ANSI: "5"}
	darkBorderInactiveColor = lipgloss.CompleteColor{TrueColor: "#4F4F4F", ANSI256: "238", ANSI: "8"}
)

// Exported adaptive colors

var (
	PrimaryColor        = lipgloss.CompleteAdaptiveColor{Light: lightPrimaryColor, Dark: darkPrimaryColor}
	SecondaryColor      = lipgloss.CompleteAdaptiveColor{Light: lightSecondaryColor, Dark: darkSecondaryColor}
	AccentColor         = lipgloss.CompleteAdaptiveColor{Light: lightAccentColor, Dark: darkAccentColor}
	TextColor           = lipgloss.CompleteAdaptiveColor{Light: lightTextColor, Dark: darkTextColor}
	SubtleColor         = lipgloss.CompleteAdaptiveColor{Light: lightSubtleColor, Dark: darkSubtleColor}
	SuccessColor        = lipgloss.CompleteAdaptiveColor{Light: lightSuccessColor, Dark: darkSuccessColor}
	WarningColor        = lipgloss.CompleteAdaptiveColor{Light: lightWarningColor, Dark: darkWarningColor}
	ErrorColor          = lipgloss.CompleteAdaptiveColor{Light: lightErrorColor, Dark: darkErrorColor}
	InfoColor           = lipgloss.CompleteAdaptiveColor{Light: lightInfoColor, Dark: darkInfoColor}
	BorderActiveColor   = lipgloss.CompleteAdaptiveColor{Light: lightBorderActiveColor, Dark: darkBorderActiveColor}
	BorderInactiveColor = lipgloss.CompleteAdaptiveColor{Light: lightBorderInactiveColor, Dark: darkBorderInactiveColor}
)

// Text styles

var (
	TextPrimary   = lipgloss.NewStyle().Foreground(PrimaryColor)
	TextSecondary = lipgloss.NewStyle().Foreground(SecondaryColor)
	TextAccent    = lipgloss.NewStyle().Foreground(AccentColor)
	TextDefault   = lipgloss.NewStyle().Foreground(TextColor)
	TextSubtle    = lipgloss.NewStyle().Foreground(SubtleColor)
	TextSuccess   = lipgloss.NewStyle().Foreground(SuccessColor)
	TextWarning   = lipgloss.NewStyle().Foreground(WarningColor)
	TextError     = lipgloss.NewStyle().Foreground(ErrorColor)
	TextInfo      = lipgloss.NewStyle().Foreground(InfoColor)
	TextBold      = lipgloss.NewStyle().Bold(true)
	TextDim       = lipgloss.NewStyle().Foreground(SubtleColor)
)

// Method color mapping for API browser.
// Uses ANSI256 colors directly to avoid TrueColor escape sequences
// being broken by the bubbles list delegate's text truncation.

func MethodStyle(method string) lipgloss.Style {
	switch method {
	case "GET":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("80")).Bold(true)
	case "POST":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
	case "PUT", "PATCH":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	case "DELETE":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("210")).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("249"))
	}
}
