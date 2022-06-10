import { FillItems } from "./library.js";
import { AddScript, setuppage, UpdateMain } from "./utils.js";
import { Parse, CreatePlaylist } from "../wailsjs/go/main/App.js";
export function InitateSidebar() {
    console.log("Initiating sidebar");
    $(".fa-search").on("click", function () {
        UpdateMain("search.html");
        AddScript("js/ts/external_search.js", true);
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
    $(".sidebar-playlist").one("mousedown", OpenDropdown);
}
async function OpenDropdown() {
    console.log("Opening dropdown");
    let contentwidth = $(".content").outerWidth();
    let _90percent = contentwidth * 0.9;
    let marginwidth = $(".content").outerWidth(true) - contentwidth;
    $(".content").css("transition", "all 0.3s ease-in-out");
    $(".content").css({
        width: _90percent + "px",
        "margin-left": marginwidth + (contentwidth - _90percent) + "px",
    });
    let finalmargin = $(".content").outerWidth(true) -
        $(".content").outerWidth();
    $(".sidebar-playlist").addClass("opened-playlist");
    if ($(".sidebar-playlist").attr("active") == "true") {
        CloseDropdown(contentwidth, marginwidth, finalmargin);
        $(".sidebar-playlist").attr("active", "false");
        return;
    }
    $("aside").css({
        width: (marginwidth + (contentwidth - _90percent)) * 0.95 + "px",
    });
    let playlisthtml = $(await Parse("playlists_sidebar.html"));
    playlisthtml.fadeOut(1);
    $(".sidebar-playlist").after(playlisthtml);
    $(".fa-search").fadeOut(1);
    $(".fa-volume-up").fadeOut(1);
    $(".fa-user").fadeOut(1);
    playlisthtml.fadeIn(300);
    $(".sidebar-playlist").one("mouseup", function () {
        console.log("Mouseup");
        $(window).one("mousedown", function () {
            if ($(".sidebar-playlist-list-item:hover").length == 0) {
                CloseDropdown(contentwidth, marginwidth, finalmargin);
            }
        });
    });
    $(".sidebar-playlist").attr("active", "true");
    $(".sidebar-playlist-list").on("click", CreatePlaylistButton);
}
function CloseDropdown(ogwidth, ogmargin, ogsidebar) {
    console.log("Closing dropdown");
    $(".content").css({
        width: ogwidth + "px",
        "margin-left": ogmargin + "px",
    });
    $("aside").css({
        width: ogsidebar + "px",
    });
    $(".sidebar-playlist").removeClass("opened-playlist");
    $(".sidebar-playlist-list").remove();
    $(".sidebar-playlist-list").fadeOut(300);
    $(".fa-search").fadeIn(300);
    $(".fa-volume-up").fadeIn(300);
    $(".fa-user").fadeIn(300);
    $(".sidebar-playlist").one("mousedown", OpenDropdown);
    $(".sidebar-playlist").attr("active", "false");
}
async function CreatePlaylistButton() {
    // @ts-ignore
    const { value: text } = await Swal.fire({
        input: "textarea",
        inputLabel: "Playlist name",
        inputPlaceholder: "What should we call the playlist?",
        inputAttributes: {
            "aria-label": "Type your message here",
        },
        showCancelButton: true,
    });
    if (text) {
        let result = await CreatePlaylist(text);
        console.log(`Created playlist: ${result}`);
        // @ts-ignore
        Swal.fire(`Created playlist ${text}`, "Now go add some funky beats!");
    }
}
