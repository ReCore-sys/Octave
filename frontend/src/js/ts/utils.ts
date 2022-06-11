import {
    FindSong,
    GetState,
    JSLog,
    Parse,
    Settings,
    UpdateSong,
} from "../wailsjs/go/main/App.js";
import { Quit, WindowMinimise } from "../wailsjs/runtime/runtime.js";
import { db, globalstate, settings } from "../wailsjs/go/models";
import { FillItems } from "./library.js";
import { InitateSidebar } from "./sidebar.js";
export { GetState };

console.log("Loaded!");
/**
 * Update the main page with the content of the file at the given path
 * @param {string} path - The path to the file that you want to load.
 */
export async function UpdateMain(path: string) {
    console.log("Updating main to " + path);
    let main = $("#main-page");
    let newcontent = await Parse(path);
    main.html(newcontent);
}
let my_settings: settings.SettingStruct;
/**
 * When the page loads, get the state, set the playlist, update the page, and add event listeners to
 * the search and home buttons.
 */ export async function setuppage() {
    $(".control-minimise").click(MinimiseWindow);
    $(".control-close").click(CloseWindow);
    my_settings = await Settings();
    let state = await GetState();
    console.log("Setting up page");

    SetPlaylist(state.activePlaylist);
    Update();
    InitateSidebar();
}

export function CloseWindow() {
    console.log("Closing");
    Quit();
}

export function MinimiseWindow() {
    console.log("Minimising");
    WindowMinimise();
}

export async function FirstLoad() {
    UpdateMain("library.html");
    FillItems();
}
FirstLoad();

/**
 * It updates the song name and artist on the page.
 */ export async function Update() {
    let state = await GetState();
    let song = state.currentSong;

    $(".playlist-song").removeClass("active");
    $(`#${song.id}`).addClass("active");
    if (
        $("#song-name").text() != song.title ||
        $(".song-art").attr("src") == null
    ) {
        $(".song-art").attr("src", "");
        $(".song-art").attr(
            "src",
            `http://localhost:${my_settings.wsport}/song_img/` + song.image
        );
    }
    $("#song-name").text(song.title);
    $("#song-artist").text(song.artist);
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
    for (let current_song of playlist.songs) {
        let song = await FindSong(current_song.id);
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
            newSong.append(`<p id="${song.id}-name">${song.title}</p>`); // Add the song name
            newSong.append(`<p id="${song.id}-artist">${song.artist}</p>`); // Add the song artist
            newSong.append(
                `<p id="${song.id}-len">${parsetime(song.length)}</p>`
            ); // Add the song length
            newSong.append(`<p id="${song.id}-album">${song.album}</p>`); // Add the song album
            songlist.append(newSong); // Add the song to the list
        } else {
            // Pretty much all we do here is check each value and see if it matches the one provided by the database
            let currentname = $("#" + song.id + "-name").text();
            if (currentname != song.title) {
                $("#" + song.id + "-name").text(song.title);
            }
            let currentartist = $("#" + song.id + "-artist").text();
            if (currentartist != song.artist) {
                $("#" + song.id + "-artist").text(song.artist);
            }
            let currentalbum = $("#" + song.id + "-album").text();
            if (currentalbum != song.album) {
                $("#" + song.id + "-album").text(song.album);
            }
            let currentlength = $("#" + song.id + "-len").text();
            if (currentlength != parsetime(song.length)) {
                $("#" + song.id + "-len").text(parsetime(song.length));
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

$(".music-list").on("scroll", function () {
    let doctop = $(".music-list").scrollTop();
    if (doctop === undefined) {
        return;
    } else {
        $(".music-head").toggleClass("shrink", doctop > 0);
        $(".playlist-art").toggleClass("shrink", doctop > 0);
        $(".music-list").toggleClass("shrink", doctop > 0);
    }
});
