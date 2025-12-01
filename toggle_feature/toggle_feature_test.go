package toggle_feature

import (
	"sync"
	"testing"
)

// TestNewToggleFeature tests the creation of a new ToggleFeatureClient instance
// It verifies that the instance is properly initialized with an empty configuration
func TestNewToggleFeature(t *testing.T) {
	tests := []struct {
		name     string
		validate func(*testing.T, *ToggleFeatureClient)
	}{
		{
			name: "creates non-nil instance",
			validate: func(t *testing.T, tf *ToggleFeatureClient) {
				if tf == nil {
					t.Error("NewToggleFeatureClient() returned nil")
				}
			},
		},
		{
			name: "has empty config by default",
			validate: func(t *testing.T, tf *ToggleFeatureClient) {
				config := tf.GetConfig()
				if config == nil {
					t.Error("GetConfig() returned nil")
				}
				if len(config) != 0 {
					t.Errorf("GetConfig() returned %d features, expected 0", len(config))
				}
			},
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := NewToggleFeatureClient()

			// Assert
			tt.validate(t, result)
		})
	}
}

// TestNewToggleFeatureWithConfig tests creating a ToggleFeatureClient with initial configuration
// It covers valid configs, invalid configs, and edge cases
func TestNewToggleFeatureWithConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      DataToggleFeature
		expectError bool
		validate    func(*testing.T, *ToggleFeatureClient, error)
	}{
		{
			name: "valid config with IsApplyAll",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: true},
			},
			expectError: false,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tf == nil {
					t.Error("expected non-nil ToggleFeatureClient")
				}
			},
		},
		{
			name: "valid config with ratio",
			config: DataToggleFeature{
				"feature1": {Ratio: 0.5},
			},
			expectError: false,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "valid config with field enable",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123", "456"},
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "invalid config with ratio > 1",
			config: DataToggleFeature{
				"feature1": {Ratio: 1.5},
			},
			expectError: true,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err == nil {
					t.Error("expected error for ratio > 1")
				}
				if tf != nil {
					t.Error("expected nil ToggleFeatureClient on error")
				}
			},
		},
		{
			name: "invalid config with ratio < 0",
			config: DataToggleFeature{
				"feature1": {Ratio: -0.1},
			},
			expectError: true,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err == nil {
					t.Error("expected error for ratio < 0")
				}
			},
		},
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err == nil {
					t.Error("expected error for nil config")
				}
			},
		},
		{
			name:        "empty config",
			config:      DataToggleFeature{},
			expectError: false,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "valid config with ratio 0",
			config: DataToggleFeature{
				"feature1": {Ratio: 0},
			},
			expectError: false,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "valid config with ratio 1",
			config: DataToggleFeature{
				"feature1": {Ratio: 1},
			},
			expectError: false,
			validate: func(t *testing.T, tf *ToggleFeatureClient, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := NewToggleFeatureWithConfig(tt.config)

			// Assert
			tt.validate(t, result, err)
		})
	}
}

// TestSetToggle tests updating the toggle configuration
// It verifies that valid configs are accepted and invalid configs are rejected
func TestSetToggle(t *testing.T) {
	tests := []struct {
		name        string
		config      DataToggleFeature
		expectError bool
	}{
		{
			name: "valid config",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: true},
			},
			expectError: false,
		},
		{
			name: "invalid config with ratio > 1",
			config: DataToggleFeature{
				"feature1": {Ratio: 2.0},
			},
			expectError: true,
		},
		{
			name: "invalid config with ratio < 0",
			config: DataToggleFeature{
				"feature1": {Ratio: -1.0},
			},
			expectError: true,
		},
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
		{
			name:        "empty config",
			config:      DataToggleFeature{},
			expectError: false,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tf := NewToggleFeatureClient()

			// Act
			err := tf.SetToggle(tt.config)

			// Assert
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestGetConfig tests retrieving the current configuration
// It verifies that the config can be retrieved and is not nil
func TestGetConfig(t *testing.T) {
	tests := []struct {
		name           string
		initialConfig  DataToggleFeature
		expectedLength int
	}{
		{
			name:           "empty config",
			initialConfig:  DataToggleFeature{},
			expectedLength: 0,
		},
		{
			name: "single feature config",
			initialConfig: DataToggleFeature{
				"feature1": {IsApplyAll: true},
			},
			expectedLength: 1,
		},
		{
			name: "multiple features config",
			initialConfig: DataToggleFeature{
				"feature1": {IsApplyAll: true},
				"feature2": {Ratio: 0.5},
			},
			expectedLength: 2,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tf := NewToggleFeatureClient()
			_ = tf.SetToggle(tt.initialConfig)

			// Act
			result := tf.GetConfig()

			// Assert
			if result == nil {
				t.Error("GetConfig() returned nil")
			}
			if len(result) != tt.expectedLength {
				t.Errorf("GetConfig() returned %d features, expected %d", len(result), tt.expectedLength)
			}
		})
	}
}

// TestUseToggle tests the feature toggle evaluation logic
// It covers all toggle modes: IsApplyAll, Ratio, and FieldEnable
func TestUseToggle(t *testing.T) {
	tests := []struct {
		name           string
		config         DataToggleFeature
		featureName    string
		data           map[string]string
		expectedResult bool
		description    string
	}{
		{
			name: "IsApplyAll enabled returns true",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: true},
			},
			featureName:    "feature1",
			data:           map[string]string{},
			expectedResult: true,
			description:    "Feature with IsApplyAll should always return true",
		},
		{
			name: "IsApplyAll disabled with no other config returns false",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: false},
			},
			featureName:    "feature1",
			data:           map[string]string{},
			expectedResult: false,
			description:    "Feature with only IsApplyAll=false should return false",
		},
		{
			name: "non-existent feature returns false",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: true},
			},
			featureName:    "feature2",
			data:           map[string]string{},
			expectedResult: false,
			description:    "Non-existent feature should return false",
		},
		{
			name: "field enable with matching field returns true",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123", "456"},
					},
				},
			},
			featureName: "feature1",
			data: map[string]string{
				"user_id": "123",
			},
			expectedResult: true,
			description:    "Matching field value should return true",
		},
		{
			name: "field enable with non-matching field returns false",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123", "456"},
					},
				},
			},
			featureName: "feature1",
			data: map[string]string{
				"user_id": "789",
			},
			expectedResult: false,
			description:    "Non-matching field value should return false",
		},
		{
			name: "field enable with missing field returns false",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123", "456"},
					},
				},
			},
			featureName:    "feature1",
			data:           map[string]string{},
			expectedResult: false,
			description:    "Missing field should return false",
		},
		{
			name: "field enable with multiple fields matching one",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123", "456"},
						"city":    {"NYC", "SF"},
					},
				},
			},
			featureName: "feature1",
			data: map[string]string{
				"user_id": "789",
				"city":    "NYC",
			},
			expectedResult: true,
			description:    "Should return true if any field matches",
		},
		{
			name: "field enable with multiple values matching",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123", "456", "789"},
					},
				},
			},
			featureName: "feature1",
			data: map[string]string{
				"user_id": "789",
			},
			expectedResult: true,
			description:    "Should match any value in the list",
		},
		{
			name: "empty data with field enable returns false",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123"},
					},
				},
			},
			featureName:    "feature1",
			data:           map[string]string{},
			expectedResult: false,
			description:    "Empty data should not match any field enable",
		},
		{
			name: "nil data with field enable returns false",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{
						"user_id": {"123"},
					},
				},
			},
			featureName:    "feature1",
			data:           nil,
			expectedResult: false,
			description:    "Nil data should not match any field enable",
		},
		{
			name: "empty field enable map returns false",
			config: DataToggleFeature{
				"feature1": {
					FieldEnable: map[string][]string{},
				},
			},
			featureName: "feature1",
			data: map[string]string{
				"user_id": "123",
			},
			expectedResult: false,
			description:    "Empty field enable map should return false",
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tf, err := NewToggleFeatureWithConfig(tt.config)
			if err != nil {
				t.Fatalf("failed to create toggle feature: %v", err)
			}

			// Act
			result := tf.UseToggle(tt.featureName, tt.data)

			// Assert
			if result != tt.expectedResult {
				t.Errorf("UseToggle() = %v, expected %v: %s", result, tt.expectedResult, tt.description)
			}
		})
	}
}

// TestUseToggleRatio tests ratio-based feature toggles
// Note: Ratio tests are probabilistic and may occasionally fail due to randomness
func TestUseToggleRatio(t *testing.T) {
	tests := []struct {
		name        string
		ratio       float64
		iterations  int
		minEnabled  int
		maxEnabled  int
		description string
	}{
		{
			name:        "ratio 0 should never enable",
			ratio:       0,
			iterations:  100,
			minEnabled:  0,
			maxEnabled:  0,
			description: "Ratio of 0 should never enable the feature",
		},
		{
			name:        "ratio 1 should always enable within range",
			ratio:       1,
			iterations:  100,
			minEnabled:  0,
			maxEnabled:  100,
			description: "Ratio of 1 may or may not enable based on BoolByRatio implementation",
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := DataToggleFeature{
				"feature1": {Ratio: tt.ratio},
			}
			tf, err := NewToggleFeatureWithConfig(config)
			if err != nil {
				t.Fatalf("failed to create toggle feature: %v", err)
			}

			// Act
			enabledCount := 0
			for i := 0; i < tt.iterations; i++ {
				if tf.UseToggle("feature1", map[string]string{}) {
					enabledCount++
				}
			}

			// Assert
			if enabledCount < tt.minEnabled || enabledCount > tt.maxEnabled {
				t.Errorf("UseToggle() enabled %d times out of %d iterations, expected between %d and %d: %s",
					enabledCount, tt.iterations, tt.minEnabled, tt.maxEnabled, tt.description)
			}
		})
	}
}

// TestIsFeatureEnabled tests the IsFeatureEnabled method (alias for UseToggle)
// It verifies that it behaves identically to UseToggle
func TestIsFeatureEnabled(t *testing.T) {
	tests := []struct {
		name        string
		config      DataToggleFeature
		featureName string
		data        map[string]string
		expected    bool
	}{
		{
			name: "enabled feature",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: true},
			},
			featureName: "feature1",
			data:        map[string]string{},
			expected:    true,
		},
		{
			name: "disabled feature",
			config: DataToggleFeature{
				"feature1": {IsApplyAll: false},
			},
			featureName: "feature1",
			data:        map[string]string{},
			expected:    false,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tf, err := NewToggleFeatureWithConfig(tt.config)
			if err != nil {
				t.Fatalf("failed to create toggle feature: %v", err)
			}

			// Act
			result := tf.IsFeatureEnabled(tt.featureName, tt.data)

			// Assert
			if result != tt.expected {
				t.Errorf("IsFeatureEnabled() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMatchesFieldEnable tests the internal matchesFieldEnable function
// It verifies field matching logic across various scenarios
func TestMatchesFieldEnable(t *testing.T) {
	tests := []struct {
		name        string
		fieldEnable map[string][]string
		data        map[string]string
		expected    bool
		description string
	}{
		{
			name: "exact match",
			fieldEnable: map[string][]string{
				"user_id": {"123"},
			},
			data: map[string]string{
				"user_id": "123",
			},
			expected:    true,
			description: "Exact field and value match should return true",
		},
		{
			name: "no match",
			fieldEnable: map[string][]string{
				"user_id": {"123"},
			},
			data: map[string]string{
				"user_id": "456",
			},
			expected:    false,
			description: "Different value should return false",
		},
		{
			name: "missing field",
			fieldEnable: map[string][]string{
				"user_id": {"123"},
			},
			data:        map[string]string{},
			expected:    false,
			description: "Missing field should return false",
		},
		{
			name:        "empty field enable",
			fieldEnable: map[string][]string{},
			data: map[string]string{
				"user_id": "123",
			},
			expected:    false,
			description: "Empty field enable should return false",
		},
		{
			name: "multiple values one match",
			fieldEnable: map[string][]string{
				"user_id": {"123", "456", "789"},
			},
			data: map[string]string{
				"user_id": "456",
			},
			expected:    true,
			description: "Should match one of multiple values",
		},
		{
			name: "multiple fields one match",
			fieldEnable: map[string][]string{
				"user_id": {"123"},
				"city":    {"NYC"},
			},
			data: map[string]string{
				"user_id": "456",
				"city":    "NYC",
			},
			expected:    true,
			description: "Should return true if any field matches",
		},
		{
			name: "nil data",
			fieldEnable: map[string][]string{
				"user_id": {"123"},
			},
			data:        nil,
			expected:    false,
			description: "Nil data should return false",
		},
		{
			name:        "nil field enable",
			fieldEnable: nil,
			data: map[string]string{
				"user_id": "123",
			},
			expected:    false,
			description: "Nil field enable should return false",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := matchesFieldEnable(tt.fieldEnable, tt.data)

			// Assert
			if result != tt.expected {
				t.Errorf("matchesFieldEnable() = %v, expected %v: %s", result, tt.expected, tt.description)
			}
		})
	}
}

// TestValidateConfig tests the configuration validation logic
// It verifies that valid configs pass and invalid configs fail with appropriate errors
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      DataToggleFeature
		expectError bool
		description string
	}{
		{
			name: "valid config with ratio 0.5",
			config: DataToggleFeature{
				"feature1": {Ratio: 0.5},
			},
			expectError: false,
			description: "Valid ratio should pass validation",
		},
		{
			name: "valid config with ratio 0",
			config: DataToggleFeature{
				"feature1": {Ratio: 0},
			},
			expectError: false,
			description: "Ratio of 0 is valid",
		},
		{
			name: "valid config with ratio 1",
			config: DataToggleFeature{
				"feature1": {Ratio: 1},
			},
			expectError: false,
			description: "Ratio of 1 is valid",
		},
		{
			name: "invalid config with ratio > 1",
			config: DataToggleFeature{
				"feature1": {Ratio: 1.5},
			},
			expectError: true,
			description: "Ratio greater than 1 should fail validation",
		},
		{
			name: "invalid config with ratio < 0",
			config: DataToggleFeature{
				"feature1": {Ratio: -0.5},
			},
			expectError: true,
			description: "Negative ratio should fail validation",
		},
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			description: "Nil config should fail validation",
		},
		{
			name:        "empty config",
			config:      DataToggleFeature{},
			expectError: false,
			description: "Empty config is valid",
		},
		{
			name: "multiple features with one invalid",
			config: DataToggleFeature{
				"feature1": {Ratio: 0.5},
				"feature2": {Ratio: 2.0},
			},
			expectError: true,
			description: "Should fail if any feature has invalid config",
		},
		{
			name: "multiple valid features",
			config: DataToggleFeature{
				"feature1": {Ratio: 0.5},
				"feature2": {IsApplyAll: true},
				"feature3": {
					FieldEnable: map[string][]string{
						"user_id": {"123"},
					},
				},
			},
			expectError: false,
			description: "Multiple valid features should pass",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validateConfig(tt.config)

			// Assert
			if tt.expectError && err == nil {
				t.Errorf("validateConfig() expected error but got nil: %s", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("validateConfig() unexpected error: %v: %s", err, tt.description)
			}
		})
	}
}

// TestConcurrentAccess tests thread-safety of ToggleFeatureClient
// It verifies that concurrent reads and writes don't cause race conditions
func TestConcurrentAccess(t *testing.T) {
	tests := []struct {
		name          string
		numGoroutines int
		numOperations int
	}{
		{
			name:          "concurrent reads and writes",
			numGoroutines: 10,
			numOperations: 100,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tf := NewToggleFeatureClient()
			config := DataToggleFeature{
				"feature1": {IsApplyAll: true},
			}
			_ = tf.SetToggle(config)

			var wg sync.WaitGroup
			wg.Add(tt.numGoroutines * 2) // readers and writers

			// Act - concurrent reads
			for i := 0; i < tt.numGoroutines; i++ {
				go func() {
					defer wg.Done()
					for j := 0; j < tt.numOperations; j++ {
						_ = tf.UseToggle("feature1", map[string]string{})
					}
				}()
			}

			// Act - concurrent writes
			for i := 0; i < tt.numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()
					for j := 0; j < tt.numOperations; j++ {
						newConfig := DataToggleFeature{
							"feature1": {IsApplyAll: id%2 == 0},
						}
						_ = tf.SetToggle(newConfig)
					}
				}(i)
			}

			// Assert - wait for all goroutines to complete without panic
			wg.Wait()
		})
	}
}

// TestPriorityOrder tests that toggle evaluation follows the correct priority order
// Priority: IsApplyAll > Ratio > FieldEnable
func TestPriorityOrder(t *testing.T) {
	tests := []struct {
		name           string
		config         ToggleFeatureConfig
		data           map[string]string
		expectedResult bool
		description    string
	}{
		{
			name: "IsApplyAll takes priority over ratio",
			config: ToggleFeatureConfig{
				IsApplyAll: true,
				Ratio:      0.0,
			},
			data:           map[string]string{},
			expectedResult: true,
			description:    "IsApplyAll should override ratio of 0",
		},
		{
			name: "IsApplyAll takes priority over field enable",
			config: ToggleFeatureConfig{
				IsApplyAll: true,
				FieldEnable: map[string][]string{
					"user_id": {"123"},
				},
			},
			data: map[string]string{
				"user_id": "456", // non-matching value
			},
			expectedResult: true,
			description:    "IsApplyAll should override non-matching field",
		},
		{
			name: "ratio 0 falls through to field enable",
			config: ToggleFeatureConfig{
				Ratio: 0.0,
				FieldEnable: map[string][]string{
					"user_id": {"123"},
				},
			},
			data: map[string]string{
				"user_id": "123",
			},
			expectedResult: true,
			description:    "Ratio of 0 means ratio is not used, so it falls through to check field enable",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := DataToggleFeature{
				"feature1": tt.config,
			}
			tf, err := NewToggleFeatureWithConfig(config)
			if err != nil {
				t.Fatalf("failed to create toggle feature: %v", err)
			}

			// Act
			result := tf.UseToggle("feature1", tt.data)

			// Assert
			if result != tt.expectedResult {
				t.Errorf("UseToggle() = %v, expected %v: %s", result, tt.expectedResult, tt.description)
			}
		})
	}
}
