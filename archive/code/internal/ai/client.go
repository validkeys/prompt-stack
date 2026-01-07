package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"go.uber.org/zap"
)

// Client wraps Anthropic SDK client with additional functionality
type Client struct {
	client      anthropic.Client
	apiKey      string
	model       string
	logger      *zap.Logger
	httpClient  *http.Client
	maxRetries  int
	baseTimeout time.Duration
}

// Config holds configuration for the AI client
type Config struct {
	APIKey     string
	Model      string
	MaxRetries int
	Timeout    time.Duration
	Logger     *zap.Logger
}

// NewClient creates a new Claude API client wrapper
func NewClient(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if cfg.Model == "" {
		cfg.Model = "claude-3-sonnet-20240229" // Default model
	}

	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3 // Default retry count
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second // Default timeout
	}

	// Create custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	// Create Anthropic client with custom HTTP client
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithHTTPClient(httpClient),
	)

	return &Client{
		client:      client,
		apiKey:      cfg.APIKey,
		model:       cfg.Model,
		logger:      cfg.Logger,
		httpClient:  httpClient,
		maxRetries:  cfg.MaxRetries,
		baseTimeout: cfg.Timeout,
	}, nil
}

// MessageRequest represents a request to send a message to Claude
type MessageRequest struct {
	SystemPrompt string
	Messages     []Message
	MaxTokens    int
	Temperature  float64
}

// Message represents a single message in conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MessageResponse represents response from Claude
type MessageResponse struct {
	Content    string
	StopReason string
	Usage      Usage
	Model      string
	StatusCode int
	RetryAfter *time.Duration // For rate limiting
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// SendMessage sends a message to Claude with retry logic
func (c *Client) SendMessage(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
	var lastErr error
	var retryCount int

	for retryCount = 0; retryCount <= c.maxRetries; retryCount++ {
		if retryCount > 0 {
			c.logger.Info("Retrying API request",
				zap.Int("attempt", retryCount),
				zap.Int("max_retries", c.maxRetries),
				zap.Error(lastErr),
			)

			// Exponential backoff: 1s, 2s, 4s, 8s...
			backoffDuration := time.Duration(1<<uint(retryCount-1)) * time.Second
			if backoffDuration > 30*time.Second {
				backoffDuration = 30 * time.Second // Cap at 30 seconds
			}

			select {
			case <-time.After(backoffDuration):
				// Continue with retry
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Log the request
		c.logger.Debug("Sending API request",
			zap.String("model", c.model),
			zap.Int("message_count", len(req.Messages)),
			zap.Int("max_tokens", req.MaxTokens),
			zap.Float64("temperature", req.Temperature),
		)

		// Convert internal Message format to SDK format
		messages := make([]anthropic.MessageParam, len(req.Messages))
		for i, msg := range req.Messages {
			if msg.Role == "user" {
				messages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
			} else if msg.Role == "assistant" {
				messages[i] = anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content))
			}
		}

		// Create message request
		messageReq := anthropic.MessageNewParams{
			Model:     anthropic.Model(c.model),
			MaxTokens: int64(req.MaxTokens),
			Messages:  messages,
		}

		if req.SystemPrompt != "" {
			messageReq.System = []anthropic.TextBlockParam{{Text: req.SystemPrompt, Type: "text"}}
		}

		if req.Temperature > 0 {
			messageReq.Temperature = anthropic.Float(req.Temperature)
		}

		// Send the request
		resp, err := c.client.Messages.New(ctx, messageReq)
		if err != nil {
			// Check if this is a retryable error
			if isRetryableError(err) {
				lastErr = err
				c.logger.Warn("API request failed, will retry",
					zap.Int("attempt", retryCount+1),
					zap.Error(err),
				)
				continue
			}

			// Non-retryable error
			c.logger.Error("API request failed (non-retryable)",
				zap.Error(err),
			)
			return nil, fmt.Errorf("API request failed: %w", err)
		}

		// Extract content from response
		content := ""
		for _, block := range resp.Content {
			if block.Type == "text" {
				content += block.Text
			}
		}

		// Log the response
		c.logger.Debug("Received API response",
			zap.String("model", string(resp.Model)),
			zap.String("stop_reason", string(resp.StopReason)),
			zap.Int("input_tokens", int(resp.Usage.InputTokens)),
			zap.Int("output_tokens", int(resp.Usage.OutputTokens)),
			zap.Int("content_length", len(content)),
		)

		// Return successful response
		return &MessageResponse{
			Content:    content,
			StopReason: string(resp.StopReason),
			Usage: Usage{
				InputTokens:  int(resp.Usage.InputTokens),
				OutputTokens: int(resp.Usage.OutputTokens),
			},
			Model:      string(resp.Model),
			StatusCode: 200,
		}, nil
	}

	// All retries exhausted
	c.logger.Error("API request failed after all retries",
		zap.Int("total_attempts", retryCount),
		zap.Error(lastErr),
	)
	return nil, fmt.Errorf("API request failed after %d retries: %w", c.maxRetries, lastErr)
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for rate limit errors (429)
	if strings.Contains(err.Error(), "429") || strings.Contains(strings.ToLower(err.Error()), "rate limit") {
		return true
	}

	// Check for timeout errors
	if strings.Contains(strings.ToLower(err.Error()), "timeout") || strings.Contains(strings.ToLower(err.Error()), "deadline exceeded") {
		return true
	}

	// Check for network errors
	if strings.Contains(strings.ToLower(err.Error()), "connection refused") ||
		strings.Contains(strings.ToLower(err.Error()), "connection reset") ||
		strings.Contains(strings.ToLower(err.Error()), "temporary failure") {
		return true
	}

	// Check for 5xx server errors
	if strings.Contains(err.Error(), "500") || strings.Contains(err.Error(), "502") || strings.Contains(err.Error(), "503") {
		return true
	}

	return false
}

// IsAuthError checks if an error is an authentication error (401)
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "401") || strings.Contains(strings.ToLower(err.Error()), "unauthorized") || strings.Contains(strings.ToLower(err.Error()), "authentication")
}

// IsRateLimitError checks if an error is a rate limit error (429)
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "429") || strings.Contains(strings.ToLower(err.Error()), "rate limit")
}

// GetRetryAfter extracts the retry-after duration from a rate limit error
func GetRetryAfter(err error) *time.Duration {
	// This is a placeholder - in a real implementation, you would parse
	// the Retry-After header from the HTTP response
	// For now, we return a default of 30 seconds
	if IsRateLimitError(err) {
		duration := 30 * time.Second
		return &duration
	}
	return nil
}

// EstimateTokens estimates the number of tokens in a text string
// This is a rough approximation: ~4 characters ≈ 1 token
func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}
	// Rough approximation: 4 characters per token
	return len(text) / 4
}

// GetModelContextLimit returns the context limit for the current model
func (c *Client) GetModelContextLimit() int {
	// Context limits for different Claude models
	limits := map[string]int{
		"claude-3-opus-20240229":     200000,
		"claude-3-sonnet-20240229":   200000,
		"claude-3-haiku-20240307":    200000,
		"claude-3-5-sonnet-20241022": 200000,
		"claude-3-5-sonnet-20240620": 200000,
		"claude-3-5-haiku-20241022":  200000,
	}

	if limit, ok := limits[c.model]; ok {
		return limit
	}

	// Default to 200K for unknown models
	return 200000
}

// SetModel updates the model used by the client
func (c *Client) SetModel(model string) {
	c.model = model
	c.logger.Info("Model updated", zap.String("model", model))
}

// GetModel returns the current model
func (c *Client) GetModel() string {
	return c.model
}

// SetTimeout updates the timeout for HTTP requests
func (c *Client) SetTimeout(timeout time.Duration) {
	c.baseTimeout = timeout
	c.httpClient.Timeout = timeout
	c.logger.Info("Timeout updated", zap.Duration("timeout", timeout))
}

// Close cleans up resources used by the client
func (c *Client) Close() error {
	// The Anthropic SDK client doesn't require explicit cleanup
	// This method is provided for future extensibility
	return nil
}

// Helper function to log HTTP requests/responses
func (c *Client) logHTTPRequest(method, url string, body interface{}) {
	if c.logger == nil {
		return
	}

	bodyStr := ""
	if body != nil {
		if jsonBytes, err := json.Marshal(body); err == nil {
			bodyStr = string(jsonBytes)
		}
	}

	c.logger.Debug("HTTP Request",
		zap.String("method", method),
		zap.String("url", url),
		zap.String("body", bodyStr),
	)
}

// Helper function to log HTTP responses
func (c *Client) logHTTPResponse(statusCode int, body interface{}) {
	if c.logger == nil {
		return
	}

	bodyStr := ""
	if body != nil {
		if jsonBytes, err := json.Marshal(body); err == nil {
			bodyStr = string(jsonBytes)
		}
	}

	c.logger.Debug("HTTP Response",
		zap.Int("status_code", statusCode),
		zap.String("body", bodyStr),
	)
}

// ReadBody reads and returns the response body as a string
func ReadBody(resp *http.Response) (string, error) {
	if resp == nil {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	return string(body), nil
}
