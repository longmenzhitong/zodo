package zodo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/go-gomail/gomail"
	"github.com/go-redis/redis"
	"github.com/jedib0t/go-pretty/v6/table"
)

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

func PrintTable(header *table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(Config.Todo.TableMaxLength)
	if header != nil {
		t.AppendHeader(*header)
	}
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

func WriteLinesToClipboard(lines []string) error {
	var text string
	for i, line := range lines {
		text += line
		if i != len(lines)-1 {
			text += "\n"
		}
	}
	err := clipboard.WriteAll(text)
	if err != nil {
		return err
	}

	fmt.Println("(Copied.)")
	return nil
}

func WriteLineToClipboard(line string) error {
	if line == "" {
		return nil
	}
	return WriteLinesToClipboard([]string{line})
}

func PrintStartMsg(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s", ColoredString(ColorGreen, "==>"), msg)
}

func PrintDoneMsg(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s", ColoredString(ColorBlue, "==>"), msg)
}
