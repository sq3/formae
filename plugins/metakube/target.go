// © 2025 Platform Engineering Labs Inc.
//
// SPDX-License-Identifier: FSL-1.1-ALv2

package main

import (
	"fmt"

	"github.com/platform-engineering-labs/formae/pkg/model"
	"github.com/platform-engineering-labs/formae/pkg/plugin/resource"
	"github.com/platform-engineering-labs/formae/plugins/metakube/pkg/config"
)

var TargetBehavior = targetBehavior{}
var _ resource.TargetBehavior = targetBehavior{}

type targetBehavior struct{}

// UpdateResourceAllowed is called when an update attempts to change a resource's target
func (t targetBehavior) UpdateResourceAllowed(current *model.Target, desired *model.Target) error {
	currentCfg := config.FromTarget(current)
	desiredCfg := config.FromTarget(desired)

	if currentCfg.Host != desiredCfg.Host {
		return fmt.Errorf("current target host and desired target host must match")
	}

	return nil
}

// UpdateTargetAllowed is called when updating a target
func (t targetBehavior) UpdateTargetAllowed(current *model.Target, desired *model.Target) error {
	currentCfg := config.FromTarget(current)
	desiredCfg := config.FromTarget(desired)

	if currentCfg.Host != desiredCfg.Host {
		return fmt.Errorf("target host cannot be modified")
	}

	return nil
}
