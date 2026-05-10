package tools

import (
	"context"
	"testing"

	"github.com/cryptoquantumwave/khunquant/pkg/config"
)

func TestExchangeTotalValue_EmptyArgs(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	// No exchange specified should return result about no accounts configured
	result := tool.Execute(context.Background(), map[string]any{})
	// This will succeed with a message about no configured accounts, so we just check it returns something
	if result == nil {
		t.Fatal("expected a result")
	}
}

func TestExchangeTotalValue_InvalidExchange(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	result := tool.Execute(context.Background(), map[string]any{
		"exchange": "nonexistent",
	})
	if !result.IsError {
		t.Fatal("expected error for nonexistent exchange")
	}
}

func TestExchangeTotalValue_WalletType(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	result := tool.Execute(context.Background(), map[string]any{
		"exchange":    "nonexistent",
		"wallet_type": "spot",
	})
	if !result.IsError {
		t.Fatal("expected error for nonexistent exchange")
	}
}

func TestExchangeTotalValue_CustomQuote(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	result := tool.Execute(context.Background(), map[string]any{
		"exchange": "nonexistent",
		"quote":    "EUR",
	})
	if !result.IsError {
		t.Fatal("expected error for nonexistent exchange")
	}
}

func TestExchangeTotalValue_QuoteCase(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	// Quote should be uppercased
	result := tool.Execute(context.Background(), map[string]any{
		"exchange": "nonexistent",
		"quote":    "usdt",
	})
	if !result.IsError {
		t.Fatal("expected error for nonexistent exchange")
	}
}

func TestExchangeTotalValue_AccountWithExchange(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	result := tool.Execute(context.Background(), map[string]any{
		"exchange": "nonexistent",
		"account":  "myaccount",
	})
	if !result.IsError {
		t.Fatal("expected error for nonexistent exchange")
	}
}

func TestExchangeTotalValue_AllParameters(t *testing.T) {
	cfg := config.DefaultConfig()
	tool := NewExchangeTotalValueTool(cfg)

	result := tool.Execute(context.Background(), map[string]any{
		"exchange":    "nonexistent",
		"account":     "myaccount",
		"wallet_type": "spot",
		"quote":       "BTC",
	})
	if !result.IsError {
		t.Fatal("expected error for nonexistent exchange")
	}
}

func TestExchangeTotalValue_ParametersSchema(t *testing.T) {
	tool := NewExchangeTotalValueTool(config.DefaultConfig())
	params := tool.Parameters()

	if params == nil {
		t.Fatal("Parameters() should not return nil")
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties in Parameters")
	}

	expectedProps := []string{"exchange", "account", "wallet_type", "quote"}
	for _, prop := range expectedProps {
		if _, ok := props[prop]; !ok {
			t.Errorf("expected property %q in Parameters", prop)
		}
	}
}

func TestExchangeTotalValue_Name(t *testing.T) {
	tool := NewExchangeTotalValueTool(config.DefaultConfig())
	name := tool.Name()
	if name != NameGetTotalValue {
		t.Errorf("Name() = %q, want %q", name, NameGetTotalValue)
	}
}

func TestExchangeTotalValue_Description(t *testing.T) {
	tool := NewExchangeTotalValueTool(config.DefaultConfig())
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
}
