package prompt

import (
	"fmt"
	"strings"
)

func DefaultQuestions() []Question {
	return []Question{
		{
			ID:       "milestone_id",
			Text:     "What is the milestone id/slug you want to use for this work? (e.g. m1, m0, readme-update)",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("milestone ID cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "milestone_title",
			Text:     "Give me a one-line title for the milestone.",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("milestone title cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "milestone_description",
			Text:     "Provide a short (1-2 sentence) description of the milestone goal.",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("milestone description cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "stakeholder",
			Text:     "Who is the primary stakeholder or product owner for this milestone? Provide name and contact (email or handle).",
			Required: true,
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("stakeholder information cannot be empty")
				}
				return nil
			},
		},
		{
			ID:       "objectives",
			Text:     "What are the top objectives for this milestone? (3 max) Please list them concisely, one per line.",
			Required: true,
			Validate: func(s string) error {
				lines := strings.Split(strings.TrimSpace(s), "\n")
				validLines := 0
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						validLines++
					}
				}
				if validLines == 0 {
					return fmt.Errorf("at least one objective is required")
				}
				if validLines > 3 {
					return fmt.Errorf("please limit to 3 objectives maximum")
				}
				return nil
			},
		},
		{
			ID:       "success_metrics",
			Text:     "How will we measure success? Give metrics and targets (e.g., quality score, delivery dates, performance targets).",
			Required: true,
		},
		{
			ID:       "style_anchors",
			Text:     "Which files or code areas should be used as style anchors? (list 1–3 paths or URLs, one per line).",
			Required: true,
		},
		{
			ID:       "scope",
			Text:     "What should be included in scope for this milestone? (brief list)",
			Required: true,
		},
		{
			ID:       "out_of_scope",
			Text:     "What is explicitly out of scope for this milestone? (brief list)",
			Required: true,
		},
		{
			ID:       "constraints",
			Text:     "Any critical constraints or assumptions? (security, infra, timelines) — list any you need us to assume.",
			Required: true,
		},
		{
			ID:       "deliverables",
			Text:     "What deliverables do you expect? (e.g., task_breakdown.yaml, quality_report.json, README). List the main artifacts.",
			Required: true,
		},
		{
			ID:       "timeline",
			Text:     "Do you have desired timelines or dates? (start / target completion) — provide if you have them; otherwise say 'assumption: start now, complete in 1-3 days.'",
			Required: true,
		},
		{
			ID:       "testing",
			Text:     "Testing requirements: unit tests, integration tests, TDD preference?",
			Required: true,
		},
		{
			ID:       "privacy",
			Text:     "Any privacy, compliance, or secrets handling notes? (short)",
			Required: true,
		},
	}
}
