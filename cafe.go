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

// Convert a CAFE file to a Go struct
func Decode(filename string) (*Parser, error) {
	input, err := readCAFEFile(filename)
	if err != nil {
		return nil, err
	}
	p := newParser(input)
	p.parseItems(false)
	return p, nil
}
