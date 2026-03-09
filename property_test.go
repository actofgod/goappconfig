package goappconfig

import "testing"

func Test_toEnvVariable(t *testing.T) {
	var tests = map[string]string{
		"InputValue":       "INPUT_VALUE",
		"outputValue":      "OUTPUT_VALUE",
		"OTLPConfig":       "OTLPCONFIG",
		"OTLP_Config":      "OTLP__CONFIG",
		"path_to_value":    "PATH_TO_VALUE",
		"Value123Value321": "VALUE_123_VALUE_321",
		"AbCdEfGhIjKlMnOp": "AB_CD_EF_GH_IJ_KL_MN_OP",
	}
	for test, expected := range tests {
		t.Run(test, func(t *testing.T) {
			result := toEnvVariable(test)
			if result != expected {
				t.Errorf("%s != %s", result, expected)
			}
		})
	}
}
