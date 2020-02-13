package cmd

import "github.com/AlecAivazis/survey/v2"

var containerQs = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "What's the name of the container?"},
		Validate: survey.Required,
	},
	{
		Name:     "image",
		Prompt:   &survey.Input{Message: "What's the name of the image?"},
		Validate: survey.Required,
	},
}
