CREATE TABLE show_seat (
    show_seat_id SERIAL PRIMARY KEY,              -- Unique ID for each seat in a specific show (auto-incremented)
    cinema_seat_id INT REFERENCES cinema_seat(cinema_seat_id) ON DELETE CASCADE,  -- Foreign key to cinema_seat, links to specific seat in the cinema
    status VARCHAR(50) NOT NULL,                   -- Status of the seat (e.g., 'Available', 'Booked', 'Selected')
    price INT,                                     -- Price for this specific seat (now as an integer value)
    show_id INT REFERENCES show(show_id) ON DELETE CASCADE,  -- Foreign key to the show table, links to a specific show
    booking_id INT REFERENCES booking(booking_id) ON DELETE SET NULL -- Foreign key to booking table, links to a booking (nullable if not booked)
);

-- Indexes for better query performance on foreign key columns
CREATE INDEX idx_show_seat_cinema_seat_id ON show_seat (cinema_seat_id);   -- Index on cinema_seat_id for fast lookups by seat
CREATE INDEX idx_show_seat_show_id ON show_seat (show_id);                 -- Index on show_id for fast lookups by show
CREATE INDEX idx_show_seat_booking_id ON show_seat (booking_id);           -- Index on booking_id for fast lookups by booking
