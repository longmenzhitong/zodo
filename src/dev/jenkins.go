package dev

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	zodo "zodo/src"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
)

const (
	stageStatusSuccess    = "SUCCESS"
	stageStatusInProgress = "IN_PROGRESS"
	stageStatusAborted    = "ABORTED"
	stageStatusFailed     = "FAILED"
)

type JenkinsBuild struct {
	Stages []struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Complete bool   `json:"complete"`
	} `json:"stages"`
}

func Deploy(service, env, branch string, checkCode, checkStatus bool) error {
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
	if service == "" {
		service = strings.ToUpper(zodo.CurrentDirName())
	} else {
		service = strings.ToUpper(service)
	}
	fmt.Printf("Service   : %s\n", service)
	if checkStatus {
		fmt.Println("Check done.")
		err := printStatus(service)
		return err
	}
	if env == "" {
		fmt.Println("Please input the env:")
		env = zodo.ReadString()
		if env == "" {
			return &zodo.InvalidInputError{Message: "env must not be empty"}
		}
	}
	fmt.Printf("Env       : %s\n", env)
	if branch == "" {
		b, err := zodo.CurrentGitBranch()
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
	input := strings.ToLower(zodo.ReadString())
	if input != "y" {
		return &zodo.CancelledError{}
	}

	// 发起构建
	err := startBuild(service, env, branch, checkCode)
	if err != nil {
		return err
	}

	if zodo.Config.Jenkins.PrintStatus {
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
	buildUrl := fmt.Sprintf("%s/job/%s/buildWithParameters", zodo.Config.Jenkins.Url, service)
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

func waitDeploy(service string) error {
	fmt.Println("Wait deploy...")
	for {
		build, err := getLastBuild(service, false)
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

func printStatus(service string) error {
	stageCount, err := getStageCount(service)
	if err != nil {
		return err
	}
	bar := progressbar.NewOptions(
		stageCount,
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	inProgress := make(map[string]bool, 0)
	succeed := make(map[string]bool, 0)
	for {
		build, err := getLastBuild(service, false)
		if err != nil {
			return err
		}
		for _, stage := range build.Stages {
			name := stage.Name
			switch stage.Status {
			case stageStatusSuccess:
				if succeed[name] {
					continue
				}
				succeed[name] = true
				err = bar.Add(1)
				if err != nil {
					return err
				}
			case stageStatusInProgress:
				if inProgress[name] {
					continue
				}
				inProgress[name] = true
				bar.Describe(name)
			case stageStatusAborted:
				return fmt.Errorf("\ndeploy aborted, please check: %s\n", getJenkinsUrl(service))
			case stageStatusFailed:
				return fmt.Errorf("\ndeploy failed, please check: %s\n", getJenkinsUrl(service))
			default:
				panic(fmt.Errorf("unexpected stage status: %s", stage.Status))
			}
		}
		if len(succeed) == stageCount {
			break
		}
		time.Sleep(time.Duration(zodo.Config.Jenkins.PollingIntervalSecond) * time.Second)
	}
	return nil
}

func getStageCount(service string) (int, error) {
	build, err := getLastBuild(service, true)
	if err != nil {
		return -1, err
	}
	return len(build.Stages), nil
}

func getLastBuild(service string, successful bool) (*JenkinsBuild, error) {
	var statusUrl string
	if successful {
		statusUrl = fmt.Sprintf("%s/job/%s/lastSuccessfulBuild/wfapi/describe", zodo.Config.Jenkins.Url, service)
	} else {
		statusUrl = fmt.Sprintf("%s/job/%s/lastBuild/wfapi/describe", zodo.Config.Jenkins.Url, service)
	}
	req, err := http.NewRequest("GET", statusUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(zodo.Config.Jenkins.Username, zodo.Config.Jenkins.Password)
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
	return fmt.Sprintf("%s/job/%s\n", zodo.Config.Jenkins.Url, service)
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

type Job struct {
	Builds []Build `xml:"build"`
}

type Build struct {
	XMLName xml.Name `xml:"build"`
	Class   string   `xml:"_class,attr"`
	Number  int      `xml:"number"`
	URL     string   `xml:"url"`
}

type BuildInfo struct {
	XMLName           xml.Name      `xml:"workflowRun"`
	Class             string        `xml:"_class,attr"`
	Actions           []BuildAction `xml:"action"`
	Building          bool          `xml:"building"`
	Description       string        `xml:"description"`
	DisplayName       string        `xml:"displayName"`
	Duration          int           `xml:"duration"`
	EstimatedDuration int           `xml:"estimatedDuration"`
	FullDisplayName   string        `xml:"fullDisplayName"`
	ID                int           `xml:"id"`
	KeepLog           bool          `xml:"keepLog"`
	Number            int           `xml:"number"`
	QueueID           int           `xml:"queueId"`
	Result            string        `xml:"result"`
	Timestamp         int64         `xml:"timestamp"`
	URL               string        `xml:"url"`
}

type BuildAction struct {
	Class      string           `xml:"_class,attr"`
	Parameters []BuildParameter `xml:"parameter,omitempty"`
	Cause      BuildCause       `xml:"cause,omitempty"`
}

type BuildParameter struct {
	Class string `xml:"_class,attr"`
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type BuildCause struct {
	Class            string `xml:"_class,attr"`
	ShortDescription string `xml:"shortDescription"`
	UserID           string `xml:"userId"`
	UserName         string `xml:"userName"`
}

// TODO 把这几个name定义成常量，把这几个方法合并在一起
func (b *BuildInfo) getBuildBranch() string {
	for _, a := range b.Actions {
		if a.Class == "hudson.model.ParametersAction" {
			for _, p := range a.Parameters {
				if p.Name == "BUILD_BRANCH" {
					return p.Value
				}
			}
		}
	}
	return ""
}

func (b *BuildInfo) getDeployEnv() string {
	for _, a := range b.Actions {
		if a.Class == "hudson.model.ParametersAction" {
			for _, p := range a.Parameters {
				if p.Name == "DEPLOYENV" {
					return p.Value
				}
			}
		}
	}
	return ""
}

func (b *BuildInfo) getServerName() string {
	for _, a := range b.Actions {
		if a.Class == "hudson.model.ParametersAction" {
			for _, p := range a.Parameters {
				if p.Name == "SERVERNAME" {
					return p.Value
				}
			}
		}
	}
	return ""
}

func (b *BuildInfo) getBuildUser() string {
	for _, a := range b.Actions {
		if a.Class == "hudson.model.CauseAction" {
			return a.Cause.UserName
		}
	}
	return ""
}

func History(job string) error {
	allBuildUrl := fmt.Sprintf("%s/job/%s/api/xml", zodo.Config.Jenkins.Url, job)
	req, err := http.NewRequest("GET", allBuildUrl, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(zodo.Config.Jenkins.Username, zodo.Config.Jenkins.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status not ok: %d, resp body: %s", resp.StatusCode, resp.Body)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var workflowJob Job
	err = xml.Unmarshal(b, &workflowJob)
	if err != nil {
		return err
	}

	limit := 5
	buildUrls := make([]string, 0)
	for _, b := range workflowJob.Builds {
		buildUrls = append(buildUrls, b.URL+"api/xml")
	}
	buildUrls = buildUrls[:limit]

	rows := make([]table.Row, 0)

	for _, url := range buildUrls {
		buildInfo, err := getBuildInfo(url)
		if err != nil {
			return err
		}

		row := table.Row{
			time.Unix(buildInfo.Timestamp/1000, 0).Format(zodo.LayoutDateTime),
			buildInfo.getDeployEnv(),
			buildInfo.getServerName(),
			buildInfo.getBuildBranch(),
			buildInfo.Result,
			buildInfo.Building,
			buildInfo.getBuildUser(),
		}
		rows = append(rows, row)
	}

	title := table.Row{"时间", "环境", "服务器", "分支", "结果", "正在构建", "发起用户"}
	zodo.PrintTable(title, rows)
	return nil
}

func getBuildInfo(url string) (*BuildInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(zodo.Config.Jenkins.Username, zodo.Config.Jenkins.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok: %d, resp body: %s", resp.StatusCode, resp.Body)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var buildInfo BuildInfo
	err = xml.Unmarshal(b, &buildInfo)
	if err != nil {
		return nil, err
	}

	return &buildInfo, nil
}
