/*
 * Copyright (c) 2019 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package pg_query

import (
	"strings"
)

func (node FuncCall) Name() (string, error) {
	difference := func(slice1 []string, slice2 []string) []string {
		var diff []string
		// Loop two times, first to find slice1 strings not in slice2,
		// second loop to find slice2 strings not in slice1
		for i := 0; i < 2; i++ {
			for _, s1 := range slice1 {
				found := false
				for _, s2 := range slice2 {
					if s1 == s2 {
						found = true
						break
					}
				}
				// String not found. We add it to return slice
				if !found {
					diff = append(diff, s1)
				}
			}
			// Swap the slices, only if it was the first loop
			if i == 0 {
				slice1, slice2 = slice2, slice1
			}
		}
		return diff
	}

	if names, err := node.Funcname.DeparseList(Context_FuncCall); err != nil {
		return "", err
	} else {
		return strings.Join(difference([]string{"pg_catalog"}, names), "."), nil
	}
}
