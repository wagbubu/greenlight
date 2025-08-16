DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'movies_runtime_check') THEN
        ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'movies_year_check') THEN
        ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'genres_length_check') THEN
        ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
    END IF;
END$$;