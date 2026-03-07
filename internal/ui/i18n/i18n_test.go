package i18n

import (
	"testing"

	"gonime/internal/config"
)

func TestTranslationByLanguage(t *testing.T) {
	viCfg := &config.Config{Language: "vi"}
	enCfg := &config.Config{Language: "en"}

	if got := T(viCfg, "title_history"); got != "LỊCH SỬ XEM" {
		t.Fatalf("unexpected vi translation: %q", got)
	}
	if got := T(enCfg, "title_history"); got != "HISTORY" {
		t.Fatalf("unexpected en translation: %q", got)
	}
}

func TestTranslationFallback(t *testing.T) {
	cfg := &config.Config{Language: "en"}
	key := "non_existing_translation_key"
	if got := T(cfg, key); got != key {
		t.Fatalf("expected key fallback, got %q", got)
	}
}
