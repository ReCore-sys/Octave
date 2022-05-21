import { parsetime, Update, UpdateSong } from "./utils.js";
export async function FillItems() {
    let songs = await window.go.main.App.GetAllSongs();
    console.log(songs);
    $(".items-list").empty();
    for (let song of songs) {
        let newSong = $(`<li song-id="${song.id}" class="item">`);
        let table = $(`<table song-id="${song.id}" class="song-table">`);
        let row = $(`<tr song-id="${song.id}">`);
        let title = $(`<td song-id="${song.id}" class="song-title">`).text(song.Title);
        let artist = $(`<td song-id="${song.id}" class="song-artist">`).text(song.Artist);
        let length = $(`<td song-id="${song.id}" class="song-length">`).text(parsetime(song.Length));
        let album = $(`<td song-id="${song.id}" class="song-album">`).text(song.Album);
        row.append(title);
        row.append(artist);
        row.append(length);
        row.append(album);
        table.append(row);
        newSong.append(table);
        newSong.on("click", async function () {
            await UpdateSong(song);
            Update();
        });
        $(".items-list").append(newSong);
    }
    $(".items-list").sortable({});
    $(".item").css("transition", "all 0.3s ease-in-out");
}
console.log("Loaded");
