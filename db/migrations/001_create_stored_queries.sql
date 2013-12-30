-- +goose Up
CREATE TABLE QUERIES (
    id int NOT NULL,
    queryString text,
    --body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE QUERIES;
