package main

import (
	"path/filepath"
	"testing"

	"github.com/BenedictKing/ccx/desktop/internal/configservice"
)

func newTestConfigService(t *testing.T) *configservice.Service {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	svc, err := configservice.New(filepath.Join(t.TempDir(), "ccx-data"))
	if err != nil {
		t.Fatalf("configservice.New failed: %v", err)
	}
	return svc
}

func TestSavedProviderKeyForPlanUsesExactPlan(t *testing.T) {
	svc := newTestConfigService(t)
	desktop := &DesktopService{configService: svc}

	if err := svc.SaveProviderKeyAsset(configservice.ProviderKeyAsset{
		Provider: configservice.ProviderMiMo,
		APIKey:   "tp-anthropic-key",
		PlanID:   "anthropic",
	}); err != nil {
		t.Fatalf("SaveProviderKeyAsset anthropic failed: %v", err)
	}
	if err := svc.SaveProviderKeyAsset(configservice.ProviderKeyAsset{
		Provider: configservice.ProviderMiMo,
		APIKey:   "tp-openai-key",
		PlanID:   "openai-chat",
	}); err != nil {
		t.Fatalf("SaveProviderKeyAsset openai-chat failed: %v", err)
	}

	if got := desktop.savedProviderKeyForPlan(configservice.ProviderMiMo, "anthropic"); got != "tp-anthropic-key" {
		t.Fatalf("anthropic key = %q", got)
	}
	if got := desktop.savedProviderKeyForPlan(configservice.ProviderMiMo, "openai-chat"); got != "tp-openai-key" {
		t.Fatalf("openai-chat key = %q", got)
	}
	if got := desktop.savedProviderKeyForPlan(configservice.ProviderMiMo, "token-cn"); got != "" {
		t.Fatalf("missing plan should not reuse another plan key, got %q", got)
	}
}

func TestSavedProviderKeyForPlanFallsBackToLegacyProviderKey(t *testing.T) {
	svc := newTestConfigService(t)
	desktop := &DesktopService{configService: svc}

	if err := svc.SaveProviderKeyAsset(configservice.ProviderKeyAsset{
		Provider: configservice.ProviderMiMo,
		APIKey:   "tp-legacy-key",
	}); err != nil {
		t.Fatalf("SaveProviderKeyAsset legacy failed: %v", err)
	}

	if got := desktop.savedProviderKeyForPlan(configservice.ProviderMiMo, "anthropic"); got != "tp-legacy-key" {
		t.Fatalf("legacy fallback key = %q", got)
	}
}
