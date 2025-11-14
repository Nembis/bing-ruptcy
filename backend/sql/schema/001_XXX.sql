-- +goose Up

CREATE TABLE "user"(
    user_id SERIAL NOT NULL,
    username VARCHAR(255) NOT NULL
);

-- +goose Down

DROP TABLE "user";
