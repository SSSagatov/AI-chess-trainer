CREATE TABLE games (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pgn TEXT NOT NULL,
    white_username TEXT NOT NULL,
    black_username TEXT NOT NULL,
    played_at TIMESTAMPTZ NOT NULL,
    white_rating INTEGER NOT NULL CHECK (white_rating BETWEEN 1 AND 65535),
    white_accuracy NUMERIC(5,2),
    black_rating INTEGER NOT NULL CHECK (black_rating BETWEEN 1 AND 65535),
    black_accuracy NUMERIC(5,2),
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 65535),
    player_color VARCHAR(5) NOT NULL,
    result VARCHAR(16) NOT NULL,
    fen TEXT NOT NULL
);

CREATE INDEX games_rating_idx ON games (rating);
