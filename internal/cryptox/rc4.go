// Copyright 2023 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cryptox

import (
	"crypto/rc4"
	"github.com/googlecloudplatform/gcsfuse/internal/logger"
)

const bufSize = 64

func RC4Stream(key []byte, off int64, dst, src []byte) {
	var tot int64

	ciph, err := rc4.NewCipher(key)
	if err != nil {
		logger.Fatal("can't create RC4 cipher: %v", err)
	}

	// produce the stream up until the point of `off`
	ciphertext := make([]byte, bufSize)
	zeros := make([]byte, bufSize)

	for tot < off {
		left := off - tot
		n := min(left, bufSize)
		ciph.XORKeyStream(ciphertext, zeros[:n])
		tot += n
	}

	ciph.XORKeyStream(dst, src)
}
