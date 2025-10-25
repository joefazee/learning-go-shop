DROP INDEX IF EXISTS idx_products_search_vector;

DROP TRIGGER IF EXISTS products_search_vector_trigger ON products;

DROP FUNCTION IF EXISTS products_search_vector_update();

ALTER TABLE products DROP COLUMN IF EXISTS search_vector;
