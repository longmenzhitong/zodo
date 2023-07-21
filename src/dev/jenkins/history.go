package jenkins

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
	zodo "zodo/src"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Job struct {
	Builds []JobBuild `xml:"build"`
}

type JobBuild struct {
	XMLName xml.Name `xml:"build"`
	Class   string   `xml:"_class,attr"`
	Number  int      `xml:"number"`
	URL     string   `xml:"url"`
}

type Build struct {
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
	Class      string                 `xml:"_class,attr"`
	Parameters []BuildActionParameter `xml:"parameter,omitempty"`
	Cause      BuildActionCause       `xml:"cause,omitempty"`
}

type BuildActionParameter struct {
	Class string `xml:"_class,attr"`
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type BuildActionCause struct {
	Class            string `xml:"_class,attr"`
	ShortDescription string `xml:"shortDescription"`
	UserID           string `xml:"userId"`
	UserName         string `xml:"userName"`
}

const (
	buildActionClassParameter           = "hudson.model.ParametersAction"
	buildActionClassCuase               = "hudson.model.CauseAction"
	buildActionParameterNameBuildBranch = "BUILD_BRANCH"
	buildActionParameterNameDeployEnv   = "DEPLOYENV"
	buildActionParameterNameServerName  = "SERVERNAME"
)

func (b *Build) getPrameterValue(parameterName string) string {
	for _, a := range b.Actions {
		if a.Class == buildActionClassParameter {
			for _, p := range a.Parameters {
				if p.Name == parameterName {
					return p.Value
				}
			}
		}
	}
	return ""
}

func (b *Build) getCauseUser() string {
	for _, a := range b.Actions {
		if a.Class == buildActionClassCuase {
			return a.Cause.UserName
		}
	}
	return ""
}

func (b *Build) getResult() string {
	var result string
	if b.Building {
		result = "❓"
	} else {
		result = boolToSymbol(b.Result == "SUCCESS")
	}
	return result
}

func History(count int) error {
	p, err := GetParam(false)
	if err != nil {
		return err
	}

	// 获取Job信息
	jobUrl := fmt.Sprintf("%s/job/%s/api/xml", zodo.Config.Jenkins.Url, p.Job)
	req, err := http.NewRequest("GET", jobUrl, nil)
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
	var job Job
	err = xml.Unmarshal(b, &job)
	if err != nil {
		return err
	}

	// 从Job信息中提取Build信息的URL
	buildUrls := make([]string, 0)
	for i := 0; i < len(job.Builds) && i < count; i++ {
		buildUrl := job.Builds[i].URL + "api/xml"
		buildUrls = append(buildUrls, buildUrl)
	}

	// 访问Build信息URL获取Build信息
	rows := make([]table.Row, 0)
	for _, buildUrl := range buildUrls {
		build, err := getBuild(buildUrl)
		if err != nil {
			return err
		}

		row := table.Row{
			zodo.SimplifyTime(time.Unix(build.Timestamp/1000, 0).Format(zodo.LayoutDateTime)),
			build.getPrameterValue(buildActionParameterNameDeployEnv),
			build.getPrameterValue(buildActionParameterNameServerName),
			build.getPrameterValue(buildActionParameterNameBuildBranch),
			build.getResult(),
			build.getCauseUser(),
		}
		rows = append(rows, row)
	}

	title := table.Row{"Time", "Env", "Server", "Branch", "Result", "User"}
	zodo.PrintTable(title, rows)
	return nil
}

func getBuild(buildUrl string) (*Build, error) {
	req, err := http.NewRequest("GET", buildUrl, nil)
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

	var build Build
	err = xml.Unmarshal(b, &build)
	if err != nil {
		return nil, err
	}

	return &build, nil
}
