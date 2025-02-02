CREATE TABLE cinema_hall (
    cinema_hall_id SERIAL PRIMARY KEY,  -- Unique ID for each cinema hall (auto-incremented)
    hall_name VARCHAR(255) NOT NULL,    -- Name of the cinema hall (e.g., 'Hall 1', 'Main Hall')
    hall_type VARCHAR(50),              -- Type of hall (e.g., IMAX, 3D, Regular)
    capacity INT,                       -- The number of seats in the hall (optional)
    UNIQUE (hall_name, hall_type)       -- Enforce unique hall names per type to avoid duplicates
);
