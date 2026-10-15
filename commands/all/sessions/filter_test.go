/*
Merlin is a post-exploitation command and control framework.

This file is part of Merlin.
Copyright (C) 2026  Nerijus Bendžiūnas

Merlin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as
published by the Free Software Foundation, either version 3 of the License, or any later version.

Merlin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Merlin.  If not, see <http://www.gnu.org/licenses/>.
*/

package sessions

import (
	"slices"
	"strings"
	"testing"
)

var testHeader = []string{"Agent GUID", "Transport", "Platform", "Host", "User", "Process", "Status", "Last Checkin", "Note"}

var testRows = [][]string{
	{"agent-web-active", "https", "linux/amd64", "web01", "root", "worker(100)", "Active", "0:00:05 ago", ""},
	{"agent-ws4-active", "https", "linux/amd64", "ws4", "root", "worker(200)", "Active", "0:00:07 ago", ""},
	{"agent-ws4-delayed", "https", "linux/amd64", "ws4", "svc", "worker(300)", "Delayed", "0:05:00 ago", ""},
	{"agent-dc-dead", "h2", "windows/amd64", "DC01", "SYSTEM", "svchost(400)", "Dead", "3:00:00 ago", ""},
	{"agent-init", "", "/", "", "", "(0)", "Init", "2562047:47:16 ago", ""},
}

func agentIDs(rows [][]string) []string {
	var ids []string
	for _, row := range rows {
		ids = append(ids, row[0])
	}
	return ids
}

func TestFilterRows(t *testing.T) {
	tests := []struct {
		name    string
		filters string
		want    []string
	}{
		{"no filter lists every session", "", agentIDs(testRows)},
		{"active", "active", []string{"agent-web-active", "agent-ws4-active"}},
		{"init", "init", []string{"agent-init"}},
		{"dead", "dead", []string{"agent-dc-dead"}},
		{"several statuses match any", "active delayed", []string{"agent-web-active", "agent-ws4-active", "agent-ws4-delayed"}},
		{"status ignores case", "ACTIVE", []string{"agent-web-active", "agent-ws4-active"}},
		{"host matches part of the name", "host ws", []string{"agent-ws4-active", "agent-ws4-delayed"}},
		{"host ignores case", "host dc01", []string{"agent-dc-dead"}},
		{"user", "user system", []string{"agent-dc-dead"}},
		{"platform", "platform windows", []string{"agent-dc-dead"}},
		{"different filters must all match", "active host ws4", []string{"agent-ws4-active"}},
		{"all four filters", "delayed host ws4 user svc platform linux", []string{"agent-ws4-delayed"}},
		{"nothing matches", "dead host web01", nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter, err := parseFilter(strings.Fields(test.filters))
			if err != nil {
				t.Fatalf("parseFilter(%q) returned error: %s", test.filters, err)
			}
			rows, err := filterRows(testHeader, testRows, filter)
			if err != nil {
				t.Fatalf("filterRows returned error: %s", err)
			}
			if got := agentIDs(rows); !slices.Equal(got, test.want) {
				t.Errorf("filters %q listed %v, want %v", test.filters, got, test.want)
			}
		})
	}
}

func TestParseFilterRejectsBadInput(t *testing.T) {
	tests := []struct {
		name    string
		filters string
		message string
	}{
		{"unknown word", "alive", "unknown sessions filter 'alive'"},
		{"host without value", "active host", "'host' filter requires a value"},
		{"user without value", "user", "'user' filter requires a value"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseFilter(strings.Fields(test.filters))
			if err == nil {
				t.Fatalf("parseFilter(%q) returned no error", test.filters)
			}
			if err.Error() != test.message {
				t.Errorf("parseFilter(%q) error %q, want %q", test.filters, err, test.message)
			}
		})
	}
}

func TestFilterRowsRejectsMissingColumn(t *testing.T) {
	header := []string{"Agent GUID", "Host"}
	rows := [][]string{{"agent-web-active", "web01"}}
	filter, err := parseFilter([]string{"active"})
	if err != nil {
		t.Fatalf("parseFilter returned error: %s", err)
	}
	if _, err := filterRows(header, rows, filter); err == nil {
		t.Fatal("filterRows returned no error for a table without a Status column")
	}
}
