-- +goose Up
alter table users add column hashed_password text not null default 'unspecified';

-- +goose Down
alter table users drop column hashed_password;
