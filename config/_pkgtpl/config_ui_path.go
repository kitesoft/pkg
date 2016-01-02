// +build ignore

package ui

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// DevJsSessionStorageLogging => Log JS Errors to Session Storage.
	// If enabled, can be used by functional tests for extended reporting
	// Path: dev/js/session_storage_logging
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevJsSessionStorageLogging model.Bool

	// DevJsSessionStorageKey => Log JS Errors to Session Storage Key.
	// Use this key to retrieve collected js errors
	// Path: dev/js/session_storage_key
	DevJsSessionStorageKey model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.DevJsSessionStorageLogging = model.NewBool(`dev/js/session_storage_logging`, model.WithPkgCfg(pkgCfg))
	pp.DevJsSessionStorageKey = model.NewStr(`dev/js/session_storage_key`, model.WithPkgCfg(pkgCfg))

	return pp
}
