// Copyright (c) Contributors to the Apptainer project, established as
//   Apptainer a Series of LF Projects LLC.
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package env

import (
	"os"
	"strings"
	"testing"

	"github.com/apptainer/apptainer/internal/pkg/test"
)

func TestSetFromList(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	tt := []struct {
		name    string
		environ []string
		wantErr bool
	}{
		{
			name: "all ok",
			environ: []string{
				"LD_LIBRARY_PATH=/.singularity.d/libs",
				"HOME=/home/tester",
				"PS1=test",
				"TERM=xterm-256color",
				"PATH=/usr/games:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				"LANG=C",
				"APPTAINER_CONTAINER=/tmp/lolcow.sif",
				"PWD=/tmp",
				"LC_ALL=C",
				"APPTAINER_NAME=lolcow.sif",
			},
			wantErr: false,
		},
		{
			name: "bad envs",
			environ: []string{
				"LD_LIBRARY_PATH=/.singularity.d/libs",
				"HOME=/home/tester",
				"PS1=test",
				"TERM=xterm-256color",
				"PATH=/usr/games:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				"LANG=C",
				"APPTAINER_CONTAINER=/tmp/lolcow.sif",
				"TEST",
				"LC_ALL=C",
				"APPTAINER_NAME=lolcow.sif",
			},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := SetFromList(tc.environ)
			if tc.wantErr && err == nil {
				t.Fatalf("Expected error, but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTrimApptainerKey(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	tests := []struct {
		name     string
		envKey   string
		expected string
	}{
		{
			name:     "good",
			envKey:   ApptainerPrefix + "TEST",
			expected: "TEST",
		},
		{
			name:     "bad",
			envKey:   "BADPREFIX_TEST",
			expected: "BADPREFIX_TEST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultKey := TrimApptainerKey(tt.envKey)
			if tt.expected != resultKey {
				t.Fatalf("Unexpected error for %s: got %s instead of %s", tt.name, resultKey, tt.expected)
			}
		})
	}
}

func TestGetenvLegacy(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	tests := []struct {
		name          string
		envVar        []string
		key           string
		legacyKey     string
		expectedValue string
	}{
		{
			name:          "apptainer only",
			envVar:        []string{"APPTAINER_TEST_APPTAINER=apptainer"},
			key:           "TEST_APPTAINER",
			legacyKey:     "TEST_APPTAINER",
			expectedValue: "apptainer",
		},
		{
			name:          "singularity only",
			envVar:        []string{"SINGULARITY_TEST_SINGULARITY=singularity"},
			key:           "TEST_SINGULARITY",
			legacyKey:     "TEST_SINGULARITY",
			expectedValue: "singularity",
		},
		{
			name:          "both",
			envVar:        []string{"APPTAINER_TEST_BOTH=apptainer", "SINGULARITY_TEST_BOTH=singularity"},
			key:           "TEST_BOTH",
			legacyKey:     "TEST_BOTH",
			expectedValue: "apptainer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.envVar {
				s := strings.SplitN(e, "=", 2)
				os.Setenv(s[0], s[1])
			}
			defer func() {
				for _, e := range tt.envVar {
					s := strings.SplitN(e, "=", 2)
					os.Unsetenv(s[0])
				}
			}()
			resultValue := GetenvLegacy(tt.key, tt.legacyKey)
			if tt.expectedValue != resultValue {
				t.Errorf("Unexpected error for %s: got %s instead of %s", tt.name, resultValue, tt.expectedValue)
			}
		})
	}
}
