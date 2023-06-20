package zodo

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	ColorBlack     = "black"
	ColorRed       = "red"
	ColorGreen     = "green"
	ColorYellow    = "yellow"
	ColorBlue      = "blue"
	ColorMagenta   = "magenta"
	ColorCyan      = "cyan"
	ColorWhite     = "white"
	ColorHiBlack   = "hiBlack"
	ColorHiRed     = "hiRed"
	ColorHiGreen   = "hiGreen"
	ColorHiYellow  = "hiYellow"
	ColorHiBlue    = "hiBlue"
	ColorHiMagenta = "hiMagenta"
	ColorHiCyan    = "hiCyan"
	ColorHiWhite   = "hiWhite"
)

var configColorStringFuncMap = map[string]func(format string, a ...interface{}) string{
	ColorBlack:     color.BlackString,
	ColorRed:       color.RedString,
	ColorGreen:     color.GreenString,
	ColorYellow:    color.YellowString,
	ColorBlue:      color.BlueString,
	ColorMagenta:   color.MagentaString,
	ColorCyan:      color.CyanString,
	ColorWhite:     color.WhiteString,
	ColorHiBlack:   color.HiBlackString,
	ColorHiRed:     color.HiRedString,
	ColorHiGreen:   color.HiGreenString,
	ColorHiYellow:  color.HiYellowString,
	ColorHiBlue:    color.HiBlueString,
	ColorHiMagenta: color.HiMagentaString,
	ColorHiCyan:    color.HiCyanString,
	ColorHiWhite:   color.HiWhiteString,
}

func ColoredString(configColor string, format string, a ...interface{}) string {
	f := configColorStringFuncMap[configColor]
	if f != nil {
		return f(format, a...)
	}
	return fmt.Sprintf(format, a...)
}
