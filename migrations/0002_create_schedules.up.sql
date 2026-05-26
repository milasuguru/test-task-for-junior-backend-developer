CREATE TABLE IF NOT EXISTS schedules (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'new',
    type        TEXT NOT NULL,
    every_n_days     INT,
    month_day        INT,
    specific_dates   DATE[],
    parity           TEXT,
    generate_days_ahead INT NOT NULL DEFAULT 30,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
