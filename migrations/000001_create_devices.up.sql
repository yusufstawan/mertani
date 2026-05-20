CREATE TABLE devices (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    location VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
