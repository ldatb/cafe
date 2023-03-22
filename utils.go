// Copyright (C) 2023 Lucas de Ataides <lucasatab@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package cafe

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readFile reads the contents of a file into an array of strings
func readCAFEFile(filepath string) ([]string, error) {
	if !strings.HasSuffix(filepath, ".cafe") {
		err := fmt.Errorf("%s is not a .cafe file", filepath)
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := []string{}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	return data, nil
}

// Checks if a given string matches any of the elements
func equalsToMany(target string, values []string) bool {
	for _, val := range values {
		if strings.TrimSpace(target) == strings.TrimSpace(val) {
			return true
		}
	}
	return false
}

// Checks if a given string has a prefix that matches any of the elements
func hasPrefixToMany(target string, values []string) bool {
	for _, val := range values {
		targetTrimmed := strings.TrimSpace(target)
		valTrimmed := strings.TrimSpace(val)
		if strings.HasPrefix(targetTrimmed, valTrimmed) {
			return true
		}
	}
	return false
}
