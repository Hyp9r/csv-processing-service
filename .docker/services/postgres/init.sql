CREATE TABLE products (
	id SERIAL PRIMARY KEY,
	product_name TEXT,
	product_category TEXT,
	product_price NUMERIC,
	product_description TEXT,
	brand_name TEXT,
	stock_quantity INTEGER,
	manufacturer TEXT,
	sku TEXT UNIQUE,
	weight NUMERIC,
	color TEXT
);