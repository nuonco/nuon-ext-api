package config

import "testing"

func TestLoadReadsInstallIDFromEnv(t *testing.T) {
	t.Setenv("NUON_INSTALL_ID", "inst_123")

	cfg := Load()
	if cfg.InstallID != "inst_123" {
		t.Fatalf("expected InstallID to be %q, got %q", "inst_123", cfg.InstallID)
	}
}
