-- name: GetShowByID :one
SELECT * FROM shows
WHERE id = $1;

-- name: GetShowByName :one
SELECT * FROM shows
WHERE titleEnglish = $1 OR titleRomanji = $1;

-- name: GetFilePaths :many
SELECT showPath FROM shows;

-- name: InsertShows :copyfrom
INSERT INTO shows (showPath, titleEnglish, titleRomanji, show_description, episodes, duration, season, seasonYear, showFormat, showStatus, genres, coverImageURL)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);