package tui

import (
	"fmt"
	"strings"

	"github.com/lwlee2608/tokentop/pkg/openrouter"
)

const (
	maxModels = 8
	maxKeys   = 10
)

func (m Model) orSection() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render(" OpenRouter"))
	b.WriteByte('\n')

	if m.orUsage == nil && m.orErr == "" {
		b.WriteString(dimStyle.Render("  Loading..."))
		b.WriteByte('\n')
		return b.String()
	}

	if m.orErr != "" {
		c := yellow
		if m.orUsage == nil {
			c = red
		}
		b.WriteString(pctStyle(c).Render(fmt.Sprintf("  ⚠️  %s (retry %d/%d)", m.orErr, m.orRetries, maxRetries)))
		b.WriteByte('\n')
		if m.orUsage == nil {
			return b.String()
		}
	}

	u := m.orUsage
	bw := m.barWidth()

	keyLabel := u.Key.Label
	switch {
	case u.Key.IsFreeTier:
		keyLabel += " (free tier)"
	case u.Key.IsManagementKey:
		keyLabel += " (management)"
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("  Key: %s", keyLabel)))
	b.WriteByte('\n')
	b.WriteByte('\n')

	if u.Key.Limit > 0 {
		usedPct := (u.Key.Limit - u.Key.LimitRemaining) / u.Key.Limit * 100
		b.WriteString(renderBar("Credit Limit", usedPct, bw,
			fmt.Sprintf("$%.4f remaining (resets %s)", u.Key.LimitRemaining, u.Key.LimitReset),
		))
		b.WriteByte('\n')
	}

	b.WriteString(dimStyle.Render(fmt.Sprintf("  Usage — Daily: $%.4f | Weekly: $%.4f | Monthly: $%.4f",
		u.Key.UsageDaily, u.Key.UsageWeekly, u.Key.UsageMonthly)))
	b.WriteByte('\n')

	if u.Key.IsManagementKey {
		b.WriteString(renderORCredits(u))
		b.WriteString(renderORActivity(u))
		b.WriteString(m.renderORModels(u))
		b.WriteString(renderORKeys(u))
	}

	b.WriteByte('\n')
	return b.String()
}

func renderORCredits(u *openrouter.Usage) string {
	if u.Credits == nil {
		return ""
	}
	return dimStyle.Render(fmt.Sprintf("\n  Credits — Total: $%.4f | Used: $%.4f | Remaining: $%.4f",
		u.Credits.Total, u.Credits.Used, u.Credits.Remaining)) + "\n"
}

func renderORActivity(u *openrouter.Usage) string {
	if u.Activity == nil {
		return ""
	}
	t := u.Activity.Totals
	var b strings.Builder
	b.WriteByte('\n')
	b.WriteString("  " + labelStyle.Render("Activity") + "\n")

	line := fmt.Sprintf("  Spend: $%.4f | Requests: %.0f | Tokens: %s prompt + %s completion",
		t.Spend, t.Requests, formatTokens(t.PromptTokens), formatTokens(t.CompletionTokens))
	if t.ReasoningTokens > 0 {
		line += fmt.Sprintf(" + %s reasoning", formatTokens(t.ReasoningTokens))
	}
	b.WriteString(dimStyle.Render(line))
	b.WriteByte('\n')
	return b.String()
}

func (m Model) renderORModels(u *openrouter.Usage) string {
	if u.Activity == nil || len(u.Activity.Models) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteByte('\n')
	b.WriteString("  " + labelStyle.Render("Top Models") + "\n")

	models := u.Activity.Models
	if len(models) > maxModels {
		models = models[:maxModels]
	}

	maxSpend := models[0].Spend
	barWidth := m.modelBarWidth()
	for _, m := range models {
		label := truncate(m.Model, 28)
		b.WriteString(dimStyle.Render(fmt.Sprintf("  %-28s  $%9.4f", label, m.Spend)))
		b.WriteByte('\n')
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			renderModelBar(m.Spend, maxSpend, barWidth),
			dimStyle.Render(fmt.Sprintf("%.0f req", m.Requests)),
		))
		b.WriteByte('\n')
	}
	return b.String()
}

func (m Model) modelBarWidth() int {
	w := m.width - 26
	if w < 12 {
		return 12
	}
	if w > 44 {
		return 44
	}
	return w
}

func renderModelBar(spend, maxSpend float64, width int) string {
	if width < 1 {
		width = 1
	}
	filled := width
	if maxSpend > 0 {
		filled = int(spend / maxSpend * float64(width))
	}
	if spend > 0 && filled == 0 {
		filled = 1
	}
	if filled > width {
		filled = width
	}

	return modelBarFilledStyle.Render(strings.Repeat("█", filled)) +
		modelBarEmptyStyle.Render(strings.Repeat("░", width-filled))
}

func renderORKeys(u *openrouter.Usage) string {
	if len(u.APIKeys) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteByte('\n')
	b.WriteString("  " + labelStyle.Render("API Keys") + "\n")
	b.WriteString(dimStyle.Render(fmt.Sprintf("  %-30s  %10s  %10s  %10s", "Key", "Daily", "Weekly", "Monthly")))
	b.WriteByte('\n')

	keys := u.APIKeys
	if len(keys) > maxKeys {
		keys = keys[:maxKeys]
	}
	for _, k := range keys {
		name := k.Label
		if name == "" {
			name = k.Name
		}
		b.WriteString(dimStyle.Render(fmt.Sprintf("  %-30s  $%9.4f  $%9.4f  $%9.4f",
			truncate(name, 30), k.UsageDaily, k.UsageWeekly, k.UsageMonthly)))
		b.WriteByte('\n')
	}
	return b.String()
}

func formatTokens(n float64) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", n/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", n/1_000)
	default:
		return fmt.Sprintf("%.0f", n)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
