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

package pgproto

import (
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/readystock/pgx/pgio"
)

const (
	AuthTypeOk                = 0
	AuthTypeCleartextPassword = 3
	AuthTypeMD5Password       = 5
)

type Authentication struct {
	Type uint32

	// MD5Password fields
	Salt [4]byte
}

func (*Authentication) Backend() {}

func (dst *Authentication) Decode(src []byte) error {
	*dst = Authentication{Type: binary.BigEndian.Uint32(src[:4])}

	switch dst.Type {
	case AuthTypeOk:
	case AuthTypeCleartextPassword:
	case AuthTypeMD5Password:
		copy(dst.Salt[:], src[4:8])
	default:
		return errors.Errorf("unknown authentication type: %d", dst.Type)
	}

	return nil
}

func (src *Authentication) Encode(dst []byte) []byte {
	dst = append(dst, 'R')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)
	dst = pgio.AppendUint32(dst, src.Type)

	switch src.Type {
	case AuthTypeMD5Password:
		dst = append(dst, src.Salt[:]...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
