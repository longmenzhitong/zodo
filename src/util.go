package zodo

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-gomail/gomail"
	"github.com/go-redis/redis"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"os/exec"
	"strings"
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

var redisClient *redis.Client

func Redis() *redis.Client {
	if redisClient == nil {
		rc := redis.NewClient(&redis.Options{
			Addr:     Config.Storage.Redis.Address,
			Password: Config.Storage.Redis.Password,
			DB:       Config.Storage.Redis.Db,
		})
		_, err := rc.Ping().Result()
		if err != nil {
			panic(err)
		}
		redisClient = rc
	}
	return redisClient
}

func SendEmail(title, text string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", Config.Email.From, "ZODO")
	m.SetHeader("To", Config.Email.To...)
	m.SetHeader("Subject", title)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(Config.Email.Server, Config.Email.Port, Config.Email.From, Config.Email.Auth)
	return d.DialAndSend(m)
}

func PrintTable(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(Config.Table.MaxLen)
	t.AppendHeader(header)
	for _, row := range rows {
		t.AppendRow(row)
		t.AppendSeparator()
	}
	t.Render()
}

func CurrentGitBranch() (string, error) {
	// 执行 git 命令来获取当前分支名称
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func ColoredString(configColor string, format string, a ...interface{}) string {
	f := configColorStringFuncMap[configColor]
	if f != nil {
		return f(format, a...)
	}
	return fmt.Sprintf(format, a...)
}
