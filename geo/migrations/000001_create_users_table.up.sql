CREATE TABLE IF NOT EXISTS search_history(
    id SERIAL PRIMARY KEY,
    query TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS address(
    id SERIAL PRIMARY KEY,
    lat TEXT,
    lon TEXT
);

CREATE TABLE IF NOT EXISTS history_search_address(
    id SERIAL PRIMARY KEY,
    search_id INT REFERENCES search_history(id),
    address_id INT REFERENCES address(id)
);

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;