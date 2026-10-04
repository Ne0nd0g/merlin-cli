/*
Merlin is a post-exploitation command and control framework.

This file is part of Merlin.
Copyright (C) 2026  Russel Van Tuyl

Merlin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as
published by the Free Software Foundation, either version 3 of the License, or any later version.

Merlin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Merlin.  If not, see <http://www.gnu.org/licenses/>.
*/

package info

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestLastCheckin verifies the check-in formatter: a parseable RFC3339 timestamp yields an
// "H:MM:SS ago" string, and anything unparseable (including an empty value) yields "unknown"
// rather than a bogus duration measured from the zero time.
func TestLastCheckin(t *testing.T) {
	got := lastCheckin(time.Now().Add(-90 * time.Minute).Format(time.RFC3339))
	if !strings.HasSuffix(got, " ago") || !strings.HasPrefix(got, "1:") {
		t.Errorf("lastCheckin(90m ago) = %q, want like \"1:30:00 ago\"", got)
	}
	for _, in := range []string{"", "not-a-timestamp", "2026-13-99T99:99:99Z"} {
		if got := lastCheckin(in); got != "unknown" {
			t.Errorf("lastCheckin(%q) = %q, want \"unknown\"", in, got)
		}
	}
}

// DoListeners takes the listener ID from its arguments, not the menu-context id passed to the
// method (which is intentionally unused). These tests exercise the argument validation + listener-ID
// parsing paths, which return before any RPC call to the server.

func TestDoListeners_ArgValidation(t *testing.T) {
	c := NewCommand()
	tests := []struct {
		name    string
		args    string
		substr  string
		wantErr bool
	}{
		{"too few arguments", "info", "requires at least one argument", false},
		{"invalid listener id", "info not-a-uuid", "error parsing UUID", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := c.DoListeners(uuid.Nil, tc.args)
			if resp.Message == nil {
				t.Fatal("response.Message is nil")
			}
			if !strings.Contains(resp.Message.Message(), tc.substr) {
				t.Errorf("message = %q, want substring %q", resp.Message.Message(), tc.substr)
			}
			if resp.Message.Error() != tc.wantErr {
				t.Errorf("Error() = %v, want %v", resp.Message.Error(), tc.wantErr)
			}
		})
	}
}
