// +build ignore

package checkoutagreements

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
	// CheckoutOptionsEnableAgreements => Enable Terms and Conditions.
	// Path: checkout/options/enable_agreements
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CheckoutOptionsEnableAgreements model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutOptionsEnableAgreements = model.NewBool(`checkout/options/enable_agreements`, model.WithPkgCfg(pkgCfg))

	return pp
}
