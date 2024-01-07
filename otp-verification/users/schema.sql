CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone_number TEXT UNIQUE NOT NULL,
    otp TEXT,
    otp_expiration_time TIMESTAMP
);