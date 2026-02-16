package vinfo_test

import (
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/vinfo"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Parallel()

	result := vinfo.String()
	assert.NotEmpty(t, result)
}

func TestFormatDisplayVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
		commit  string
	}{
		{
			name:    "dev version",
			version: "dev",
			commit:  "none",
		},
		{
			name:    "semver only",
			version: "v1.2.3",
			commit:  "none",
		},
		{
			name:    "with commit",
			version: "v1.2.3",
			commit:  "abc123def",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := vinfo.FormatDisplayVersion(testCase.version, testCase.commit)
			assert.NotEmpty(t, result)
		})
	}
}

func TestShortCommit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		commit string
		want   string
	}{
		{
			name:   "short commit",
			commit: "abc123",
			want:   "abc123",
		},
		{
			name:   "long commit",
			commit: "abc123def456789",
			want:   "abc123d",
		},
		{
			name:   "empty commit",
			commit: "",
			want:   "",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := vinfo.ShortCommit(testCase.commit); got != testCase.want {
				t.Errorf("ShortCommit() = %v, want %v", got, testCase.want)
			}
		})
	}
}
