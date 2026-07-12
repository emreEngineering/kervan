CREATE TABLE products(
    id BIGSERIAL PRIMARY KEY ,
    name VARCHAR(255) NOT NULL ,
    description TEXT NOT NULL DEFAULT '',
    price DECIMAL(10,2) NOT NULL CHECK (price>0),
    category VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);