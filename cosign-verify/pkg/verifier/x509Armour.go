//
// Copyright 2021 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package verifier

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/prometheus/common/log"
	"github.com/sigstore/sigstore/pkg/signature"
)

type PublicKey interface {
	signature.Verifier
	signature.PublicKeyProvider
}

func decodeArmoredKey(keyBytes []byte) (PublicKey, error) { // (PublicKey, error) (*ecdsa.PublicKey, error)
	if len(keyBytes) == 0 {
		return nil, fmt.Errorf("Key: empty")
	}
	pems := parsePems(keyBytes)
	for _, p := range pems {
		// TODO check header
		key, err := x509.ParsePKIXPublicKey(p.Bytes)
		if err != nil {
			log.Error(err, "parsing key", "key", p)
		}
		// return key.(*ecdsa.PublicKey), nil
		return signature.ECDSAVerifier{Key: key.(*ecdsa.PublicKey), HashAlg: crypto.SHA256}, nil
	}
	return nil, fmt.Errorf("Key: empty")
}

func parsePems(b []byte) []*pem.Block {
	p, rest := pem.Decode(b)
	if p == nil {
		return nil
	}
	pems := []*pem.Block{p}

	if rest != nil {
		return append(pems, parsePems(rest)...)
	}
	return pems
}
