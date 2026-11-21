-- +goose Up
create table shows(
    id uuid primary key default uuidv7(),
    titleEnglish text,
    titleRomanji text,
    show_description text,
    episodes int,
    duration int,
    season text,
    seasonYear int,
    showFormat text,
    showStatus text,
    genres text[],
    coverImageURL text
);

-- +goose Down
drop table shows;
