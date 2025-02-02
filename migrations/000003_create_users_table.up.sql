CREATE TABLE users (
    id SERIAL PRIMARY KEY,           -- Unique identifier for each user
    name VARCHAR(50),
    surname VARCHAR(50),
    email VARCHAR(255) UNIQUE NOT NULL,   -- User's email, used for authentication
    phone_number VARCHAR(20),           -- User's phone number (stored as a string to handle different formats)
    password_hash VARCHAR(255) NOT NULL, -- Hashed password for security
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'helpdesk', 'user')),  -- Including multiple roles
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp for user creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Timestamp for user updates
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone_number ON users (phone_number);
CREATE INDEX idx_role ON users(role);
