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

package types

import (
	"database/sql/driver"

	"github.com/pkg/errors"
)

func DatabaseSQLValue(ci *ConnInfo, src Value) (interface{}, error) {
	if valuer, ok := src.(driver.Valuer); ok {
		return valuer.Value()
	}

	if textEncoder, ok := src.(TextEncoder); ok {
		buf, err := textEncoder.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		return string(buf), nil
	}

	if binaryEncoder, ok := src.(BinaryEncoder); ok {
		buf, err := binaryEncoder.EncodeBinary(ci, nil)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}

	return nil, errors.New("cannot convert to database/sql compatible value")
}

func EncodeValueText(src TextEncoder) (interface{}, error) {
	buf, err := src.EncodeText(nil, make([]byte, 0, 32))
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return string(buf), err
}
