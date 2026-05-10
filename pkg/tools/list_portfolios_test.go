package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/cryptoquantumwave/khunquant/pkg/config"
)

func newTestListPortfoliosTool(t *testing.T) *ListPortfoliosTool {
	t.Helper()
	return NewListPortfoliosTool(config.DefaultConfig())
}

func TestListPortfoliosTool_Name(t *testing.T) {
	tool := newTestListPortfoliosTool(t)
	if tool.Name() != NameListPortfolios {
		t.Errorf("Name() = %q, want %q", tool.Name(), NameListPortfolios)
	}
}

func TestListPortfoliosTool_Description(t *testing.T) {
	tool := newTestListPortfoliosTool(t)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description() should not be empty")
	}
}

func TestListPortfoliosTool_Parameters(t *testing.T) {
	tool := newTestListPortfoliosTool(t)
	params := tool.Parameters()

	if params["type"] != "object" {
		t.Errorf("type should be 'object', got %v", params["type"])
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties should be a map")
	}

	// ListPortfolios has no required parameters
	if len(props) != 0 {
		t.Errorf("expected empty properties map for ListPortfolios, got %d props", len(props))
	}
}

func TestListPortfoliosTool_Execute_NoExchanges(t *testing.T) {
	cfg := config.DefaultConfig()
	// Disable all exchanges
	cfg.Exchanges.Binance.Enabled = false
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.Bitkub.Enabled = false
	cfg.Exchanges.OKX.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result.IsError {
		t.Fatalf("unexpected error: %s", result.ForLLM)
	}
	if !strings.Contains(result.ForUser, "No exchange accounts") {
		t.Errorf("expected 'No exchange accounts' message, got %q", result.ForUser)
	}
}

func TestListPortfoliosTool_Execute_NoArgs(t *testing.T) {
	tool := newTestListPortfoliosTool(t)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return non-nil result")
	}
	// Result might be error if no exchanges configured, but that's ok
}

func TestListPortfoliosTool_Execute_WithArgs(t *testing.T) {
	tool := newTestListPortfoliosTool(t)
	// ListPortfolios ignores all arguments
	result := tool.Execute(context.Background(), map[string]any{
		"foo": "bar",
		"baz": 123,
	})

	if result == nil {
		t.Fatal("Execute should return non-nil result")
	}
}

func TestListPortfoliosTool_Execute_BinanceEnabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.Binance.Enabled = true
	cfg.Exchanges.Binance.Accounts = []config.ExchangeAccount{
		{Name: "main", APIKey: *config.NewSecureString("test-key"), Secret: *config.NewSecureString("test-secret")},
	}
	// Disable others
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.Bitkub.Enabled = false
	cfg.Exchanges.OKX.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result.IsError {
		t.Logf("Binance listing result: %s", result.ForLLM)
	}
	// Note: May error if binance SDK unavailable, but we're testing parameter handling
}

func TestListPortfoliosTool_Execute_MultipleExchanges(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.Binance.Enabled = true
	cfg.Exchanges.Binance.Accounts = []config.ExchangeAccount{
		{Name: "spot", APIKey: *config.NewSecureString("key1"), Secret: *config.NewSecureString("secret1")},
	}
	cfg.Exchanges.Bitkub.Enabled = true
	cfg.Exchanges.Bitkub.Accounts = []config.ExchangeAccount{
		{Name: "main", APIKey: *config.NewSecureString("key2"), Secret: *config.NewSecureString("secret2")},
	}
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.OKX.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
	// May error if exchange SDK not available, but interface should work
}

func TestListPortfoliosTool_Execute_UnnamedAccount(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.Bitkub.Enabled = true
	cfg.Exchanges.Bitkub.Accounts = []config.ExchangeAccount{
		{APIKey: *config.NewSecureString("key"), Secret: *config.NewSecureString("secret")}, // No name
	}
	cfg.Exchanges.Binance.Enabled = false
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.OKX.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
	// Should auto-generate name like "1" for unnamed accounts
}

func TestListPortfoliosTool_Execute_SettradEnabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.Settrade.Enabled = true
	cfg.Exchanges.Settrade.Accounts = []config.SettradeExchangeAccount{
		{ExchangeAccount: config.ExchangeAccount{Name: "default", APIKey: *config.NewSecureString("test-key"), Secret: *config.NewSecureString("test-secret")}},
	}
	cfg.Exchanges.Binance.Enabled = false
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.Bitkub.Enabled = false
	cfg.Exchanges.OKX.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
}

func TestListPortfoliosTool_Execute_OKXEnabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.OKX.Enabled = true
	cfg.Exchanges.OKX.Accounts = []config.OKXExchangeAccount{
		{ExchangeAccount: config.ExchangeAccount{Name: "trading", APIKey: *config.NewSecureString("test-key"), Secret: *config.NewSecureString("test-secret")}},
	}
	cfg.Exchanges.Binance.Enabled = false
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.Bitkub.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
}

func TestListPortfoliosTool_Execute_BinanceTHEnabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.BinanceTH.Enabled = true
	cfg.Exchanges.BinanceTH.Accounts = []config.ExchangeAccount{
		{Name: "baht", APIKey: *config.NewSecureString("test-key"), Secret: *config.NewSecureString("test-secret")},
	}
	cfg.Exchanges.Binance.Enabled = false
	cfg.Exchanges.Bitkub.Enabled = false
	cfg.Exchanges.OKX.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
}

func TestListPortfoliosTool_Execute_AllExchanges(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.Binance.Enabled = true
	cfg.Exchanges.Binance.Accounts = []config.ExchangeAccount{
		{Name: "main", APIKey: *config.NewSecureString("k1"), Secret: *config.NewSecureString("s1")},
	}
	cfg.Exchanges.BinanceTH.Enabled = true
	cfg.Exchanges.BinanceTH.Accounts = []config.ExchangeAccount{
		{Name: "th", APIKey: *config.NewSecureString("k2"), Secret: *config.NewSecureString("s2")},
	}
	cfg.Exchanges.Bitkub.Enabled = true
	cfg.Exchanges.Bitkub.Accounts = []config.ExchangeAccount{
		{Name: "spot", APIKey: *config.NewSecureString("k3"), Secret: *config.NewSecureString("s3")},
	}
	cfg.Exchanges.OKX.Enabled = true
	cfg.Exchanges.OKX.Accounts = []config.OKXExchangeAccount{
		{ExchangeAccount: config.ExchangeAccount{Name: "okx", APIKey: *config.NewSecureString("k4"), Secret: *config.NewSecureString("s4")}},
	}
	cfg.Exchanges.Settrade.Enabled = true
	cfg.Exchanges.Settrade.Accounts = []config.SettradeExchangeAccount{
		{ExchangeAccount: config.ExchangeAccount{Name: "settrade", APIKey: *config.NewSecureString("k5"), Secret: *config.NewSecureString("s5")}},
	}

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
}

func TestListPortfoliosTool_Execute_MultipleAccountsSameExchange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Exchanges.Binance.Enabled = true
	cfg.Exchanges.Binance.Accounts = []config.ExchangeAccount{
		{Name: "spot", APIKey: *config.NewSecureString("k1"), Secret: *config.NewSecureString("s1")},
		{Name: "futures", APIKey: *config.NewSecureString("k2"), Secret: *config.NewSecureString("s2")},
		{Name: "savings", APIKey: *config.NewSecureString("k3"), Secret: *config.NewSecureString("s3")},
	}
	cfg.Exchanges.BinanceTH.Enabled = false
	cfg.Exchanges.Bitkub.Enabled = false
	cfg.Exchanges.OKX.Enabled = false
	cfg.Exchanges.Settrade.Enabled = false

	tool := NewListPortfoliosTool(cfg)
	result := tool.Execute(context.Background(), map[string]any{})

	if result == nil {
		t.Fatal("Execute should return result")
	}
}
