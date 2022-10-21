package color

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
	purple = "\033[35m"
	gray   = "\033[37m"
	cyan   = "\033[36m"
	white  = "\033[97m"
)

func Red(s string) string {
	return red + s + reset
}

func Yellow(s string) string {
	return yellow + s + reset
}

func Green(s string) string {
	return green + s + reset
}

func Blue(s string) string {
	return blue + s + reset
}

func Purple(s string) string {
	return purple + s + reset
}

func Gray(s string) string {
	return gray + s + reset
}
