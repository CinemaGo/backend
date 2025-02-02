CREATE TABLE cinema_seat (
    cinema_seat_id SERIAL PRIMARY KEY,      -- Unique ID for each cinema seat (auto-incremented)
    seat_row VARCHAR(10) NOT NULL,          -- Row identifier (e.g., 'A', 'B', 'C', 'D')
    seat_number INT NOT NULL,               -- Seat number within the row (e.g., 1, 2, 3, ...)
    seat_type VARCHAR(20),                  -- Type of seat (e.g., 'Regular', 'VIP', 'Accessible')
    hall_id INT REFERENCES cinema_hall(cinema_hall_id) ON DELETE CASCADE, -- Foreign key to the cinema hall
    UNIQUE (hall_id, seat_row, seat_number) -- Ensures no duplicate seats in the same hall
);

CREATE INDEX idx_cinema_seat_hall_id ON cinema_seat (hall_id);
