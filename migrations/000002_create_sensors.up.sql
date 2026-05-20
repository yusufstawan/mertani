CREATE TABLE sensors (
    id UUID PRIMARY KEY,
    device_id UUID NOT NULL,
    sensor_type VARCHAR NOT NULL,
    value NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_sensors_device
        FOREIGN KEY (device_id)
        REFERENCES devices (id)
        ON DELETE CASCADE
);

CREATE INDEX idx_sensors_device_id ON sensors (device_id);
