// Copyright 2015 CoreStore Authors
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

package tools

type (
	// TableToStructMap uses a string key as easy identifier for maybe later manipulation in config_user.go
	// and a pointer to a TableToStruct struct. Developers can later use the init() func in config_user.go to change
	// the value of variable ConfigTableToStruct.
	TableToStructMap map[string]*TableToStruct
	TableToStruct    struct {
		// Package defines the name of the target package
		Package string
		// OutputFile specifies the path where to write the newly generated code
		OutputFile string
		// QueryString SQL query to filter all the tables which you desire, e.g. SHOW TABLES LIKE 'catalog\_%'
		// @todo maybe change that to []byte
		QueryString string
		// EntityTypeCodes If provided then eav_entity_type.value_table_prefix will be evaluated for further tables.
		EntityTypeCodes []string
	}

	// EntityTypeMap uses a string key as easy identifier, which must also exists in the table eav_EntityTable,
	// for maybe later manipulation in config_user.go
	// and as value a pointer to a EntityType struct. Developers can later use the init() func in config_user.go to change
	// the value of variable ConfigEntityType.
	// The values of EntityType will be uses for materialization in Go code of the eav_entity_type table data.
	EntityTypeMap map[string]*EntityType
	// EntityType is configuration struct which maps the PHP classes to Go types, interfaces and table names.
	EntityType struct {
		// ImportPath path to the package
		ImportPath string
		// EntityModel Go type which implements eav.EntityTypeModeller
		EntityModel string
		// AttributeModel Go type which implements eav.EntityTypeAttributeModeller
		AttributeModel string
		// EntityTable Go type which implements eav.EntityTypeTabler
		EntityTable string
		// IncrementModel Go type which implements eav.EntityTypeIncrementModeller
		IncrementModel string
		// AdditionalAttributeTable Go type which implements eav.EntityTypeAdditionalAttributeTabler
		AdditionalAttributeTable string
		// EntityAttributeCollection Go type which implements eav.EntityAttributeCollectioner
		EntityAttributeCollection string

		// TempAdditionalAttributeTable string which defines the existing table name
		// and specifies more attribute configuration options besides eav_attribute table.
		// This table name is used in a DB query while materializing attribute configuration to Go code.
		// Mage_Eav_Model_Resource_Attribute_Collection::_initSelect()
		TempAdditionalAttributeTable string

		// TempAdditionalAttributeTableWebsite string which defines the existing table name
		// and stores website-dependent attribute parameters.
		// If an EAV model doesn't demand this functionality, let this string empty.
		// This table name is used in a DB query while materializing attribute configuration to Go code.
		// Mage_Customer_Model_Resource_Attribute::_getEavWebsiteTable()
		// Mage_Eav_Model_Resource_Attribute_Collection::_getEavWebsiteTable()
		TempAdditionalAttributeTableWebsite string
	}

	// AttributeModelDefMap contains data to map the three eav_attribute columns
	// (backend|frontend|source)_model to the correct Go function and package.
	// It contains mappings for Magento 1 & 2. A developer has the option to to change/extend the value
	// using the file config_user.go with the init() func.
	// Rethink the Go code here ... because catalog.Product().Attribute().Frontend().Image() is pretty long ... BUT
	// developers coming from Magento are already familiar with this code base and naming ...
	// Def for Definition to avoid a naming conflict :-( Better name?
	AttributeModelDefMap map[string]*AttributeModelDef
	// AttributeModelDef defines which Go type/func has which import path
	AttributeModelDef struct {
		// Is used when you want to specify your own package
		ImportPath string
		// GoModel is a function string which implements, when later executed, one of the n interfaces
		// for (backend|frontend|source|data)_model
		GoModel string
	}
)

// Keys returns all keys from a EntityTypeMap
func (m EntityTypeMap) Keys() []string {
	ret := make([]string, len(m), len(m))
	i := 0
	for k, _ := range m {
		ret[i] = k
		i++
	}
	return ret
}

// EavAttributeColumnNameToInterface mapping column name to Go interface name. Do not add attribute_model
// as this column is unused in Magento 1+2
var EavAttributeColumnNameToInterface = map[string]string{
	"backend_model":  "eav.AttributeBackendModeller",
	"frontend_model": "eav.AttributeFrontendModeller",
	"source_model":   "eav.AttributeSourceModeller",
}

// TablePrefix defines the global table name prefix. See Magento install tool. Can be overridden via func init()
var TablePrefix string = ""

// TableMapMagento1To2 provides mapping between table names. If a table name is in the map then also the struct
// will be rewritten to that new Magneto2 compatible table name.
// Magento2 dev/tools/Magento/Tools/Migration/factory_table_names/replace_ce.php
// @todo implement
var TableMapMagento1To2 = map[string]string{
	"core_cache":             "cache",     // not needed but added in case
	"core_cache_tag":         "cache_tag", // not needed but added in case
	"core_config_data":       "core_config_data",
	"core_design_change":     "design_change",
	"core_directory_storage": "media_storage_directory_storage",
	"core_email_template":    "email_template", // not needed but added in case
	"core_file_storage":      "media_storage_file_storage",
	"core_flag":              "flag",          // not needed but added in case
	"core_layout_link":       "layout_link",   // not needed but added in case
	"core_layout_update":     "layout_update", // not needed but added in case
	"core_resource":          "setup_module",  // not needed but added in case
	"core_session":           "session",       // not needed but added in case
	"core_store":             "store",
	"core_store_group":       "store_group",
	"core_variable":          "variable",
	"core_variable_value":    "variable_value",
	"core_website":           "store_website",
}

// ConfigTableToStruct contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigTableToStruct = TableToStructMap{
	"eav": &TableToStruct{
		Package:         "eav",
		OutputFile:      "eav/generated_tables.go",
		QueryString:     `SHOW TABLES LIKE "{{tableprefix}}eav%"`,
		EntityTypeCodes: nil,
	},
	"store": &TableToStruct{
		Package:    "store",
		OutputFile: "store/generated_tables.go",
		QueryString: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
		    TABLE_SCHEMA = DATABASE() AND
		    TABLE_NAME IN (
		    	'{{tableprefix}}core_store','{{tableprefix}}store',
		    	'{{tableprefix}}core_store_group','{{tableprefix}}store_group',
		    	'{{tableprefix}}core_website','{{tableprefix}}website'
		    ) GROUP BY TABLE_NAME;`,
		EntityTypeCodes: nil,
	},
	"catalog": &TableToStruct{
		Package:    "catalog",
		OutputFile: "catalog/generated_tables.go",
		QueryString: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
		    TABLE_SCHEMA = DATABASE() AND
		    (TABLE_NAME LIKE '{{tableprefix}}catalog\_%' OR TABLE_NAME LIKE '{{tableprefix}}catalogindex%' ) AND
		    TABLE_NAME NOT LIKE '{{tableprefix}}%bundle%' AND
		    TABLE_NAME NOT LIKE '{{tableprefix}}%\_flat\_%' GROUP BY TABLE_NAME;`,
		EntityTypeCodes: []string{"catalog_category", "catalog_product"},
	},
	"customer": &TableToStruct{
		Package:         "customer",
		OutputFile:      "customer/generated_tables.go",
		QueryString:     `SHOW TABLES LIKE "{{tableprefix}}customer%"`,
		EntityTypeCodes: []string{"customer", "customer_address"},
	},
}

// ConfigEntityTypeMaterialization configuration for materializeEntityType() to write the materialized entity types
// into a folder. Other fields of the struct TableToStruct are ignored. Use the file config_user.go with the
// func init() to change/extend it.
var ConfigEntityTypeMaterialization = &TableToStruct{
	Package:    "materialized",
	OutputFile: "materialized/generated_eav_entity.go",
}

// ConfigEntityType contains default configuration. Use the file config_user.go with the func init() to change/extend it.
// Needed in materializeEntityType()
var ConfigEntityType = EntityTypeMap{
	"customer": &EntityType{
		ImportPath:                          "github.com/corestoreio/csfw/customer",
		EntityModel:                         "customer.Customer()",
		AttributeModel:                      "customer.Attribute()",
		EntityTable:                         "customer.Customer()",
		IncrementModel:                      "customer.Customer()",
		AdditionalAttributeTable:            "customer.Customer()",
		EntityAttributeCollection:           "customer.Customer()",
		TempAdditionalAttributeTable:        "{{tableprefix}}customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: "{{tableprefix}}customer_eav_attribute_website",
	},
	"customer_address": &EntityType{
		ImportPath:                          "github.com/corestoreio/csfw/customer",
		EntityModel:                         "customer.Address()",
		AttributeModel:                      "customer.AddressAttribute()",
		EntityTable:                         "customer.Address()",
		AdditionalAttributeTable:            "customer.Address()",
		EntityAttributeCollection:           "customer.Address()",
		TempAdditionalAttributeTable:        "{{tableprefix}}customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: "{{tableprefix}}customer_eav_attribute_website",
	},
	"catalog_category": &EntityType{
		ImportPath:                   "github.com/corestoreio/csfw/catalog",
		EntityModel:                  "catalog.Category()",
		AttributeModel:               "catalog.Attribute()",
		EntityTable:                  "catalog.Category()",
		AdditionalAttributeTable:     "catalog.Category()",
		EntityAttributeCollection:    "catalog.Category()",
		TempAdditionalAttributeTable: "{{tableprefix}}catalog_eav_attribute",
	},
	"catalog_product": &EntityType{
		ImportPath:                   "github.com/corestoreio/csfw/catalog",
		EntityModel:                  "catalog.Product()",
		AttributeModel:               "catalog.Attribute()",
		EntityTable:                  "catalog.Product()",
		AdditionalAttributeTable:     "catalog.Product()",
		EntityAttributeCollection:    "catalog.Product()",
		TempAdditionalAttributeTable: "{{tableprefix}}catalog_eav_attribute",
	},
	// @todo extend for all sales entities
}

// ConfigAttributeModel contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigAttributeModel = AttributeModelDefMap{
	"catalog/product_attribute_frontend_image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Frontend().Image()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Frontend\\Image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Frontend().Image()",
	},
	"eav/entity_attribute_frontend_datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Entity().Attribute().Frontend().Datetime()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Frontend\\Datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Entity().Attribute().Frontend().Datetime()",
	},
	"catalog/attribute_backend_customlayoutupdate": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Attribute().Backend().Customlayoutupdate()",
	},
	"Magento\\Catalog\\Model\\Attribute\\Backend\\Customlayoutupdate": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Attribute().Backend().Customlayoutupdate()",
	},
	"catalog/category_attribute_backend_image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Image()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Image()",
	},
	"catalog/category_attribute_backend_sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Sortby()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Sortby()",
	},
	"catalog/category_attribute_backend_urlkey": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"catalog/product_attribute_backend_boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Boolean()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Boolean()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Category": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Category()",
	},
	"catalog/product_attribute_backend_groupprice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Groupprice()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\GroupPrice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().GroupPrice()",
	},
	"catalog/product_attribute_backend_media": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Media()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Media": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Media()",
	},
	"catalog/product_attribute_backend_msrp": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Price()",
	},
	"catalog/product_attribute_backend_price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Price()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Price()",
	},
	"catalog/product_attribute_backend_recurring": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Recurring()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Recurring": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Recurring()",
	},
	"catalog/product_attribute_backend_sku": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Sku()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Sku": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Sku()",
	},
	"catalog/product_attribute_backend_startdate": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Startdate()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Startdate": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Startdate()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Stock": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Stock()",
	},
	"catalog/product_attribute_backend_tierprice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Tierprice()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Tierprice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Tierprice()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Weight": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Weight()",
	},
	"catalog/product_attribute_backend_urlkey": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"customer/attribute_backend_data_boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Data().Boolean()",
	},
	"Magento\\Customer\\Model\\Attribute\\Backend\\Data\\Boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Data().Boolean()",
	},
	"customer/customer_attribute_backend_billing": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Billing()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Billing": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Billing()",
	},
	"customer/customer_attribute_backend_password": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Password()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Password": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Password()",
	},
	"customer/customer_attribute_backend_shipping": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Shipping()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Shipping": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Shipping()",
	},
	"customer/customer_attribute_backend_store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Store()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Store()",
	},
	"customer/customer_attribute_backend_website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Website()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Website()",
	},
	"customer/entity_address_attribute_backend_region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Backend().Region()",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Backend\\Region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Backend().Region()",
	},
	"customer/entity_address_attribute_backend_street": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Backend().Street()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\DefaultBackend": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().DefaultBackend()",
	},
	"eav/entity_attribute_backend_datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Datetime()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Datetime()",
	},
	"eav/entity_attribute_backend_time_created": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Entity().Attribute().Backend().Time().Created()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Created": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Time().Created()",
	},
	"eav/entity_attribute_backend_time_updated": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Time().Updated()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Updated": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Time().Updated()",
	},
	"bundle/product_attribute_source_price_view": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/bundle",
		GoModel:    "bundle.Product().Attribute().Source().Price().View()",
	},
	"Magento\\Bundle\\Model\\Product\\Attribute\\Source\\Price\\View": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/bundle",
		GoModel:    "bundle.Product().Attribute().Source().Price().View()",
	},
	"Magento\\CatalogInventory\\Model\\Source\\Stock": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/cataloginventory",
		GoModel:    "cataloginventory.Source().Stock()",
	},
	"catalog/category_attribute_source_layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"catalog/category_attribute_source_mode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Mode()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Mode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Mode()",
	},
	"catalog/category_attribute_source_page": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Page()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Page": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Page()",
	},
	"catalog/category_attribute_source_sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Sortby()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Sortby()",
	},
	"catalog/entity_product_attribute_design_options_container": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Design().Options().Container()",
	},
	"Magento\\Catalog\\Model\\Entity\\Product\\Attribute\\Design\\Options\\Container": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Design().Options().Container()",
	},
	"catalog/product_attribute_source_countryofmanufacture": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Countryofmanufacture()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Countryofmanufacture": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Countryofmanufacture()",
	},
	"catalog/product_attribute_source_layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Layout()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Layout()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Status": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Status()",
	},
	"catalog/product_attribute_source_msrp_type_enabled": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Msrp().Type().Enabled()",
	},
	"catalog/product_attribute_source_msrp_type_price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Msrp().Type().Price()",
	},
	"catalog/product_status": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Status()",
	},
	"catalog/product_visibility": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Visibility()",
	},
	"Magento\\Catalog\\Model\\Product\\Visibility": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Visibility()",
	},
	"core/design_source_design": &AttributeModelDef{
		ImportPath: "",
		GoModel:    "",
	},
	"customer/customer_attribute_source_group": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Group()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Group": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Group()",
	},
	"customer/customer_attribute_source_store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Store()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Store()",
	},
	"customer/customer_attribute_source_website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Website()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Website()",
	},
	"customer/entity_address_attribute_source_country": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Country()",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Country": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Country()",
	},
	"customer/entity_address_attribute_source_region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Region()",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Region()",
	},
	"eav/entity_attribute_source_boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Boolean()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Boolean()",
	},
	"eav/entity_attribute_source_table": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Table()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Table": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Table()",
	},
	"Magento\\Msrp\\Model\\Product\\Attribute\\Source\\Type\\Price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/msrp",
		GoModel:    "msrp.Product().Attribute().Source().Type().Price()",
	},
	"tax/class_source_product": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.Class().Source().Product()",
	},
	"Magento\\Tax\\Model\\TaxClass\\Source\\Product": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.TaxClass().Source().Product()",
	},
	"Magento\\Theme\\Model\\Theme\\Source\\Theme": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/theme",
		GoModel:    "",
	},
	"customer/attribute_data_postcode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Data.Postcode()",
	},
	"Magento\\Customer\\Model\\Attribute\\Data\\Postcode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Data.Postcode()",
	},
}
