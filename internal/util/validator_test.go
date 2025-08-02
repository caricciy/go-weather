package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testRow struct {
	name     string
	cep      string
	expected bool
}

func TestCheckCEPIsValid(t *testing.T) {
	testTable := []testRow{
		{"Valid CEP", "12345678", true},
		{"Invalid CEP - too short", "12345", false},
		{"Invalid CEP - too long", "123456789", false},
		{"Invalid CEP - contains letters", "1234abcd", false},
		{"Invalid CEP - special characters", "1234!@#$", false},
		{"Invalid CEP - empty string", "", false},
	}

	for _, tr := range testTable {
		t.Run(tr.name, func(t *testing.T) {
			result := CheckCEPIsValid(tr.cep)
			assert.Equal(t, tr.expected, result, "Expected result for %s to be %v", tr.cep, tr.expected)
		})
	}
}
