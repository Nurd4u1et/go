CREATE TABLE IF NOT EXISTS movies (
                                      id SERIAL PRIMARY KEY,
                                      title TEXT NOT NULL,
                                      year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS actors (
                                      id SERIAL PRIMARY KEY,
                                      movie_id INT NOT NULL REFERENCES movies(id),
    name TEXT NOT NULL
    );

INSERT INTO movies (title, year) VALUES ('Inception', 2010) ON CONFLICT DO NOTHING;
INSERT INTO movies (title, year) VALUES ('Batman', 2014) ON CONFLICT DO NOTHING;

INSERT INTO actors (movie_id, name) VALUES (1, 'Leonardo DiCaprio') ON CONFLICT DO NOTHING;
INSERT INTO actors (movie_id, name) VALUES (1, 'Christian Bale') ON CONFLICT DO NOTHING;
