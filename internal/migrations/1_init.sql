CREATE TABLE shipments (
    id TEXT PRIMARY KEY,
    reference TEXT NOT NULL,
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    status TEXT NOT NULL,
    cost NUMERIC NOT NULL,
    driver_revenue NUMERIC NOT NULL
);

CREATE TABLE shipment_events (
    id SERIAL PRIMARY KEY,
    shipment_id TEXT NOT NULL REFERENCES shipments(id),
    status TEXT NOT NULL,
    timestamp BIGINT NOT NULL
);