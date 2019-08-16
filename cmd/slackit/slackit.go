package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slackit/pkg/formatter"
	"github.com/softleader/slackit/pkg/paths"
	"github.com/softleader/slackit/pkg/slackapi"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
)

const (
	usage = `Easily share files from command line to Slack

使用 '--channel' 可以將指定的 FILE 分享至該 Slack 頻道上:

	$ slctl slackit /PATH/TO/FILE /PATH/TO/ANOTHER/FILE -c CHANNEL

不傳入 '--channel' 會產生一個互動式的選擇介面, 包含了你有參與的頻道及所有使用者名稱
配合使用 '--all' 可以改成列出所有, 甚至你沒參與的頻道名稱

	$ slctl slackit FILE...
	$ slctl slackit FILE... --all

頻道及使用者清單會 cache 在本機維持 1 天, 傳入 '--force' 可以強制重新取得清單

	$ slctl slackit FILE... -f

使用 '--message' 可以為上傳檔案註記一段文字

	$ slctl slackit FILE... -m hello

使用 '--rm' 可以在上傳成功後, 自動刪掉本機的檔案

	$ slctl slackit FILE... --rm

第一次使用時, 必須傳入 '--slack-token' 讓 plugin 知道如何與 Slack api 互動
你可以在 https://api.slack.com/custom-integrations/legacy-tokens 建立一個 token
Plugin 在使用後也會自動的記錄下來, 之後就不用再次傳入
相對的, 你也可以再次傳入 '--slack-token' 來 renew 已存在的 token

	$ slctl slackit FILE... --slack-token TOKEN
`
)

var (
	offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	token      = os.Getenv("SL_TOKEN")
	slackToken string
	mount      = os.Getenv("SL_PLUGIN_MOUNT")
	force      bool
	size       = 10
	all        bool
	channel    string
	message    string
	rm         bool
)

func main() {
	if err := newRootCmd(os.Args[1:]).Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slackit",
		Short: "Easily share files from command line to Slack",
		Long:  usage,
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&formatter.PlainFormatter{})
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return send(args)
		},
	}

	f := cmd.PersistentFlags()
	f.BoolVar(&offline, "offline", offline, "work offline, Overrides $SL_OFFLINE")
	f.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVar(&token, "token", token, "github access token. Overrides $SL_TOKEN")
	f.StringVar(&slackToken, "slack-token", slackToken, "slack token")
	f.StringVarP(&channel, "channel", "c", channel, "specify the slack channel, leave blank to open an interactive prompt")
	f.BoolVarP(&force, "force", "f", force, "force to evict the cache")
	f.IntVar(&size, "size", size, "specify the number of items appear on the select prompt")
	f.BoolVar(&all, "all", all, "show all channels instead of related channels")
	f.StringVarP(&message, "message", "m", message, "specify the message text introducing the file")
	f.BoolVar(&rm, "rm", rm, "automatically remove the file after upload succeed")
	f.Parse(args)

	return cmd
}

func isMemberOfChannel(channel slackapi.Channel) bool {
	return all || channel.IsMember
}

func send(files []string) error {
	api, err := slackapi.NewClient(logrus.StandardLogger(), slackToken, mount)
	if err != nil {
		return err
	}
	if channel == "" {
		c, err := slackapi.LoadChannels(logrus.StandardLogger(), api, mount, force)
		if err != nil {
			return err
		}
		selected, err := promptChannel(c.Filter(isMemberOfChannel))
		if err != nil {
			return err
		}
		channel = selected.ID
	}
	for _, file := range files {
		if err := upload(api, file, channel); err != nil {
			return err
		}
	}
	return nil
}

func upload(api *slack.Client, path, channel string) error {
	path, _ = homedir.Expand(path)
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if exist, err := paths.Exists(abs); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("path is not exist: %s", abs)
	}
	if message == "" {
		if message, err = prompt("Message text introducing the file"); err != nil {
			return err
		}
	}
	file := slack.FileUploadParameters{
		File:           abs,
		Filename:       filepath.Base(abs),
		Channels:       []string{channel},
		InitialComment: message,
	}
	if _, err = api.UploadFile(file); err != nil {
		return err
	}
	logrus.Printf("Successfully uploaded file: %s", abs)
	if rm {
		logrus.Debugf("removing file: %s", abs)
		if err := os.Remove(abs); err != nil {
			logrus.Debugln(err)
		}
	}
	return nil
}
