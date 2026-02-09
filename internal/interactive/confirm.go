package interactive

import (
	"github.com/AlecAivazis/survey/v2"
)

// ConfirmAction prompts the user for yes/no confirmation
func ConfirmAction(message string) (bool, error) {
	var confirmed bool
	prompt := &survey.Confirm{
		Message: message,
		Default: true,
	}

	if err := survey.AskOne(prompt, &confirmed); err != nil {
		return false, err
	}

	return confirmed, nil
}
