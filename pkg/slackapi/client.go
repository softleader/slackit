package slackapi

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

const (
	tokenFileName = `slack-token`
)

func NewClient(log *logrus.Logger, token, savePath string) (*slack.Client, error) {
	t, err := retrieveToken(log, token, savePath)
	if err != nil {
		return nil, err
	}
	return slack.New(t), nil
}

func retrieveToken(log *logrus.Logger, token, savePath string) (string, error) {
	if token != "" {
		defer saveToken(log, token, savePath)
		return token, nil
	}
	if token, found := loadToken(log, savePath); found {
		return token, nil
	}
	return "", fmt.Errorf("you might need to get a Slack api token from: %s", "https://api.slack.com/custom-integrations/legacy-tokens")
}

func saveToken(log *logrus.Logger, token, path string) {
	if path != "" {
		log.Debug("you have specified a slack token, saving it for the next time")
		ioutil.WriteFile(filepath.Join(path, tokenFileName), []byte(token), 0644)
	}
}

func loadToken(log *logrus.Logger, path string) (string, bool) {
	if path != "" {
		log.Debugf("loading cached token from: %s\n", path)
		b, err := ioutil.ReadFile(filepath.Join(path, tokenFileName));
		if err == nil {
			return string(b), true
		}
		log.Debug(err)
	}
	return "", false
}
