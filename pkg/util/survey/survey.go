package survey

import (
	"os"
	"regexp"

	"github.com/pkg/errors"
	surveypkg "gopkg.in/AlecAivazis/survey.v1"
)

// QuestionOptions defines a question and its options
type QuestionOptions struct {
	Question               string
	DefaultValue           string
	ValidationRegexPattern string
	ValidationMessage      string
	ValidationFunc         func(value string) error
	Options                []string
	IsPassword             bool
}

// DefaultValidationRegexPattern is the default regex pattern to validate the input
var DefaultValidationRegexPattern = regexp.MustCompile("^.*$")

// Survey is the interface for asking questions
type Survey interface {
	Question(params *QuestionOptions) (string, error)
}

type survey struct{}

// NewSurvey creates a new survey object
func NewSurvey() Survey {
	return &survey{}
}

// Question asks the user a question and returns the answer
func (s *survey) Question(params *QuestionOptions) (string, error) {
	var prompt surveypkg.Prompt
	compiledRegex := DefaultValidationRegexPattern
	if params.ValidationRegexPattern != "" {
		compiledRegex = regexp.MustCompile(params.ValidationRegexPattern)
	}

	if params.Options != nil {
		prompt = &surveypkg.Select{
			Message: params.Question + "\n",
			Options: params.Options,
			Default: params.DefaultValue,
		}
	} else if params.IsPassword {
		prompt = &surveypkg.Password{
			Message: params.Question,
		}
	} else {
		prompt = &surveypkg.Input{
			Message: params.Question,
			Default: params.DefaultValue,
		}
	}

	question := []*surveypkg.Question{
		{
			Name:   "question",
			Prompt: prompt,
		},
	}

	if params.Options == nil {
		question[0].Validate = func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return errors.New("Input was not a string")
			}

			// Check regex
			if compiledRegex.MatchString(str) == false {
				if params.ValidationMessage != "" {
					return errors.New(params.ValidationMessage)
				}

				return errors.Errorf("Answer has to match pattern: %s", compiledRegex.String())
			}

			// Check function
			if params.ValidationFunc != nil {
				err := params.ValidationFunc(str)
				if err != nil {
					if params.ValidationMessage != "" {
						return errors.New(params.ValidationMessage)
					}

					return errors.Errorf("%v", err)
				}
			}

			return nil
		}
	}

	// Ask it
	answers := struct {
		Question string
	}{}

	err := surveypkg.Ask(question, &answers)
	if err != nil {
		// Keyboard interrupt
		os.Exit(0)
	}

	return answers.Question, nil
}
