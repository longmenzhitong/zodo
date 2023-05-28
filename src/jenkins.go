package zodo

import (
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	stageStatusSuccess    = "SUCCESS"
	stageStatusInProgress = "IN_PROGRESS"
	stageStatusAborted    = "ABORTED"
)

type JenkinsBuild struct {
	Stages []struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Complete bool   `json:"complete"`
	} `json:"stages"`
}

func Deploy(service, env, branch string, checkCode, statusOnly bool) error {
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
	if statusOnly {
		fmt.Println("Check done.")
		err := printStatus(service)
		return err
	}
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

	// 确认构建
	fmt.Println("Sure to deploy? [y/n]")
	input := strings.ToLower(readString())
	if input != "y" {
		return &CancelledError{}
	}

	// 发起构建
	err := startBuild(service, env, branch, checkCode)
	if err != nil {
		return err
	}

	if Config.Jenkins.PrintStatus {
		// 等待构建
		err = waitDeploy(service)
		if err != nil {
			return err
		}

		// 打印状态
		err = printStatus(service)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("To check the progress, please visit: %s", getJenkinsUrl(service))
	}

	fmt.Println("Deploy done!")
	return nil
}

func startBuild(service, env, branch string, checkCode bool) error {
	buildUrl := fmt.Sprintf("%s/job/%s/buildWithParameters", Config.Jenkins.Url, service)
	requestBody := url.Values{
		"BUILD_BRANCH":  {branch},
		"SERVERNAME":    {env},
		"IS_CHECK_CODE": {strings.ToUpper(boolToText(checkCode))},
	}
	req, err := http.NewRequest("POST", buildUrl, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(Config.Jenkins.Username, Config.Jenkins.Password)
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

func waitDeploy(service string) error {
	fmt.Println("Wait deploy...")
	for {
		build, err := getBuild(service)
		if err != nil {
			return err
		}
		for _, stage := range build.Stages {
			if stage.Status == stageStatusInProgress {
				fmt.Println("Start deploy...")
				return nil
			}
		}
		time.Sleep(time.Duration(Config.Jenkins.PollingIntervalSecond) * time.Second)
	}
}

func printStatus(service string) error {
	if Config.Jenkins.StageCount == 0 {
		return &InvalidConfigError{Message: "jenkins.stageCount doesn't exist"}
	}
	bar := progressbar.Default(int64(Config.Jenkins.StageCount))
	succeedStage := make(map[string]bool, 0)
	for {
		build, err := getBuild(service)
		if err != nil {
			return err
		}
		for _, stage := range build.Stages {
			//fmt.Printf("name: %s, status: %s\n", stage.Name, stage.Status)
			if stage.Status == stageStatusSuccess {
				if succeedStage[stage.Name] {
					continue
				}
				succeedStage[stage.Name] = true
				err = bar.Add(1)
				if err != nil {
					return err
				}
			} else if stage.Status == stageStatusAborted {
				return fmt.Errorf("\ndeploy aborted, please check: %s\n", getJenkinsUrl(service))
			}
		}
		if len(succeedStage) == Config.Jenkins.StageCount {
			break
		}
		time.Sleep(time.Duration(Config.Jenkins.PollingIntervalSecond) * time.Second)
	}
	return nil
}

func getBuild(service string) (*JenkinsBuild, error) {
	statusUrl := fmt.Sprintf("%s/job/%s/lastBuild/wfapi/describe", Config.Jenkins.Url, service)
	req, err := http.NewRequest("GET", statusUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(Config.Jenkins.Username, Config.Jenkins.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok: %d, resp body: %s", resp.StatusCode, resp.Body)
	}
	var build JenkinsBuild
	err = json.NewDecoder(resp.Body).Decode(&build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func getJenkinsUrl(service string) string {
	return fmt.Sprintf("%s/job/%s\n", Config.Jenkins.Url, service)
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
