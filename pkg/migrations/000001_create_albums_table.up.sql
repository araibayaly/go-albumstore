CREATE TABLE IF NOT EXISTS albums
(
    ID        bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    Title text NOT NULL,
    Artist text NOT NULL,
    Genre text NOT NULL,
    Year text NOT NULL
); 