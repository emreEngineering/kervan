CREATE TABLE stocks(
    product_id BIGINT PRIMARY KEY,
    available INT NOT NULL CHECK(available >= 0)
);