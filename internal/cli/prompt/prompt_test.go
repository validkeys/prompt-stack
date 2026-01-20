package prompt

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestNewPrompt(t *testing.T) {
	questions := []Question{
		{ID: "q1", Text: "Question 1", Required: true},
	}
	p := NewPrompt(questions)

	if p == nil {
		t.Fatal("NewPrompt returned nil")
	}
	if len(p.questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(p.questions))
	}
	if p.responses == nil {
		t.Error("responses map not initialized")
	}
}

func TestPromptRun_EmptyQuestions(t *testing.T) {
	p := NewPrompt([]Question{})
	ctx := context.Background()

	result, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(result.Responses) != 0 {
		t.Errorf("Expected 0 responses, got %d", len(result.Responses))
	}
}

func TestPromptRun_SingleQuestion(t *testing.T) {
	p := NewPrompt([]Question{
		{ID: "q1", Text: "What is your name?", Required: true},
	})

	oldReadStringFunc := readStringFunc
	defer func() { readStringFunc = oldReadStringFunc }()

	inputs := []string{"John Doe\n"}
	inputIndex := 0

	readStringFunc = func(reader *bufio.Reader, delim byte) (string, error) {
		if inputIndex >= len(inputs) {
			return "\n", nil
		}
		result := inputs[inputIndex]
		inputIndex++
		return result, nil
	}

	ctx := context.Background()
	result, err := p.Run(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Responses["q1"] != "John Doe" {
		t.Errorf("Expected response 'John Doe', got '%s'", result.Responses["q1"])
	}
	if !strings.Contains(result.Transcript, "Q: What is your name?") {
		t.Error("Transcript missing question")
	}
	if !strings.Contains(result.Transcript, "A: John Doe") {
		t.Error("Transcript missing answer")
	}
}

func TestPromptRun_Validation(t *testing.T) {
	tests := []struct {
		name        string
		question    Question
		inputs      []string
		expectError bool
		expected    string
	}{
		{
			name: "valid email",
			question: Question{
				ID:   "email",
				Text: "Enter your email",
				Validate: func(s string) error {
					if !strings.Contains(s, "@") {
						return fmt.Errorf("invalid email")
					}
					return nil
				},
				Required: true,
			},
			inputs:      []string{"invalid", "test@example.com\n"},
			expectError: false,
			expected:    "test@example.com",
		},
		{
			name: "min length",
			question: Question{
				ID:   "name",
				Text: "Enter your name",
				Validate: func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("name too short")
					}
					return nil
				},
				Required: true,
			},
			inputs:      []string{"ab", "Alice\n"},
			expectError: false,
			expected:    "Alice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPrompt([]Question{tt.question})

			oldReadStringFunc := readStringFunc
			defer func() { readStringFunc = oldReadStringFunc }()

			inputs := tt.inputs
			inputIndex := 0

			readStringFunc = func(reader *bufio.Reader, delim byte) (string, error) {
				if inputIndex >= len(inputs) {
					return "\n", nil
				}
				result := inputs[inputIndex]
				inputIndex++
				return result, nil
			}

			ctx := context.Background()
			result, err := p.Run(ctx)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && result.Responses[tt.question.ID] != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result.Responses[tt.question.ID])
			}
		})
	}
}

func TestPromptRun_OptionalQuestion(t *testing.T) {
	p := NewPrompt([]Question{
		{ID: "optional", Text: "Optional question", Required: false},
	})

	oldReadStringFunc := readStringFunc
	defer func() { readStringFunc = oldReadStringFunc }()

	readStringFunc = func(reader *bufio.Reader, delim byte) (string, error) {
		return "\n", nil
	}

	ctx := context.Background()
	result, err := p.Run(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Responses["optional"] != "" {
		t.Errorf("Expected empty response for optional question, got '%s'", result.Responses["optional"])
	}
}

func TestPromptRun_RequiredQuestionEmpty(t *testing.T) {
	p := NewPrompt([]Question{
		{ID: "required", Text: "Required question", Required: true},
	})

	oldReadStringFunc := readStringFunc
	defer func() { readStringFunc = oldReadStringFunc }()

	inputs := []string{"", "Answer\n"}
	callCount := 0

	readStringFunc = func(reader *bufio.Reader, delim byte) (string, error) {
		if callCount < len(inputs) {
			result := inputs[callCount]
			callCount++
			return result, nil
		}
		return "\n", nil
	}

	ctx := context.Background()
	result, err := p.Run(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if callCount < 2 {
		t.Error("Expected to retry on empty required field")
	}
	if result.Responses["required"] != "Answer" {
		t.Errorf("Expected 'Answer', got '%s'", result.Responses["required"])
	}
}

func TestPromptRun_ContextCancellation(t *testing.T) {
	p := NewPrompt([]Question{
		{ID: "q1", Text: "Question 1", Required: true},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.Run(ctx)
	if err == nil {
		t.Error("Expected error from cancelled context")
	}
}

func TestPromptRun_TranscriptFormat(t *testing.T) {
	p := NewPrompt([]Question{
		{ID: "q1", Text: "First question?", Required: true},
		{ID: "q2", Text: "Second question?", Required: false},
	})

	oldReadStringFunc := readStringFunc
	defer func() { readStringFunc = oldReadStringFunc }()

	inputs := []string{"Answer 1\n", "Answer 2\n"}
	inputIndex := 0

	readStringFunc = func(reader *bufio.Reader, delim byte) (string, error) {
		if inputIndex >= len(inputs) {
			return "\n", nil
		}
		result := inputs[inputIndex]
		inputIndex++
		return result, nil
	}

	ctx := context.Background()
	result, err := p.Run(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	lines := strings.Split(result.Transcript, "\n")
	if len(lines) < 6 {
		t.Errorf("Expected at least 6 lines in transcript, got %d", len(lines))
	}

	if !strings.Contains(result.Transcript, "Q: First question?") {
		t.Error("Transcript missing first question")
	}
	if !strings.Contains(result.Transcript, "A: Answer 1") {
		t.Error("Transcript missing first answer")
	}
	if !strings.Contains(result.Transcript, "Q: Second question?") {
		t.Error("Transcript missing second question")
	}
}
