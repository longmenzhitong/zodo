package zodo

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Deploy(service, env, branch string, checkCode bool) error {
	// 检查配置
	fmt.Println("Check config...")
	fmt.Printf("Url       : %s\n", boolToSymbol(Config.Jenkins.Url != ""))
	if Config.Jenkins.Url == "" {
		return &InvalidConfigError{Message: "jenkins.url doesn't exist"}
	}
	fmt.Printf("Username  : %s\n", boolToSymbol(Config.Jenkins.Username != ""))
	if Config.Jenkins.Username == "" {
		return &InvalidConfigError{Message: "jenkins.username doesn't exist"}
	}
	fmt.Printf("Password  : %s\n", boolToSymbol(Config.Jenkins.Password != ""))
	if Config.Jenkins.Password == "" {
		return &InvalidConfigError{Message: "jenkins.password doesn't exist"}
	}
	fmt.Println("Check done.")

	// 检查参数
	fmt.Println("Check params...")
	if service == "" {
		service = strings.ToUpper(CurrentDirName())
	} else {
		service = strings.ToUpper(service)
	}
	fmt.Printf("Service   : %s\n", service)
	if env == "" {
		fmt.Println("Please input the env:")
		env = readString()
		if env == "" {
			return &InvalidInputError{Message: "env must not be empty"}
		}
	}
	fmt.Printf("Env       : %s\n", env)
	if branch == "" {
		b, err := CurrentGitBranch()
		if err != nil {
			return err
		}
		branch = b
	}
	fmt.Printf("Branch    : %s\n", branch)
	fmt.Printf("CheckCode : %v\n", checkCode)
	fmt.Println("Check done.")

	// 确认发布
	fmt.Println("Sure to deploy? [y/n]")
	input := strings.ToLower(readString())
	if input != "y" {
		return &CancelledError{}
	}
	fmt.Println("Start deploy...")

	// 构建请求
	jenkinsUrl := fmt.Sprintf("%s/job/%s/buildWithParameters", Config.Jenkins.Url, service)
	requestBody := url.Values{
		"BUILD_BRANCH":  {branch},
		"SERVERNAME":    {env},
		"IS_CHECK_CODE": {strings.ToUpper(boolToText(checkCode))},
	}
	req, err := http.NewRequest("POST", jenkinsUrl, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(Config.Jenkins.Username, Config.Jenkins.Password)

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBody := string(body)
	if respBody != "" {
		fmt.Println(respBody)
	} else {
		fmt.Println("Deploy done.")
		fmt.Println("To check the progress, visit:")
		fmt.Printf("%s/job/%s\n", Config.Jenkins.Url, service)
	}
	return nil
}

func boolToSymbol(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}

func boolToText(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
