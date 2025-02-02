CREATE TABLE actors_crew (
    id SERIAL PRIMARY KEY,               -- Unique ID for actor/crew
    full_name VARCHAR(255),               -- Full name of the actor/crew member
    image_url VARCHAR(255),               -- URL to the image
    occupation VARCHAR(100),              -- Occupation (e.g., Actor, Director, etc.)
    role_description VARCHAR(100),        -- Role description (e.g., "Nina Johnson as Dina")
    born_date VARCHAR(10),                -- Birthdate of the actor/crew member (as a string, e.g., "1985-04-15")
    birthplace VARCHAR(255),              -- Birthplace of the actor/crew member
    about TEXT,                           -- Short biography of the actor/crew member
    is_actor BOOLEAN                      -- Indicates whether the person is an actor (TRUE for actor, FALSE for other roles)
);
