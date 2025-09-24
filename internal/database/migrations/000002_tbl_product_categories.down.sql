ALTER TABLE products
ADD COLUMN category_id BIGINT REFERENCES categories(id) 
ON DELETE SET NULL;

DROP TABLE IF EXISTS product_categories;