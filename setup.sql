CREATE TABLE "songs" (
	"Title"	TEXT,
	"Album"	TEXT,
	"Artist"	TEXT,
	"Year"	INTEGER,
	"Length"	INTEGER,
	"Genre"	INTEGER,
	"Image"	INTEGER
, "ID"	INTEGER)

CREATE TABLE "playlists" (
	"id"	INTEGER,
	"name"	INTEGER,
	"songlen"	INTEGER
)

CREATE UNIQUE INDEX "Songs_playlist" ON "playlists" (
	"id"	ASC
)

CREATE UNIQUE INDEX "Songs_id" ON "songs" (
	"ID"	ASC
)
