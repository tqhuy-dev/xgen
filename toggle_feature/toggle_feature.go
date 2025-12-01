package toggle_feature

import (
	"fmt"
	"sync/atomic"

	"github.com/tqhuy-dev/xgen/utilities"
)

// ToggleFeatureConfig represents the configuration for a single feature toggle.
// It supports three modes of operation:
// 1. IsApplyAll: Enable for all users/requests
// 2. Ratio: Enable based on a probability ratio (0.0 to 1.0)
// 3. FieldEnable: Enable based on specific field-value matches
type ToggleFeatureConfig struct {
	// Ratio represents the probability of enabling the feature (0.0 to 1.0)
	// A value of 0.5 means 50% chance of being enabled
	Ratio float64

	// IsApplyAll enables the feature for all requests when true
	IsApplyAll bool

	// FieldEnable maps field names to allowed values
	// Feature is enabled if any field-value pair matches
	FieldEnable map[string][]string
}

// DataToggleFeature is a map of feature names to their configurations
type DataToggleFeature map[string]ToggleFeatureConfig

// ToggleFeatureClient provides thread-safe feature toggle functionality
// using atomic operations for concurrent access
type ToggleFeatureClient struct {
	atomic.Value
}

// NewToggleFeatureClient creates a new ToggleFeatureClient instance with empty configuration
func NewToggleFeatureClient() *ToggleFeatureClient {
	tf := &ToggleFeatureClient{}
	tf.Store(make(DataToggleFeature))
	return tf
}

// NewToggleFeatureWithConfig creates a new ToggleFeatureClient instance with provided configuration
func NewToggleFeatureWithConfig(config DataToggleFeature) (*ToggleFeatureClient, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	tf := &ToggleFeatureClient{}
	tf.Store(config)
	return tf, nil
}

// SetToggle updates the toggle configuration atomically
func (t *ToggleFeatureClient) SetToggle(config DataToggleFeature) error {
	if err := validateConfig(config); err != nil {
		return err
	}
	t.Store(config)
	return nil
}

// GetConfig returns a copy of the current configuration
func (t *ToggleFeatureClient) GetConfig() DataToggleFeature {
	if config, ok := t.Load().(DataToggleFeature); ok {
		return config
	}
	return make(DataToggleFeature)
}

// UseToggle checks if a feature should be enabled based on the provided data
// It evaluates conditions in the following order:
// 1. IsApplyAll: Returns true immediately if enabled for all
// 2. Ratio: Returns result based on probability if ratio > 0
// 3. FieldEnable: Checks if any field-value pair matches
// Returns false if feature is not found or no conditions match
func (t *ToggleFeatureClient) UseToggle(feature string, data map[string]string) bool {
	toggleConfig, ok := t.Load().(DataToggleFeature)
	if !ok {
		return false
	}

	toggle, ok := toggleConfig[feature]
	if !ok {
		return false
	}

	// Check IsApplyAll first (fastest check)
	if toggle.IsApplyAll {
		return true
	}

	// Check ratio-based toggle (probabilistic)
	if toggle.Ratio > 0 {
		return utilities.BoolByRatio(toggle.Ratio)
	}

	// Check field-based matching
	if len(toggle.FieldEnable) > 0 {
		return matchesFieldEnable(toggle.FieldEnable, data)
	}

	return false
}

// IsFeatureEnabled is an alias for UseToggle for better readability
func (t *ToggleFeatureClient) IsFeatureEnabled(feature string, data map[string]string) bool {
	return t.UseToggle(feature, data)
}

// matchesFieldEnable checks if any field-value pair in the toggle config matches the provided data
func matchesFieldEnable(fieldEnable map[string][]string, data map[string]string) bool {
	for field, allowedValues := range fieldEnable {
		if value, exists := data[field]; exists {
			for _, allowedValue := range allowedValues {
				if value == allowedValue {
					return true
				}
			}
		}
	}
	return false
}

// validateConfig validates the toggle feature configuration
func validateConfig(config DataToggleFeature) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	for featureName, cfg := range config {
		// Validate ratio is between 0 and 1
		if cfg.Ratio < 0 || cfg.Ratio > 1 {
			return fmt.Errorf("invalid ratio for feature '%s': %f (must be between 0.0 and 1.0)", featureName, cfg.Ratio)
		}
	}

	return nil
}
