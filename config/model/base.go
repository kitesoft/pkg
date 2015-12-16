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
	"github.com/corestoreio/csfw/util/cast"
)

// path defines the path in the core_config_data table like a/b/c. All other
// types in this package inherits from this path type
type path string

func (p path) set(w config.Writer, v interface{}, s scope.Scope, id int64) error {
	return w.Write(config.Path(string(p)), config.Value(v), config.Scope(s, id))
}

func (p path) lookupString(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string) {
	aPath := string(p)
	if fields, err := pkgCfg.FindFieldByPath(aPath); err == nil {
		v, _ = cast.ToStringE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.StringSlice.SectionSlice.FindFieldByPath", "err", err, "path", aPath)
		}
	}

	if val, err := sg.String(aPath); err == nil {
		v = val
	}
	return
}
