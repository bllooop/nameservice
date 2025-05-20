-- +goose Up
-- +goose StatementBegin
CREATE TABLE people (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    age INTEGER,
    gender TEXT,
    nationality TEXT
);
CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_surname ON people(surname);
CREATE INDEX idx_people_gender ON people(gender);
CREATE INDEX idx_people_nationality ON people(nationality);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE people;
-- +goose StatementEnd
