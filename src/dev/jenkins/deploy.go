package jenkins

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	zodo "zodo/src"
)

const (
	stageStatusSuccess    = "SUCCESS"
	stageStatusInProgress = "IN_PROGRESS"
	stageStatusAborted    = "ABORTED"
	stageStatusFailed     = "FAILED"
)

func Deploy(job, server, branch string, checkCode bool) error {
	// 检查配置
	fmt.Println("Check config...")
	fmt.Printf("Url       : %s\n", boolToSymbol(zodo.Config.Jenkins.Url != ""))
	if zodo.Config.Jenkins.Url == "" {
		return &zodo.InvalidConfigError{Message: "jenkins.url doesn't exist"}
	}
	fmt.Printf("Username  : %s\n", boolToSymbol(zodo.Config.Jenkins.Username != ""))
	if zodo.Config.Jenkins.Username == "" {
		return &zodo.InvalidConfigError{Message: "jenkins.username doesn't exist"}
	}
	fmt.Printf("Password  : %s\n", boolToSymbol(zodo.Config.Jenkins.Password != ""))
	if zodo.Config.Jenkins.Password == "" {
		return &zodo.InvalidConfigError{Message: "jenkins.password doesn't exist"}
	}
	fmt.Println("Check done.")

	// 检查参数
	fmt.Println("Check params...")
	fmt.Printf("Job       : %s\n", job)
	if server == "" {
		fmt.Println("Please input server:")
		server = zodo.ReadString()
		if server == "" {
			return &zodo.InvalidInputError{Message: "server must not be empty"}
		}
	}
	fmt.Printf("Server    : %s\n", server)
	fmt.Printf("Branch    : %s\n", branch)
	fmt.Printf("CheckCode : %v\n", checkCode)
	fmt.Println("Check done.")

	// 确认构建
	fmt.Println("Sure to deploy? [y/n]")
	input := strings.ToLower(zodo.ReadString())
	if input != "y" {
		return &zodo.CancelledError{}
	}

	// 发起构建
	err := startBuild(job, server, branch, checkCode)
	if err != nil {
		return err
	}

	if zodo.Config.Jenkins.PrintStatus {
		// 等待构建
		err = waitDeploy(job)
		if err != nil {
			return err
		}

		// 打印状态
		err = Status(job)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("To check the progress, please visit: %s", getJenkinsUrl(job))
	}

	fmt.Println("🍺 Deploy done!")
	return nil
}

func startBuild(job, server, branch string, checkCode bool) error {
	buildUrl := fmt.Sprintf("%s/job/%s/buildWithParameters", zodo.Config.Jenkins.Url, job)
	requestBody := url.Values{
		"BUILD_BRANCH":  {branch},
		"SERVERNAME":    {server},
		"IS_CHECK_CODE": {strings.ToUpper(boolToText(checkCode))},
	}
	req, err := http.NewRequest("POST", buildUrl, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(zodo.Config.Jenkins.Username, zodo.Config.Jenkins.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("http status not created: %d, resp body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func waitDeploy(job string) error {
	fmt.Println("Wait deploy...")
	for {
		build, err := getLastBuild(job, false)
		if err != nil {
			return err
		}
		for _, stage := range build.Stages {
			if stage.Status == stageStatusInProgress {
				fmt.Println("Start deploy...")
				return nil
			}
		}
		time.Sleep(time.Duration(zodo.Config.Jenkins.PollingIntervalSecond) * time.Second)
	}
}
