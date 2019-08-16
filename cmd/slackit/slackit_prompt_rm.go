package main

import (
	"github.com/manifoldco/promptui"
)

func promptRm() (bool, error) {
	prompt := promptui.Prompt{
		Label: "Message text introducing the file: ",
	}
	return prompt.Run()
}
