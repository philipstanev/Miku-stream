-- +goose Up
alter table shows 
add column showPath text;

-- +goose Down
alter table shows
drop column showPath;
