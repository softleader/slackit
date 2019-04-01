package slackapi

import (
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	channelsFileName = `slack-channels.yaml`
)

type Channels struct {
	Channels []Channel
	Expires  time.Time
}

func (c *Channels) expired() bool {
	return c == nil || c.Expires.Before(time.Now())
}

func (c *Channels) save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func (c *Channels) Filter(filter func(Channel) bool) []Channel {
	if filter == nil {
		return c.Channels
	}
	var filtered []Channel
	for _, channel := range c.Channels {
		if filter(channel) {
			filtered = append(filtered, channel)
		}
	}
	return filtered
}

type Channel struct {
	Name     string
	ID       string
	IsMember bool
}

func LoadChannels(log *logrus.Logger, api *slack.Client, path string, force bool) (*Channels, error) {
	cachePath := filepath.Join(path, channelsFileName)
	if force {
		c, err := fetchOnline(log, api)
		if err != nil {
			return nil, err
		}
		c.save(cachePath)
		return c, nil
	}
	c, err := loadLocal(log, cachePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if c.expired() {
		log.Debugln("cache is out of date")
		c, err := fetchOnline(log, api)
		if err != nil {
			return nil, err
		}
		c.save(cachePath)
	}
	return c, nil
}

func loadLocal(log *logrus.Logger, path string) (c *Channels, err error) {
	log.Debugf("loading cached slack channels from: %s\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	c = &Channels{}
	err = yaml.Unmarshal(data, c)
	return
}

func fetchOnline(log *logrus.Logger, api *slack.Client) (*Channels, error) {
	log.Debugf("fetching slack channels")
	c := &Channels{
		Expires: time.Now().AddDate(0, 0, 1),
	}
	channels, err := api.GetChannels(true)
	if err != nil {
		return nil, err
	}
	for _, channel := range channels {
		c.Channels = append(c.Channels, Channel{
			Name:     channel.Name,
			ID:       channel.ID,
			IsMember: channel.IsMember,
		})
	}
	users, err := api.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if !user.IsBot && !user.Deleted && user.ID != "USLACKBOT" {
			c.Channels = append(c.Channels, Channel{
				Name:     user.Name,
				ID:       user.ID,
				IsMember: true,
			})
		}
	}
	return c, nil
}
