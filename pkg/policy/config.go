// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package policy

import (
	"time"

	"github.com/cilium/cilium/pkg/labels"
	"github.com/cilium/cilium/pkg/lock"
	"github.com/cilium/cilium/pkg/logging"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/source"
)

var (
	log          = logging.DefaultLogger.WithField(logfields.LogSubsys, "policy")
	mutex        lock.RWMutex // Protects enablePolicy
	enablePolicy string       // Whether policy enforcement is enabled.
)

// SetPolicyEnabled sets the policy enablement configuration. Valid values are:
// - endpoint.AlwaysEnforce
// - endpoint.NeverEnforce
// - endpoint.DefaultEnforcement
func SetPolicyEnabled(val string) {
	mutex.Lock()
	enablePolicy = val
	mutex.Unlock()
}

// GetPolicyEnabled returns the policy enablement configuration
func GetPolicyEnabled() string {
	mutex.RLock()
	val := enablePolicy
	mutex.RUnlock()
	return val
}

// AddOptions are options which can be passed to PolicyAdd
type AddOptions struct {
	// Replace if true indicates that existing rules with identical labels should be replaced
	Replace bool
	// ReplaceWithLabels if present indicates that existing rules with the
	// given LabelArray should be deleted.
	ReplaceWithLabels labels.LabelArray
	// Generated should be set as true to signalize a the policy being inserted
	// was generated by cilium-agent, e.g. dns poller.
	Generated bool

	// The source of this policy, one of api, fqdn or k8s
	Source source.Source

	// The time the policy initially began to be processed in Cilium, such as when the
	// policy was received from the API server.
	ProcessingStartTime time.Time
}
