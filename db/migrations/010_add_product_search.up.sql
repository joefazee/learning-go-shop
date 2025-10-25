-- Add tsvector column for full-text search
ALTER TABLE products ADD COLUMN search_vector tsvector;

-- Create function to update search vector
CREATE OR REPLACE FUNCTION products_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
            setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
            setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B') ||
            setweight(to_tsvector('english', coalesce(NEW.sku, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Create trigger to automatically update search_vector
CREATE TRIGGER products_search_vector_trigger
    BEFORE INSERT OR UPDATE ON products
    FOR EACH ROW
EXECUTE FUNCTION products_search_vector_update();


-- Update existing products with search vectors
UPDATE products SET search_vector =
                        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
                        setweight(to_tsvector('english', coalesce(description, '')), 'B') ||
                        setweight(to_tsvector('english', coalesce(sku, '')), 'C');


CREATE INDEX idx_products_search_vector ON products USING GIN(search_vector);

-- Add comment for documentation
COMMENT ON COLUMN products.search_vector IS
    'Full-text search vector with weighted fields: A=name, B=description, C=sku';
