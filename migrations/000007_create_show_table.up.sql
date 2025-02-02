CREATE TABLE show (
    show_id SERIAL PRIMARY KEY,           -- Unique ID for each show (auto-incremented)
    show_date DATE NOT NULL,               -- Date of the show (e.g., '2025-02-02')
    start_time TIME NOT NULL,              -- Start time of the show (e.g., '14:30')
    hall_id INT REFERENCES cinema_hall(cinema_hall_id) ON DELETE CASCADE,  -- Foreign key to the cinema_hall table, links to the hall where the show is held
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,  -- Foreign key to the movies table, links to the movie being shown
    UNIQUE (hall_id, show_date, start_time)  -- Prevents double booking a hall at the same time (same hall, same date, same time)
);

CREATE INDEX idx_show_hall_id ON show (hall_id);
CREATE INDEX idx_show_movie_id ON show (movie_id);
