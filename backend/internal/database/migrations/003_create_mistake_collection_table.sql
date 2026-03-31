-- Migration script to create mistake_collection table

CREATE TABLE mistake_collection (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    mistake_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Add any additional constraints or indices if necessary