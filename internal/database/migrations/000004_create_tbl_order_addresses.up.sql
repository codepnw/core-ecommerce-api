CREATE TABLE order_addresses (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    address_id UUID NOT NULL REFERENCES addresses(id),
    address_line TEXT,
    city TEXT,
    state TEXT,
    postal_code TEXT,
    phone TEXT
);