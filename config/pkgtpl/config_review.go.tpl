package review

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "review",
				Label:     `Product Reviews`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/review/allow_guest`,
						ID:           "allow_guest",
						Label:        `Allow Guests to Write Reviews`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
