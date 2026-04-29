package confluence

import (
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func strPtr(s string) *string { return &s }

func TestGetConfig_NilConnection(t *testing.T) {
	got := GetConfig(nil)
	if got.BaseUrl != nil || got.Token != nil || got.Username != nil || got.DeploymentType != nil {
		t.Error("GetConfig(nil) should return a zero-value confluenceConfig")
	}
}

func TestGetConfig_NilConnectionConfig(t *testing.T) {
	conn := &plugin.Connection{Config: nil}
	got := GetConfig(conn)
	if got.BaseUrl != nil || got.Token != nil {
		t.Error("GetConfig with nil Config should return a zero-value confluenceConfig")
	}
}

func TestGetConfig_WrongType(t *testing.T) {
	conn := &plugin.Connection{Config: "not-a-confluence-config"}
	got := GetConfig(conn)
	if got.BaseUrl != nil || got.Token != nil {
		t.Error("GetConfig with wrong Config type should return a zero-value confluenceConfig")
	}
}

func TestGetConfig_ValidConfig(t *testing.T) {
	cfg := &confluenceConfig{
		BaseUrl:        strPtr("https://example.atlassian.net/"),
		Username:       strPtr("user@example.com"),
		Token:          strPtr("secret-token"),
		DeploymentType: strPtr("cloud"),
	}
	conn := &plugin.Connection{Config: cfg}
	got := GetConfig(conn)

	if got.BaseUrl == nil || *got.BaseUrl != "https://example.atlassian.net/" {
		t.Errorf("BaseUrl = %v; want https://example.atlassian.net/", got.BaseUrl)
	}
	if got.Username == nil || *got.Username != "user@example.com" {
		t.Errorf("Username = %v; want user@example.com", got.Username)
	}
	if got.Token == nil || *got.Token != "secret-token" {
		t.Errorf("Token = %v; want secret-token", got.Token)
	}
	if got.DeploymentType == nil || *got.DeploymentType != "cloud" {
		t.Errorf("DeploymentType = %v; want cloud", got.DeploymentType)
	}
}
