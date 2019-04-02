# slackit

The [slctl](https://github.com/softleader/slctl) plugin to easily share files from command line to Slack

> 簡單且快速的分享檔案到 Slack 頻道中

## Install

```sh
$ slctl plugin install github.com/softleader/slackit
```

## Usage

使用 `--channel` 可以將指定的 FILE 分享至該 Slack 頻道上:

```sh
$ slctl slackit /PATH/TO/FILE /PATH/TO/ANOTHER/FILE -c CHANNEL
```

不傳入 `--channel` 會產生一個互動式的選擇介面, 包含了所有參與的頻道及使用者名稱

```sh
$ slctl slackit FILE...
```

使用 `--all` 可以列出所有, 甚至你沒參與的頻道及使用者名稱

```sh
$ slctl slackit FILE... --all
```

頻道及使用者清單會 cache 在本機維持 1 天, 傳入 `--force` 可以強制重新取得清單

```sh
$ slctl slackit FILE... -f
```

使用 `--message` 可以為上傳檔案註記一段文字

```sh
$ slctl slackit FILE... -m hello
```

使用 `--rm` 可以在上傳成功後, 自動刪掉本機的檔案

```sh
$ slctl slackit FILE... --rm
```

第一次使用時, 必須傳入 `--slack-token` 讓 plugin 知道如何與 Slack api 互動, 你可以在 https://api.slack.com/custom-integrations/legacy-tokens 建立一個 token, Plugin 在使用後也會自動的記錄下來, 之後就不用再次傳入, 相對的, 你也可以再次傳入 `--slack-token` 來 renew 已存在的 token

```sh
$ slctl slackit FILE... --slack-token TOKEN
```
