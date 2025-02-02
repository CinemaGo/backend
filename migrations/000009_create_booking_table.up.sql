CREATE TABLE booking (
    booking_id SERIAL PRIMARY KEY,          -- Unique ID for each booking (auto-incremented)
    number_of_seats INT NOT NULL,           -- The number of seats booked
    status VARCHAR(50),                     -- Status of the booking (e.g., "Pending", "Confirmed", "Cancelled")
    user_id INT REFERENCES users(id) ON DELETE CASCADE,  -- Foreign key to users table (links booking to a user), deletes booking when user is deleted
    show_id INT REFERENCES show(show_id) ON DELETE CASCADE  -- Foreign key to show table (links booking to a specific show), deletes booking when show is deleted
);

-- Indexes for faster querying by user and show
CREATE INDEX idx_booking_user_id ON booking (user_id);   -- Index on user_id for fast lookups by user
CREATE INDEX idx_booking_show_id ON booking (show_id);   -- Index on show_id for fast lookups by show
