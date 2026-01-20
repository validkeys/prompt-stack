package prompt

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type Question struct {
	ID       string
	Text     string
	Required bool
	Validate func(string) error
}

type Prompt struct {
	questions []Question
	responses map[string]string
}

type InterviewResult struct {
	Responses  map[string]string
	Transcript string
}

func NewPrompt(questions []Question) *Prompt {
	return &Prompt{
		questions: questions,
		responses: make(map[string]string),
	}
}

var readStringFunc = func(reader *bufio.Reader, delim byte) (string, error) {
	return reader.ReadString(delim)
}

func (p *Prompt) Run(ctx context.Context) (*InterviewResult, error) {
	reader := bufio.NewReader(os.Stdin)
	var transcript strings.Builder

	for _, q := range p.questions {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		transcript.WriteString(fmt.Sprintf("Q: %s\n\n", q.Text))

		var response string
		var err error

		for {
			fmt.Printf("%s\n", q.Text)
			if q.Required {
				fmt.Printf("(Required) Your answer: ")
			} else {
				fmt.Printf("(Optional, press Enter to skip) Your answer: ")
			}

			response, err = readStringFunc(reader, '\n')
			if err != nil {
				return nil, fmt.Errorf("failed to read input: %w", err)
			}

			response = strings.TrimSpace(response)

			if response == "" && !q.Required {
				break
			}

			if q.Validate != nil {
				if err := q.Validate(response); err != nil {
					fmt.Printf("Validation error: %v\n", err)
					continue
				}
			}

			if response == "" && q.Required {
				fmt.Println("This field is required. Please provide an answer.")
				continue
			}

			break
		}

		p.responses[q.ID] = response
		transcript.WriteString(fmt.Sprintf("A: %s\n\n", response))
	}

	return &InterviewResult{
		Responses:  p.responses,
		Transcript: transcript.String(),
	}, nil
}
