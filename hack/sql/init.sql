CREATE TABLE IF NOT EXISTS ingredients_thresholds (
   id SERIAL PRIMARY KEY,
   ingredient_type VARCHAR(10) NOT NULL,
   min_value NUMERIC(4,2) NOT NULL,
   max_value NUMERIC(4,2) NOT null,
   for_crust_size NUMERIC(4,2) not NULL
)