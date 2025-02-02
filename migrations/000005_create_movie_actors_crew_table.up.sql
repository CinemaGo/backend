CREATE TABLE movie_actors_crew (
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,       -- Foreign key to movies table (with CASCADE delete)
    actor_crew_id INT REFERENCES actors_crew(id) ON DELETE CASCADE, -- Foreign key to actors/crew table (with CASCADE delete)
    PRIMARY KEY (movie_id, actor_crew_id)                        -- Composite primary key
);

CREATE INDEX idx_movie_id ON movie_actors_crew(movie_id);
CREATE INDEX idx_actor_crew_id ON movie_actors_crew(actor_crew_id);
