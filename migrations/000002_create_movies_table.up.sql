CREATE TABLE movies (
    id SERIAL PRIMARY KEY,           -- Unique movie ID
    title VARCHAR(255) NOT NULL,      -- Movie title
    description TEXT,                 -- Movie description
    genre VARCHAR(50),                -- Movie genre (e.g., Action, Comedy, etc.)
    language VARCHAR(50),             -- Language of the movie (e.g., English, Hindi)
    trailer_url VARCHAR(255),         -- URL to the movie trailer   
    poster_url VARCHAR(255),          -- URL to the movie poster image
    rating INT,                       -- Movie rating (e.g., 7.5/10)
    rating_provider VARCHAR(100),     -- Rating provider (e.g., IMDb, Rotten Tomatoes)
    duration INT,                     -- Movie duration in minutes
    release_date VARCHAR(20),         -- Release date as string (e.g., "2025-01-23")
    age_limit VARCHAR(20),            -- Age rating (e.g., UA13+, A, etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp for movie creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Timestamp for movie updates
);
