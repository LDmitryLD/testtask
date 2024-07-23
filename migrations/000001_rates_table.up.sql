CREATE TABLE IF NOT EXISTS rates (
    id SERIAL PRIMARY KEY,
    timestamp BIGINT,
    asks JSON,
    bids JSON
);