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
	"fmt"
	"maps"
	"slices"
	"strings"
)

var sessionStatuses = []string{"active", "delayed", "dead", "init"}

var sessionContainsFields = []string{"host", "user", "platform"}

type sessionFilter struct {
	statuses []string
	contains map[string]string
}

func parseFilter(args []string) (filter sessionFilter, err error) {
	filter.contains = make(map[string]string)
	for index := 0; index < len(args); index++ {
		word := strings.ToLower(args[index])
		switch {
		case slices.Contains(sessionStatuses, word):
			filter.statuses = append(filter.statuses, word)
		case slices.Contains(sessionContainsFields, word):
			if index+1 >= len(args) {
				return filter, fmt.Errorf("'%s' filter requires a value", word)
			}
			index++
			filter.contains[word] = strings.ToLower(args[index])
		default:
			return filter, fmt.Errorf("unknown sessions filter '%s'", args[index])
		}
	}
	return filter, nil
}

func (f sessionFilter) columns() []string {
	names := slices.Sorted(maps.Keys(f.contains))
	if len(f.statuses) > 0 {
		names = append(names, "status")
	}
	return names
}

func (f sessionFilter) matches(row []string, columns map[string]int) bool {
	if len(f.statuses) > 0 && !slices.Contains(f.statuses, cellValue(row, columns["status"])) {
		return false
	}
	for name, value := range f.contains {
		if !strings.Contains(cellValue(row, columns[name]), value) {
			return false
		}
	}
	return true
}

func cellValue(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(row[index]))
}

func filterRows(header []string, rows [][]string, filter sessionFilter) ([][]string, error) {
	columns := make(map[string]int, len(header))
	for index, name := range header {
		columns[strings.ToLower(strings.TrimSpace(name))] = index
	}
	for _, name := range filter.columns() {
		if _, found := columns[name]; !found {
			return nil, fmt.Errorf("sessions table from the server has no '%s' column", name)
		}
	}
	var matched [][]string
	for _, row := range rows {
		if filter.matches(row, columns) {
			matched = append(matched, row)
		}
	}
	return matched, nil
}
