-- +goose Up
create table users (
    id            uuid primary key default uuidv7(),
    username      text not null unique,
    password_hash text not null,
    role          text not null default 'user'
);

-- +goose Down
drop table users;