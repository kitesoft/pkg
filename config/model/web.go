// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package model

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

// URL represents a path in config.Getter handles URLs and internal validation
type URL path

func (p URL) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string, err error) {
	// todo URL checks
	return path(p).lookupString(pkgCfg, sg), nil
}

func (p URL) Set(w config.Writer, v string, s scope.Scope, id int64) error {
	return path(p).set(w, v, s, id)
}
