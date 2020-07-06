package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"libs.altipla.consulting/errors"
)

func Confirm(msg string) (bool, error) {
	var reply bool
	prompt := &survey.Confirm{
		Message: msg,
	}
	if err := survey.AskOne(prompt, &reply); err != nil {
		return false, errors.Trace(err)
	}
	return reply, nil
}
