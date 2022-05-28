import { db, settings, global_state } from "../frontend/wailsjs/go/models";
import {
    FindSong,
    GetState,
    JSLog,
    Parse,
    Settings,
} from "../frontend/wailsjs/go/main/App";
import "../frontend/wailsjs/runtime";
import { FillItems } from "./search.js";

var startTime: Date, endTime: Date;

function start() {
    startTime = new Date();
}

function end() {
    endTime = new Date();
    var timeDiff = endTime.getMilliseconds() - startTime.getMilliseconds();
    return timeDiff;
}

export function diff(id: number) {
    let d = end();
    JSLog("info", `${id} took ${d}ms`);
    start();
}
start();
console.log("Loaded!");
/**
 * Update the main page with the content of the file at the given path
 * @param {string} path - The path to the file that you want to load.
 */ export async function UpdateMain(path: string) {
    console.log("Updating main");
    let main = $("#main-page");
    let newcontent = await Parse(path);
    main.html(newcontent);
}
let my_settings: settings.SettingsStruct;
/**
 * When the page loads, get the state, set the playlist, update the page, and add event listeners to
 * the search and home buttons.
 */ export async function setuppage() {
    my_settings = await Settings();
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
/**
 * It returns the state of the application.
 * @returns {Statemap} The state of the application.
 */
export async function State(): Promise<global_state.Statemap> {
    let state = await GetState();
    return state;
}

/**
 * It updates the song name and artist on the page.
 */ export async function Update() {
    let state = await State();
    let song = state.currentSong;

    $(".playlist-song").removeClass("active");
    $(`#${song.id}`).addClass("active");
    if (
        $("#song-name").text() != song.Title ||
        $(".song-art").attr("src") == null
    ) {
        $(".song-art").attr("src", "");
        $(".song-art").attr(
            "src",
            `http://localhost:${my_settings.wsport}/song_img/` + song.Image
        );
    }
    $("#song-name").text(song.Title);
    $("#song-artist").text(song.Artist);
}

/**
 * "UpdateSong" is a function that takes a song object as a parameter and calls the "UpdateSong"
 * function in the "App" class in the "main" namespace in the "go" window object.
 * @param {Song} song - Song
 */ export async function UpdateSong(song: db.Song) {
    UpdateSong(song);
}

/**
 * SetPlaylist is an async function that takes a playlist as a parameter and then sets the playlist
 * name, art, and song list, as well as updating the playlist screen
 * @param {Playlist} playlist - Playlist
 */ export async function SetPlaylist(playlist: db.Playlist) {
    let settings = await Settings();
    if ($("#playlist-name").text() != playlist.name) {
        $("#playlist-art").attr(
            "src",
            `http://localhost:${settings.wsport}/song_img/` + playlist.art
        );
    }
    $("#playlist-name").text(playlist.name);
    let songlist = $(".music-list");
    let activeid = $(".active").attr("id");
    let emptied = false;
    for (let songid of playlist.songs) {
        let song = await FindSong(songid);
        // Check if the song already exists. This saves us some computational power and reduces flickering effect
        if ($("#" + song.id).length == 0) {
            // If it doesn't exists, redo the list
            if (!emptied) {
                songlist.empty();
                emptied = true;
            }
            let newSong = $(`<li id="${song.id}" class="playlist-song">`).on(
                "click",
                function () {
                    clik(song.id);
                }
            ); // Create the new song
            if (song.id == activeid) {
                $(".active").removeClass("active");
                newSong.addClass("active");
            } // If the song is active, add the active class
            newSong.append(`<p id="${song.id}-name">${song.Title}</p>`); // Add the song name
            newSong.append(`<p id="${song.id}-artist">${song.Artist}</p>`); // Add the song artist
            newSong.append(
                `<p id="${song.id}-len">${parsetime(song.Length)}</p>`
            ); // Add the song length
            newSong.append(`<p id="${song.id}-album">${song.Album}</p>`); // Add the song album
            songlist.append(newSong); // Add the song to the list
        } else {
            // Pretty much all we do here is check each value and see if it matches the one provided by the database
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
/**
 * It removes the class "active" from all elements with the class "playlist-song" and then adds the
 * class "active" to the element with the id that was passed to the function.
 * @param {string} id - The id of the song
 */
export async function clik(id: string) {
    console.log("Clicked!");
    SetSong(id);
}
/**
 * It takes a song id, finds the song, and then updates the song and the UI.
 * @param {string} id - The id of the song to set.
 */ export async function SetSong(id: string) {
    let song = await FindSong(id);
    UpdateSong(song);
    Update();
}

/**
 * It takes a number of seconds and returns a string of minutes and seconds
 * @param {number} time - The time in seconds
 * @returns A string with the minutes and seconds.
 */ export function parsetime(time: number): string {
    let minutes = Math.floor(time / 60);
    let seconds = time % 60;
    return `${minutes}:${seconds}`;
}

setuppage();
/**======================
 *    Log functions
 *========================**/
export function error(msg: string) {
    JSLog("error", msg);
}
export function info(msg: string) {
    JSLog("info", msg);
}
export function debug(msg: string) {
    JSLog("debug", msg);
}
export function warn(msg: string) {
    JSLog("warn", msg);
}
export function fatal(msg: string) {
    JSLog("fatal", msg);
}

export function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * It adds a script to the head of the document
 * @param {string} script - The script to add.
 * @param {boolean} [module=false] - boolean = false
 */
export function AddScript(script: string, module: boolean = false) {
    if (!module) {
        $("head").append(
            `<script type='text/javascript' src='${script}'></script>`
        );
    } else {
        $("head").append(`<script type='module' src='${script}'></script>`);
    }
}
