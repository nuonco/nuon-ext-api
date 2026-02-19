package config

import (
	"os"

	"github.com/nuonco/nuon-ext-api/internal/debug"
)

type Config struct {
	APIURL     string
	APIToken   string
	OrgID      string
	AppID      string
	InstallID  string
	ConfigFile string
	ExtName    string
	ExtDir     string
}

func Load() *Config {
	cfg := &Config{
		APIURL:     os.Getenv("NUON_API_URL"),
		APIToken:   os.Getenv("NUON_API_TOKEN"),
		OrgID:      os.Getenv("NUON_ORG_ID"),
		AppID:      os.Getenv("NUON_APP_ID"),
		InstallID:  os.Getenv("NUON_INSTALL_ID"),
		ConfigFile: os.Getenv("NUON_CONFIG_FILE"),
		ExtName:    os.Getenv("NUON_EXT_NAME"),
		ExtDir:     os.Getenv("NUON_EXT_DIR"),
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.nuon.co"
	}

	debug.Log("config: api_url=%s org_id=%s app_id=%s install_id=%s token=%s",
		cfg.APIURL, cfg.OrgID, cfg.AppID, cfg.InstallID, maskToken(cfg.APIToken))

	return cfg
}

func maskToken(token string) string {
	if token == "" {
		return "(empty)"
	}
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
