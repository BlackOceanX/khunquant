//go:build !mipsle && !netbsd && !(freebsd && arm)

package agent

import (
	"testing"

	"github.com/cryptoquantumwave/khunquant/pkg/providers/protocoltypes"
)

func TestProviderToSeahorseMessage_Basic(t *testing.T) {
	msg := protocoltypes.Message{
		Role:    "user",
		Content: "hello",
	}
	got := providerToSeahorseMessage(msg)
	if got.Role != "user" {
		t.Errorf("Role = %q, want user", got.Role)
	}
	if got.Content != "hello" {
		t.Errorf("Content = %q, want hello", got.Content)
	}
	if len(got.Parts) != 0 {
		t.Errorf("Parts = %d, want 0 for plain message", len(got.Parts))
	}
}

func TestProviderToSeahorseMessage_WithToolCallID(t *testing.T) {
	msg := protocoltypes.Message{
		Role:       "tool",
		Content:    "result-content",
		ToolCallID: "call-123",
	}
	got := providerToSeahorseMessage(msg)
	if len(got.Parts) != 1 {
		t.Fatalf("expected 1 part for tool result, got %d", len(got.Parts))
	}
	if got.Parts[0].Type != "tool_result" {
		t.Errorf("Parts[0].Type = %q, want tool_result", got.Parts[0].Type)
	}
	if got.Parts[0].ToolCallID != "call-123" {
		t.Errorf("Parts[0].ToolCallID = %q, want call-123", got.Parts[0].ToolCallID)
	}
}

func TestProviderToSeahorseMessage_WithToolCalls(t *testing.T) {
	msg := protocoltypes.Message{
		Role:    "assistant",
		Content: "",
		ToolCalls: []protocoltypes.ToolCall{
			{
				ID: "tc-1",
				Function: &protocoltypes.FunctionCall{
					Name:      "get_weather",
					Arguments: `{"city":"BKK"}`,
				},
			},
		},
	}
	got := providerToSeahorseMessage(msg)
	if len(got.Parts) != 1 {
		t.Fatalf("expected 1 part for tool call, got %d", len(got.Parts))
	}
	if got.Parts[0].Type != "tool_use" {
		t.Errorf("Parts[0].Type = %q, want tool_use", got.Parts[0].Type)
	}
	if got.Parts[0].Name != "get_weather" {
		t.Errorf("Parts[0].Name = %q, want get_weather", got.Parts[0].Name)
	}
}

func TestProviderToSeahorseMessage_WithMedia(t *testing.T) {
	msg := protocoltypes.Message{
		Role:    "user",
		Content: "look at this",
		Media:   []string{"file:///tmp/image.png"},
	}
	got := providerToSeahorseMessage(msg)
	if len(got.Parts) != 1 {
		t.Fatalf("expected 1 part for media, got %d", len(got.Parts))
	}
	if got.Parts[0].Type != "media" {
		t.Errorf("Parts[0].Type = %q, want media", got.Parts[0].Type)
	}
	if got.Parts[0].MediaURI != "file:///tmp/image.png" {
		t.Errorf("Parts[0].MediaURI = %q, want file:///tmp/image.png", got.Parts[0].MediaURI)
	}
}

func TestProviderToSeahorseMessage_ReasoningContent(t *testing.T) {
	msg := protocoltypes.Message{
		Role:             "assistant",
		Content:          "answer",
		ReasoningContent: "step-by-step",
	}
	got := providerToSeahorseMessage(msg)
	if got.ReasoningContent != "step-by-step" {
		t.Errorf("ReasoningContent = %q, want step-by-step", got.ReasoningContent)
	}
}
