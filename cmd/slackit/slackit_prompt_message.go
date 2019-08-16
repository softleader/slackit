package main

import (
	"github.com/manifoldco/promptui"
)

func promptMessage() (string, error) {
	prompt := promptui.Prompt{
		Label: "Message text introducing the file: ",
	}
	return prompt.Run()
}
