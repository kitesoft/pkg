// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package util_test

import (
	"testing"

	"github.com/corestoreio/csfw/util"
	"github.com/stretchr/testify/assert"
)

var magento1TableList = []string{"admin_assert", "admin_role", "admin_rule", "admin_user", "adminnotification_inbox", "api2_acl_attribute", "api2_acl_role", "api2_acl_rule", "api2_acl_user", "api_assert", "api_role", "api_rule", "api_session", "api_user", "captcha_log", "catalog_category_anc_categs_index_idx", "catalog_category_anc_categs_index_tmp", "catalog_category_anc_products_index_idx", "catalog_category_anc_products_index_tmp", "catalog_category_entity", "catalog_category_entity_datetime", "catalog_category_entity_decimal", "catalog_category_entity_int", "catalog_category_entity_text", "catalog_category_entity_varchar", "catalog_category_flat_store_1", "catalog_category_product", "catalog_category_product_index", "catalog_category_product_index_enbl_idx", "catalog_category_product_index_enbl_tmp", "catalog_category_product_index_idx", "catalog_category_product_index_tmp", "catalog_compare_item", "catalog_eav_attribute", "catalog_product_bundle_option", "catalog_product_bundle_option_value", "catalog_product_bundle_price_index", "catalog_product_bundle_selection", "catalog_product_bundle_selection_price", "catalog_product_bundle_stock_index", "catalog_product_enabled_index", "catalog_product_entity", "catalog_product_entity_datetime", "catalog_product_entity_decimal", "catalog_product_entity_gallery", "catalog_product_entity_group_price", "catalog_product_entity_int", "catalog_product_entity_media_gallery", "catalog_product_entity_media_gallery_value", "catalog_product_entity_text", "catalog_product_entity_tier_price", "catalog_product_entity_varchar", "catalog_product_flat_1", "catalog_product_index_eav", "catalog_product_index_eav_decimal", "catalog_product_index_eav_decimal_idx", "catalog_product_index_eav_decimal_tmp", "catalog_product_index_eav_idx", "catalog_product_index_eav_tmp", "catalog_product_index_group_price", "catalog_product_index_price", "catalog_product_index_price_bundle_idx", "catalog_product_index_price_bundle_opt_idx", "catalog_product_index_price_bundle_opt_tmp", "catalog_product_index_price_bundle_sel_idx", "catalog_product_index_price_bundle_sel_tmp", "catalog_product_index_price_bundle_tmp", "catalog_product_index_price_cfg_opt_agr_idx", "catalog_product_index_price_cfg_opt_agr_tmp", "catalog_product_index_price_cfg_opt_idx", "catalog_product_index_price_cfg_opt_tmp", "catalog_product_index_price_downlod_idx", "catalog_product_index_price_downlod_tmp", "catalog_product_index_price_final_idx", "catalog_product_index_price_final_tmp", "catalog_product_index_price_idx", "catalog_product_index_price_opt_agr_idx", "catalog_product_index_price_opt_agr_tmp", "catalog_product_index_price_opt_idx", "catalog_product_index_price_opt_tmp", "catalog_product_index_price_tmp", "catalog_product_index_tier_price", "catalog_product_index_website", "catalog_product_link", "catalog_product_link_attribute", "catalog_product_link_attribute_decimal", "catalog_product_link_attribute_int", "catalog_product_link_attribute_varchar", "catalog_product_link_type", "catalog_product_option", "catalog_product_option_price", "catalog_product_option_title", "catalog_product_option_type_price", "catalog_product_option_type_title", "catalog_product_option_type_value", "catalog_product_relation", "catalog_product_super_attribute", "catalog_product_super_attribute_label", "catalog_product_super_attribute_pricing", "catalog_product_super_link", "catalog_product_website", "cataloginventory_stock", "cataloginventory_stock_item", "cataloginventory_stock_status", "cataloginventory_stock_status_idx", "cataloginventory_stock_status_tmp", "catalogrule", "catalogrule_affected_product", "catalogrule_customer_group", "catalogrule_group_website", "catalogrule_product", "catalogrule_product_price", "catalogrule_website", "catalogsearch_fulltext", "catalogsearch_query", "catalogsearch_result", "checkout_agreement", "checkout_agreement_store", "cms_block", "cms_block_store", "cms_page", "cms_page_store", "core_cache", "core_cache_option", "core_cache_tag", "core_config_data", "core_email_queue", "core_email_queue_recipients", "core_email_template", "core_flag", "core_layout_link", "core_layout_update", "core_resource", "core_session", "core_store", "core_store_group", "core_translate", "core_url_rewrite", "core_variable", "core_variable_value", "core_website", "coupon_aggregated", "coupon_aggregated_order", "coupon_aggregated_updated", "cron_schedule", "customer_address_entity", "customer_address_entity_datetime", "customer_address_entity_decimal", "customer_address_entity_int", "customer_address_entity_text", "customer_address_entity_varchar", "customer_eav_attribute", "customer_eav_attribute_website", "customer_entity", "customer_entity_datetime", "customer_entity_decimal", "customer_entity_int", "customer_entity_text", "customer_entity_varchar", "customer_form_attribute", "customer_group", "dataflow_batch", "dataflow_batch_export", "dataflow_batch_import", "dataflow_import_data", "dataflow_profile", "dataflow_profile_history", "dataflow_session", "dbr_people", "design_change", "directory_country", "directory_country_format", "directory_country_region", "directory_country_region_name", "directory_currency_rate", "downloadable_link", "downloadable_link_price", "downloadable_link_purchased", "downloadable_link_purchased_item", "downloadable_link_title", "downloadable_sample", "downloadable_sample_title", "eav_attribute", "eav_attribute_group", "eav_attribute_label", "eav_attribute_option", "eav_attribute_option_value", "eav_attribute_set", "eav_entity", "eav_entity_attribute", "eav_entity_datetime", "eav_entity_decimal", "eav_entity_int", "eav_entity_store", "eav_entity_text", "eav_entity_type", "eav_entity_varchar", "eav_form_element", "eav_form_fieldset", "eav_form_fieldset_label", "eav_form_type", "eav_form_type_entity", "gift_message", "importexport_importdata", "index_event", "index_process", "index_process_event", "log_customer", "log_quote", "log_summary", "log_summary_type", "log_url", "log_url_info", "log_visitor", "log_visitor_info", "log_visitor_online", "newsletter_problem", "newsletter_queue", "newsletter_queue_link", "newsletter_queue_store_link", "newsletter_subscriber", "newsletter_template", "null_types", "oauth_consumer", "oauth_nonce", "oauth_token", "paypal_cert", "paypal_payment_transaction", "paypal_settlement_report", "paypal_settlement_report_row", "persistent_session", "poll", "poll_answer", "poll_store", "poll_vote", "product_alert_price", "product_alert_stock", "rating", "rating_entity", "rating_option", "rating_option_vote", "rating_option_vote_aggregated", "rating_store", "rating_title", "report_compared_product_index", "report_event", "report_event_types", "report_viewed_product_aggregated_daily", "report_viewed_product_aggregated_monthly", "report_viewed_product_aggregated_yearly", "report_viewed_product_index", "review", "review_detail", "review_entity", "review_entity_summary", "review_status", "review_store", "sales_bestsellers_aggregated_daily", "sales_bestsellers_aggregated_monthly", "sales_bestsellers_aggregated_yearly", "sales_billing_agreement", "sales_billing_agreement_order", "sales_flat_creditmemo", "sales_flat_creditmemo_comment", "sales_flat_creditmemo_grid", "sales_flat_creditmemo_item", "sales_flat_invoice", "sales_flat_invoice_comment", "sales_flat_invoice_grid", "sales_flat_invoice_item", "sales_flat_order", "sales_flat_order_address", "sales_flat_order_grid", "sales_flat_order_item", "sales_flat_order_payment", "sales_flat_order_status_history", "sales_flat_quote", "sales_flat_quote_address", "sales_flat_quote_address_item", "sales_flat_quote_item", "sales_flat_quote_item_option", "sales_flat_quote_payment", "sales_flat_quote_shipping_rate", "sales_flat_shipment", "sales_flat_shipment_comment", "sales_flat_shipment_grid", "sales_flat_shipment_item", "sales_flat_shipment_track", "sales_invoiced_aggregated", "sales_invoiced_aggregated_order", "sales_order_aggregated_created", "sales_order_aggregated_updated", "sales_order_status", "sales_order_status_label", "sales_order_status_state", "sales_order_tax", "sales_order_tax_item", "sales_payment_transaction", "sales_recurring_profile", "sales_recurring_profile_order", "sales_refunded_aggregated", "sales_refunded_aggregated_order", "sales_shipping_aggregated", "sales_shipping_aggregated_order", "salesrule", "salesrule_coupon", "salesrule_coupon_usage", "salesrule_customer", "salesrule_customer_group", "salesrule_label", "salesrule_product_attribute", "salesrule_website", "sendfriend_log", "shipping_tablerate", "sitemap", "tag", "tag_properties", "tag_relation", "tag_summary", "tax_calculation", "tax_calculation_rate", "tax_calculation_rate_title", "tax_calculation_rule", "tax_class", "tax_order_aggregated_created", "tax_order_aggregated_updated", "weee_discount", "weee_tax", "widget", "widget_instance", "widget_instance_page", "widget_instance_page_layout", "wishlist", "wishlist_item", "wishlist_item_option"}
var magento2TableList = []string{"admin_system_messages", "admin_user", "adminnotification_inbox", "authorization_role", "authorization_rule", "cache", "cache_tag", "captcha_log", "catalog_category_entity", "catalog_category_entity_datetime", "catalog_category_entity_decimal", "catalog_category_entity_int", "catalog_category_entity_text", "catalog_category_entity_varchar", "catalog_category_product", "catalog_category_product_index", "catalog_category_product_index_tmp", "catalog_compare_item", "catalog_eav_attribute", "catalog_product_bundle_option", "catalog_product_bundle_option_value", "catalog_product_bundle_price_index", "catalog_product_bundle_selection", "catalog_product_bundle_selection_price", "catalog_product_bundle_stock_index", "catalog_product_entity", "catalog_product_entity_datetime", "catalog_product_entity_decimal", "catalog_product_entity_gallery", "catalog_product_entity_group_price", "catalog_product_entity_int", "catalog_product_entity_media_gallery", "catalog_product_entity_media_gallery_value", "catalog_product_entity_text", "catalog_product_entity_tier_price", "catalog_product_entity_varchar", "catalog_product_index_eav", "catalog_product_index_eav_decimal", "catalog_product_index_eav_decimal_idx", "catalog_product_index_eav_decimal_tmp", "catalog_product_index_eav_idx", "catalog_product_index_eav_tmp", "catalog_product_index_group_price", "catalog_product_index_price", "catalog_product_index_price_bundle_idx", "catalog_product_index_price_bundle_opt_idx", "catalog_product_index_price_bundle_opt_tmp", "catalog_product_index_price_bundle_sel_idx", "catalog_product_index_price_bundle_sel_tmp", "catalog_product_index_price_bundle_tmp", "catalog_product_index_price_cfg_opt_agr_idx", "catalog_product_index_price_cfg_opt_agr_tmp", "catalog_product_index_price_cfg_opt_idx", "catalog_product_index_price_cfg_opt_tmp", "catalog_product_index_price_downlod_idx", "catalog_product_index_price_downlod_tmp", "catalog_product_index_price_final_idx", "catalog_product_index_price_final_tmp", "catalog_product_index_price_idx", "catalog_product_index_price_opt_agr_idx", "catalog_product_index_price_opt_agr_tmp", "catalog_product_index_price_opt_idx", "catalog_product_index_price_opt_tmp", "catalog_product_index_price_tmp", "catalog_product_index_tier_price", "catalog_product_index_website", "catalog_product_link", "catalog_product_link_attribute", "catalog_product_link_attribute_decimal", "catalog_product_link_attribute_int", "catalog_product_link_attribute_varchar", "catalog_product_link_type", "catalog_product_option", "catalog_product_option_price", "catalog_product_option_title", "catalog_product_option_type_price", "catalog_product_option_type_title", "catalog_product_option_type_value", "catalog_product_relation", "catalog_product_super_attribute", "catalog_product_super_attribute_label", "catalog_product_super_link", "catalog_product_website", "catalog_url_rewrite_product_category", "cataloginventory_stock", "cataloginventory_stock_item", "cataloginventory_stock_status", "cataloginventory_stock_status_idx", "cataloginventory_stock_status_tmp", "catalogrule", "catalogrule_customer_group", "catalogrule_group_website", "catalogrule_product", "catalogrule_product_price", "catalogrule_website", "catalogsearch_fulltext_scope1", "checkout_agreement", "checkout_agreement_store", "cms_block", "cms_block_store", "cms_page", "cms_page_store", "core_config_data", "cron_schedule", "csCustomer_value_datetime", "csCustomer_value_decimal", "csCustomer_value_int", "csCustomer_value_text", "csCustomer_value_varchar", "customer_address_entity", "customer_address_entity_datetime", "customer_address_entity_decimal", "customer_address_entity_int", "customer_address_entity_text", "customer_address_entity_varchar", "customer_eav_attribute", "customer_eav_attribute_website", "customer_entity", "customer_form_attribute", "customer_group", "customer_log", "customer_visitor", "design_change", "directory_country", "directory_country_format", "directory_country_region", "directory_country_region_name", "directory_currency_rate", "downloadable_link", "downloadable_link_price", "downloadable_link_purchased", "downloadable_link_purchased_item", "downloadable_link_title", "downloadable_sample", "downloadable_sample_title", "eav_attribute", "eav_attribute_group", "eav_attribute_label", "eav_attribute_option", "eav_attribute_option_value", "eav_attribute_set", "eav_entity", "eav_entity_attribute", "eav_entity_datetime", "eav_entity_decimal", "eav_entity_int", "eav_entity_store", "eav_entity_text", "eav_entity_type", "eav_entity_varchar", "eav_form_element", "eav_form_fieldset", "eav_form_fieldset_label", "eav_form_type", "eav_form_type_entity", "email_template", "flag", "gift_message", "googleoptimizer_code", "import_history", "importexport_importdata", "indexer_state", "integration", "layout_link", "layout_update", "log_customer", "log_quote", "log_summary", "log_summary_type", "log_url", "log_url_info", "log_visitor", "log_visitor_info", "log_visitor_online", "mview_state", "newsletter_problem", "newsletter_queue", "newsletter_queue_link", "newsletter_queue_store_link", "newsletter_subscriber", "newsletter_template", "oauth_consumer", "oauth_nonce", "oauth_token", "paypal_billing_agreement", "paypal_billing_agreement_order", "paypal_cert", "paypal_payment_transaction", "paypal_settlement_report", "paypal_settlement_report_row", "persistent_session", "product_alert_price", "product_alert_stock", "quote", "quote_address", "quote_address_item", "quote_id_mask", "quote_item", "quote_item_option", "quote_payment", "quote_shipping_rate", "rating", "rating_entity", "rating_option", "rating_option_vote", "rating_option_vote_aggregated", "rating_store", "rating_title", "report_compared_product_index", "report_event", "report_event_types", "report_viewed_product_aggregated_daily", "report_viewed_product_aggregated_monthly", "report_viewed_product_aggregated_yearly", "report_viewed_product_index", "review", "review_detail", "review_entity", "review_entity_summary", "review_status", "review_store", "sales_bestsellers_aggregated_daily", "sales_bestsellers_aggregated_monthly", "sales_bestsellers_aggregated_yearly", "sales_creditmemo", "sales_creditmemo_comment", "sales_creditmemo_grid", "sales_creditmemo_item", "sales_invoice", "sales_invoice_comment", "sales_invoice_grid", "sales_invoice_item", "sales_invoiced_aggregated", "sales_invoiced_aggregated_order", "sales_order", "sales_order_address", "sales_order_aggregated_created", "sales_order_aggregated_updated", "sales_order_grid", "sales_order_item", "sales_order_payment", "sales_order_status", "sales_order_status_history", "sales_order_status_label", "sales_order_status_state", "sales_order_tax", "sales_order_tax_item", "sales_payment_transaction", "sales_refunded_aggregated", "sales_refunded_aggregated_order", "sales_sequence_meta", "sales_sequence_profile", "sales_shipment", "sales_shipment_comment", "sales_shipment_grid", "sales_shipment_item", "sales_shipment_track", "sales_shipping_aggregated", "sales_shipping_aggregated_order", "salesrule", "salesrule_coupon", "salesrule_coupon_aggregated", "salesrule_coupon_aggregated_order", "salesrule_coupon_aggregated_updated", "salesrule_coupon_usage", "salesrule_customer", "salesrule_customer_group", "salesrule_label", "salesrule_product_attribute", "salesrule_website", "search_query", "sendfriend_log", "sequence_creditmemo_0", "sequence_creditmemo_1", "sequence_invoice_0", "sequence_invoice_1", "sequence_order_0", "sequence_order_1", "sequence_shipment_0", "sequence_shipment_1", "session", "setup_module", "shipping_tablerate", "sitemap", "store", "store_group", "store_website", "tax_calculation", "tax_calculation_rate", "tax_calculation_rate_title", "tax_calculation_rule", "tax_class", "tax_order_aggregated_created", "tax_order_aggregated_updated", "theme", "theme_file", "translation", "ui_bookmark", "url_rewrite", "variable", "variable_value", "vde_theme_change", "weee_tax", "widget", "widget_instance", "widget_instance_page", "widget_instance_page_layout", "wishlist", "wishlist_item", "wishlist_item_option"}

func TestMagentoVersion(t *testing.T) {
	tests := []struct {
		want1  bool
		want2  bool
		prefix string
		ts     []string
	}{
		{true, false, "", []string{"core_store", "core_website", "core_store_group", "api_user", "admin_user"}},
		{true, false, "prefix_", []string{"prefix_core_store", "prefix_core_website", "prefix_core_store_group", "prefix_api_user", "prefix_admin_user"}},
		{false, true, "", []string{"integration", "store_website", "store_group", "authorization_role"}},
		{true, false, "", magento1TableList},
		{false, true, "", magento2TableList},
		{false, true, "Prefix_", []string{"Prefix_integration", "Prefix_store_website", "Prefix_store_group", "Prefix_admin_user", "Prefix_authorization_role"}},
	}

	for _, test := range tests {
		v1, v2 := util.MagentoVersion(test.prefix, test.ts)
		assert.Equal(t, test.want1, v1, "%v", test)
		assert.Equal(t, test.want2, v2, "%v", test)
	}
}

var benchmarkMagentoVersion bool

// BenchmarkMagentoVersionTrue-4 	   30000	     51502 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMagentoVersionTrue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkMagentoVersion, _ = util.MagentoVersion("", magento1TableList)
	}
	if benchmarkMagentoVersion != true {
		b.Error("Must be Magento1!")
	}
}

// BenchmarkMagentoVersionFalse-4	   30000	     47737 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMagentoVersionFalse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, benchmarkMagentoVersion = util.MagentoVersion("", magento2TableList)
	}
	if benchmarkMagentoVersion != true {
		b.Error("Cannot be Magento1!")
	}
}