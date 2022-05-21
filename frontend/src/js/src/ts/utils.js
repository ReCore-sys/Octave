import { FillItems } from "./search.js";
var startTime, endTime;
function start() {
    startTime = new Date();
}
function end() {
    endTime = new Date();
    var timeDiff = endTime.getMilliseconds() - startTime.getMilliseconds();
    return timeDiff;
}
export function diff(id) {
    let d = end();
    window.go.main.App.JSLog("info", `${id} took ${d}ms`);
    start();
}
start();
console.log("Loaded!");
export async function UpdateMain(path) {
    console.log("Updating main");
    let main = $("#main-page");
    let newcontent = await window.go.main.App.Parse(path);
    main.html(newcontent);
}
let settings;
export async function setuppage() {
    settings = await window.go.main.App.Settings();
    let state = await State();
    console.log("Setting up page");
    SetPlaylist(state.activePlaylist);
    Update();
    $(".fa-search").on("click", function () {
        UpdateMain("search.html");
        AddScript("js/src/ts/external_search.js", true);
    });
    $(".fa-home").on("click", function () {
        UpdateMain("home.html");
        setuppage();
    });
    $(".fa-bars").on("click", function () {
        UpdateMain("library.html");
        AddScript("js/external/sort.js");
        FillItems();
    });
}
export async function State() {
    let state = await window.go.main.App.GetState();
    return state;
}
export async function Update() {
    let state = await State();
    let song = state.currentSong;
    $(".playlist-song").removeClass("active");
    $(`#${song.id}`).addClass("active");
    if ($("#song-name").text() != song.Title ||
        $(".song-art").attr("src") == null) {
        $(".song-art").attr("src", "");
        $(".song-art").attr("src", `http://localhost:${settings.wsport}/song_img/` + song.Image);
    }
    $("#song-name").text(song.Title);
    $("#song-artist").text(song.Artist);
}
export async function UpdateSong(song) {
    window.go.main.App.UpdateSong(song);
}
export async function SetPlaylist(playlist) {
    let settings = await window.go.main.App.Settings();
    if ($("#playlist-name").text() != playlist.name) {
        $("#playlist-art").attr("src", `http://localhost:${settings.wsport}/song_img/` + playlist.art);
    }
    $("#playlist-name").text(playlist.name);
    let songlist = $(".music-list");
    let activeid = $(".active").attr("id");
    let emptied = false;
    for (let songid of playlist.songs) {
        let song = await window.go.main.App.FindSong(songid);
        if ($("#" + song.id).length == 0) {
            if (!emptied) {
                songlist.empty();
                emptied = true;
            }
            let newSong = $(`<li id="${song.id}" class="playlist-song">`).on("click", function () {
                clik(song.id);
            });
            if (song.id == activeid) {
                $(".active").removeClass("active");
                newSong.addClass("active");
            }
            newSong.append(`<p id="${song.id}-name">${song.Title}</p>`);
            newSong.append(`<p id="${song.id}-artist">${song.Artist}</p>`);
            newSong.append(`<p id="${song.id}-len">${parsetime(song.Length)}</p>`);
            newSong.append(`<p id="${song.id}-album">${song.Album}</p>`);
            songlist.append(newSong);
        }
        else {
            let currentname = $("#" + song.id + "-name").text();
            if (currentname != song.Title) {
                $("#" + song.id + "-name").text(song.Title);
            }
            let currentartist = $("#" + song.id + "-artist").text();
            if (currentartist != song.Artist) {
                $("#" + song.id + "-artist").text(song.Artist);
            }
            let currentalbum = $("#" + song.id + "-album").text();
            if (currentalbum != song.Album) {
                $("#" + song.id + "-album").text(song.Album);
            }
            let currentlength = $("#" + song.id + "-len").text();
            if (currentlength != parsetime(song.Length)) {
                $("#" + song.id + "-len").text(parsetime(song.Length));
            }
        }
    }
}
export async function clik(id) {
    console.log("Clicked!");
    SetSong(id);
}
export async function SetSong(id) {
    let song = await window.go.main.App.FindSong(id);
    UpdateSong(song);
    Update();
}
export function parsetime(time) {
    let minutes = Math.floor(time / 60);
    let seconds = time % 60;
    return `${minutes}:${seconds}`;
}
setuppage();
export function error(msg) {
    window.go.main.App.JSLog("error", msg);
}
export function info(msg) {
    window.go.main.App.JSLog("info", msg);
}
export function debug(msg) {
    window.go.main.App.JSLog("debug", msg);
}
export function warn(msg) {
    window.go.main.App.JSLog("warn", msg);
}
export function fatal(msg) {
    window.go.main.App.JSLog("fatal", msg);
}
export function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
export function AddScript(script, module = false) {
    if (!module) {
        $("head").append(`<script type='text/javascript' src='${script}'></script>`);
    }
    else {
        $("head").append(`<script type='module' src='${script}'></script>`);
    }
}
