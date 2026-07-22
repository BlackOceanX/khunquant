package api

// ─────────────────────────────────────────────────────────────────────────────
// Credential-card model presets
//
// UPDATE THIS FILE when new models are released (frontier models change ~monthly).
// Each provider block is independent — just edit labels and model IDs.
//
// ModelID format: {provider-prefix}/{api-model-id}
//   openai/       → OpenAI API (Codex CLI uses ChatGPT session auth)
//   anthropic/    → Anthropic API
//   antigravity/  → Google Antigravity (cloudcode-pa.googleapis.com)
//
// Last verified: 2026-06-28
// ─────────────────────────────────────────────────────────────────────────────

// oauthProviderModelPresets defines the clickable model chips shown on each
// credential card. Order matters — first entry is shown leftmost/topmost.
//
// Sources checked when last updated:
//
//	OpenAI:      developers.openai.com/codex/models
//	Anthropic:   platform.claude.com/docs/en/about-claude/models/overview
//	Antigravity: github.com/NoeFabris/opencode-antigravity-auth (PR #574)
var oauthProviderModelPresets = map[string][]oauthModelPreset{
	// ── OpenAI ────────────────────────────────────────────────────────────────
	// gpt-5.5  = current default for ChatGPT-authenticated Codex CLI sessions
	// gpt-5.4  = stable fallback
	// gpt-5.4-mini = lighter tasks, subagents, interactive edits
	oauthProviderOpenAI: {
		{Label: "GPT-5.5", ModelID: "openai/gpt-5.5"},
		{Label: "GPT-5.4", ModelID: "openai/gpt-5.4"},
		{Label: "GPT-5.4 Mini", ModelID: "openai/gpt-5.4-mini"},
	},

	// ── Anthropic ─────────────────────────────────────────────────────────────
	// claude-fable-5   = Mythos class, released 2026-06-09 (top tier, $10/$50/MTok)
	// claude-sonnet-4.6 = recommended production default ($3/$15/MTok)
	// claude-opus-4.8  = heavy reasoning ($5/$25/MTok)
	// claude-haiku-4.5 = low-latency, low-cost ($1/$5/MTok)
	oauthProviderAnthropic: {
		{Label: "Claude Fable 5", ModelID: "anthropic/claude-fable-5"},
		{Label: "Claude Sonnet 4.6", ModelID: "anthropic/claude-sonnet-4.6"},
		{Label: "Claude Opus 4.8", ModelID: "anthropic/claude-opus-4.8"},
		{Label: "Claude Haiku 4.5", ModelID: "anthropic/claude-haiku-4.5"},
	},

	// ── Google Antigravity ────────────────────────────────────────────────────
	// Model IDs MUST match what v1internal:fetchAvailableModels returns for the
	// account — the backend 404s ("NOT_FOUND: Requested entity was not found")
	// on any ID it doesn't expose. The provider strips the "antigravity/" prefix.
	// Labels mirror the Antigravity app's model picker; the IDs are the API IDs
	// (the API's displayName differs from the ID, e.g. id gemini-3.5-flash-low →
	// "Gemini 3.5 Flash (Medium)"). All entries verified live (2026-06-28) with a
	// tool-laden request — each returns a real tool call, not just a ping.
	// Gotchas: bare "gemini-3-pro"/"gemini-3.1-pro" do NOT exist (404); the High
	// Pro tier works via "gemini-pro-agent" — the literal "gemini-3.1-pro-high" ID
	// returns INVALID_ARGUMENT via this code path. gemini-3-flash-preview is a
	// Gemini CLI quota model (different API) — NOT valid here.
	// Note: the exact ID set is rollout-dependent. If a chip 404s, re-check with
	// fetchAvailableModels for the account and update the IDs here.
	//
	// Gemini 3.6 Flash added 2026-07-22 (Google GA'd it 2026-07-21; spotted in
	// Antigravity's own tieredModelIds config as a single ID, unlike 3.5 which
	// exposes separate -low/-extra-low/-agent IDs per tier). "Flash" replaces
	// 3.5 as Antigravity's tiered "flash" family; there is only one ID here —
	// thinking effort is controlled server-side/internally, not via separate
	// High/Medium/Low model IDs, so we only list one chip for it.
	oauthProviderGoogleAntigravity: {
		{Label: "Gemini 3.6 Flash", ModelID: "antigravity/gemini-3.6-flash-tiered"},
		{Label: "Gemini 3 Flash", ModelID: "antigravity/gemini-3-flash"},
		{Label: "Gemini 3.5 Flash (Medium)", ModelID: "antigravity/gemini-3.5-flash-low"},
		{Label: "Gemini 3.5 Flash (High)", ModelID: "antigravity/gemini-3-flash-agent"},
		{Label: "Gemini 3.5 Flash (Low)", ModelID: "antigravity/gemini-3.5-flash-extra-low"},
		{Label: "Gemini 3.1 Pro (High)", ModelID: "antigravity/gemini-pro-agent"},
		{Label: "Gemini 3.1 Pro (Low)", ModelID: "antigravity/gemini-3.1-pro-low"},
		{Label: "Claude Sonnet 4.6 (Thinking)", ModelID: "antigravity/claude-sonnet-4-6"},
		{Label: "Claude Opus 4.6 (Thinking)", ModelID: "antigravity/claude-opus-4-6-thinking"},
		{Label: "GPT-OSS 120B (Medium)", ModelID: "antigravity/gpt-oss-120b-medium"},
	},
}

// defaultModelForProvider returns the model ID to use when a provider is first
// connected and no existing model_list entry exists for it.
// Keep in sync with the first entry of each provider block above.
func defaultModelForProvider(provider string) string {
	switch provider {
	case oauthProviderOpenAI:
		return "openai/gpt-5.5"
	case oauthProviderAnthropic:
		return "anthropic/claude-sonnet-4.6" // sonnet is recommended production default
	case oauthProviderGoogleAntigravity:
		return "antigravity/gemini-3-flash"
	case oauthProviderGoogleGemini:
		return "gemini-code-assist/gemini-2.5-flash"
	default:
		return ""
	}
}
