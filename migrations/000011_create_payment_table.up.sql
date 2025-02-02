CREATE TABLE payment (
    payment_id SERIAL PRIMARY KEY,            -- Unique ID for each payment (auto-incremented)
    amount INT NOT NULL,                      -- Amount paid for the booking (in integer, e.g., cents or whole currency units)
    remote_transaction_id VARCHAR(255),       -- External unique transaction ID from the payment provider (e.g., PayPal or credit card gateway)
    payment_method VARCHAR(50),               -- Method used for payment (e.g., "Credit Card", "PayPal", "Bank Transfer", etc.)
    booking_id INT REFERENCES booking(booking_id) ON DELETE CASCADE  -- Foreign key to the booking table (links payment to a specific booking)
);
