package zodo

import (
	"github.com/go-gomail/gomail"
	"github.com/go-redis/redis"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
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
