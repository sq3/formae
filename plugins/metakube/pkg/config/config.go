// © 2025 Platform Engineering Labs Inc.
//
// SPDX-License-Identifier: FSL-1.1-ALv2

package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/platform-engineering-labs/formae/pkg/model"
)

const DefaultHost = "https://metakube.syseleven.de"

// Config represents MetaKube target configuration
type Config struct {
	// API endpoint (defaults to https://metakube.syseleven.de)
	Host string `json:"Host"`

	// Bearer token for authentication
	Token string `json:"Token"`

	// Optional: Project ID for scoping operations
	ProjectID string `json:"ProjectID"`

	// OpenStack Application Credential ID
	ApplicationCredentialID string `json:"ApplicationCredentialID,omitempty"`

	// OpenStack Application Credential Secret
	ApplicationCredentialSecret string `json:"ApplicationCredentialSecret,omitempty"`
}

// FromTarget extracts MetaKube configuration from a formae target
func FromTarget(target *model.Target) *Config {
	if target == nil || target.Config == nil {
		return &Config{
			Host: DefaultHost,
		}
	}

	cfg := &Config{}
	_ = json.Unmarshal(target.Config, cfg)

	// Set default host if not specified
	if cfg.Host == "" {
		cfg.Host = DefaultHost
	}

	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate(ctx context.Context) error {
	if c.Host == "" {
		return fmt.Errorf("MetaKube Host is required")
	}
	if c.Token == "" {
		return fmt.Errorf("MetaKube Token is required")
	}
	return nil
}
