import * as underscore from "underscore";
import { GetAllSongs, UpdateSong } from "../wailsjs/go/main/App.js";
import { clik, parsetime, Update } from "./utils.js";
import "../external/datatables.net/js/jquery.dataTables.min.js";

export async function FillItems() {
    let songs = await GetAllSongs();
    console.log(songs);
    $(".items-holder").empty();

    let table = $(`<table class="song-table">`);
    let header = $(`<thead>`);
    let header_row = $(`<tr class="table-head">`);
    header_row.append($("<th>").text("Title"));
    header_row.append($("<th>").text("Artist"));
    header_row.append($("<th>").text("Length"));
    header_row.append($("<th>").text("Album"));
    header.append(header_row);
    let body = $("<tbody>");
    table.append(header);
    table.append(body);

    $(".items-holder").append(table);
    for (let song of songs) {
        let row = $(`<tr song-id="${song.id}">`);
        let songtitle = $(
            `<td song-id="${song.id}" class="song-title text-truncate">`
        ).text(song.title);
        let artist = $(
            `<td song-id="${song.id}" class="song-artist text-truncate">`
        ).text(song.artist);
        let length = $(
            `<td song-id="${song.id}" class="song-length text-truncate">`
        ).text(parsetime(song.length));
        let album = $(
            `<td song-id="${song.id}" class="song-album text-truncate">`
        ).text(song.album);
        row.append(songtitle);
        row.append(artist);
        row.append(length);
        row.append(album);
        body.append(row);
        row.on("click", async function () {
            await UpdateSong(song);
            Update();
        });
    }
    $("table").DataTable({
        paging: false,
        searching: false,
        info: false,
    });
    $(".item").css("transition", "all 0.3s ease-in-out");
}
console.log("Loaded");
