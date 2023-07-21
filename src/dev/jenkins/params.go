package jenkins

import (
	"fmt"
	"os"
	"strings"
	zodo "zodo/src"

	"gopkg.in/yaml.v3"
)

const paramFileName = ".jenkins_param.yml"
const templateValPrefix = "$"

type param struct {
	Job         string            `yaml:"job"`
	BuildParams map[string]string `yaml:"buildParams"`
}

func (p *param) checkJob() error {
	if p.Job == "" {
		fmt.Println("Please input job:")
		p.Job = zodo.ReadString()
		if p.Job == "" {
			return &zodo.InvalidInputError{Message: "job must not be empty"}
		}
	} else if strings.HasPrefix(p.Job, templateValPrefix) {
		job, err := getTemplateValue(p.Job)
		if err != nil {
			return err
		}
		p.Job = job
	}

	return nil
}

func (p *param) checkBuildParams() error {
	if len(p.BuildParams) == 0 {
		return &zodo.InvalidInputError{Message: "build params must not be empty"}
	}

	for k, v := range p.BuildParams {
		var err error
		if v == "" {
			v, err = askParamValue(k)
		} else if strings.HasPrefix(v, templateValPrefix) {
			v, err = getTemplateValue(v)
		}
		if err != nil {
			return err
		}

		p.BuildParams[k] = v
	}

	return nil
}

func GetParam(build bool) (*param, error) {
	var p *param
	var err error

	p, err = getParamFromPath(zodo.CurrentPath(paramFileName))
	if err != nil {
		p, err = getParamFromPath(zodo.Path(paramFileName))
		if err != nil {
			return nil, err
		}
	}

	err = p.checkJob()
	if err != nil {
		return nil, err
	}

	if build {
		err := p.checkBuildParams()
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func getParamFromPath(path string) (*param, error) {
	var p param

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func askParamValue(paramKey string) (string, error) {
	fmt.Printf("Please input %s:\n", paramKey)
	v := zodo.ReadString()
	if v == "" {
		return v, &zodo.InvalidInputError{Message: fmt.Sprintf("%s must not be empty", paramKey)}
	}
	return v, nil
}

func getTemplateValue(templateKey string) (string, error) {
	var v string
	if strings.HasPrefix(templateKey, "$CURRENT_DIR_NAME") {
		v = zodo.CurrentDirName()
	} else if strings.HasPrefix(templateKey, "$CURRENT_GIT_BRANCH") {
		b, err := zodo.CurrentGitBranch()
		if err != nil {
			return "", err
		}
		v = b
	}

	if strings.HasSuffix(templateKey, ".UPPER") {
		v = strings.ToUpper(v)
	}

	return v, nil
}
