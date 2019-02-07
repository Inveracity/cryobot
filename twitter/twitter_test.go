package twitter

import (
	"cryobot/config"
	"testing"
)

// Check if "test is in the slice of blacklisted entries"
func TestIsblacklisted(t *testing.T) {
	cfg := config.Config{}
	cfg.Twitter.Blacklist = []string{"test", "test2", "test3"}
	account := "test"
	output := isBlacklisted(account, cfg)
	if !output {
		t.Error()
	}
}
