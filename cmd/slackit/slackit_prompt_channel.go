package main

import (
	"github.com/manifoldco/promptui"
	"github.com/softleader/slackit/pkg/slackapi"
	"strings"
)

func promptChannel(channels []slackapi.Channel) (*slackapi.Channel, error) {
	prompt := promptui.Select{
		Label: "Select Channel",
		Items: channels,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   promptui.IconSelect + " {{ .Name }}",
			Inactive: "  {{ .Name }}",
			Selected: promptui.IconGood + " {{ .Name }}",
		},
		Searcher: func(input string, index int) bool {
			channel := channels[index]
			name := strings.Replace(strings.ToLower(channel.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
		Size: size,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &channels[i], nil
}
