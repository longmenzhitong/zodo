package jenkins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	zodo "zodo/src"
)

func DefaultJob() string {
	return strings.ToUpper(zodo.CurrentDirName())
}

func DefaultBranch() (string, error) {
	b, err := zodo.CurrentGitBranch()
	if err != nil {
		return "", err
	}
	return b, nil
}

type LastBuild struct {
	Stages []struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Complete bool   `json:"complete"`
	} `json:"stages"`
}

func getLastBuild(job string, successful bool) (*LastBuild, error) {
	var statusUrl string
	if successful {
		statusUrl = fmt.Sprintf("%s/job/%s/lastSuccessfulBuild/wfapi/describe", zodo.Config.Jenkins.Url, job)
	} else {
		statusUrl = fmt.Sprintf("%s/job/%s/lastBuild/wfapi/describe", zodo.Config.Jenkins.Url, job)
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
	var build LastBuild
	err = json.NewDecoder(resp.Body).Decode(&build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func getJenkinsUrl(job string) string {
	return fmt.Sprintf("%s/job/%s\n", zodo.Config.Jenkins.Url, job)
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