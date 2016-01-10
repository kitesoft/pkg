// +build ignore

package cron

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "cron",
					Label:     `Cron (Scheduled Tasks) - all the times are in minutes`,
					Comment:   element.LongText(`For correct URLs generated during cron runs please make sure that Web > Secure and Unsecure Base URLs are explicitly set.`),
					SortOrder: 15,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields:    element.NewFieldSlice(),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}