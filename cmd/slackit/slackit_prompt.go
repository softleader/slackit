package main

import (
	"github.com/manifoldco/promptui"
)

func prompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	return prompt.Run()
}
