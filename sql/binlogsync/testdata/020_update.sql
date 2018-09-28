/*
// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
*/
UPDATE `catalog_product_entity` SET
	`attribute_set_id`=111,
	 `sku`='MH01-XL-CS111',
	 `has_options`=1
WHERE entity_id = 66;

-- this query should delete also entity value rows in their tables and
-- must propagate these deletions also in the binary log.
delete from catalog_product_entity WHERE entity_id = 65;


UPDATE `catalog_product_entity_varchar` SET value='Dodo Sport Watch and See' WHERE `value_id`=436;

UPDATE catalog_product_entity_int set value=null where value_id=212;
