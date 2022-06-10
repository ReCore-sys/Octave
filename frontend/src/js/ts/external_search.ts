import {
    InCache,
    JSLog,
    Save,
    Search,
    SongDownloaded,
} from "../wailsjs/go/main/App.js";
import { debug, fatal } from "./utils.js";

$(".search-bar").on("keypress", async function (e) {
    $(".search-bar").css("transition", "all 0.3s ease-in-out");
    if ($(this).is(":focus")) {
        $(this).css({
            "margin-top": "10px",
            width: "700px",
        });
    }
    if (e.code == "Enter") {
        e.preventDefault();
        console.log(0);
        $(".search-bar").addClass("after-search");
        $(".search-bar").trigger("blur");
        var searchq = $(".search-bar").val();
        console.log(searchq);
        if (typeof searchq == "string") {
            debug(`Searching for ${searchq}`);
            let resprom = Search(searchq);
            let songlist = $(".search-results");
            songlist.empty();
            let loader = $(".har-loader");
            loader.css("display", "flex");
            let results = await resprom;
            loader.css("display", "none");
            if (results.length > 0) {
                for (var i = 0; i < results.length; i++) {
                    let newline = $("<li>");
                    let textarea = $("<div>");
                    newline.addClass("search-result");
                    newline.attr("id", results[i].id);
                    textarea.addClass("text-area");
                    let newimg = $("<img>");
                    newimg.addClass("search-result-art");
                    newimg.attr("src", results[i].Image);
                    let newtitle = $("<p>");
                    newtitle.addClass("search-result-title");
                    newtitle.text(results[i].Title).wrapInner("<strong />");
                    let newartist = $("<p>");
                    newartist.addClass("search-result-artist");
                    newartist.text(results[i].Artist);
                    if (results[i].Image != "") {
                        newline.append(newimg);
                    }
                    textarea.append(newtitle);
                    textarea.append(newartist);
                    newline.append(textarea);
                    songlist.append(newline);
                }
            }
            AddBadges();
            console.log("DONE");
            setupSearch();
        } else {
            JSLog("warn", `Search was not a string\n${searchq}`);
        }
    }
});

/**
 * It sets up the search results to be tilted and clickable
 */
function setupSearch() {
    // @ts-ignore
    $(".search-result").tilt({
        maxTilt: 5,
    });
    // Because jquery is a pain, i can't pass a function to the click event, it has to be anonymous
    $(".search-result").on("click", async function () {
        await ClickObesity($(this));
    });
}

async function ClickObesity(e: JQuery) {
    console.log(0);
    // Now we need to pop it out and just generally make it L A R G E
    console.log(1);
    let id = e.attr("id");
    if (typeof id == "undefined") {
        console.log("Well fuck");
        fatal(`Clicked on a search result that doesn't have an id`);
        return;
    }
    console.log(2);
    let song = await InCache(id);
    console.log(3);
    let newmodal = $("<div>");
    newmodal.addClass("modal");
    newmodal.hide();
    newmodal.addClass("modal");
    let header = $(`<h1>`).text(song.Title).addClass("modal-title");
    newmodal.append(header);

    let artistheader = $(`<h2>`).text(song.Artist).addClass("modal-artist");
    newmodal.append(artistheader);

    let downloadbutton = $(`<button>`).addClass("modal-download");
    if (!(await SongDownloaded(id))) {
        downloadbutton.text("Downloaded");
    } else {
        downloadbutton.append($(`<embed src="assets/svg/thumbs_up.svg">`));
    }
    downloadbutton.attr("song-id", song.id);
    downloadbutton.on("click", async function () {
        let id = $(this).attr("song-id");
        if (typeof id == "undefined") {
            fatal(`Clicked on a download button that doesn't have an id`);
            return;
        }

        let worked = await Save(id);
        if (worked) {
            DownloadAnim();
        }
        AddBadges();
    });
    newmodal.append(downloadbutton);
    /* Appending the modal to the content div. */
    $(".content").append(newmodal);
    $(".search-results").fadeOut(0);

    newmodal.fadeIn(200);
}

function UnFatten() {
    $(".modal").fadeOut(200);
    $(".search-results").fadeIn(200);
    $(".modal").remove();
}

$("body").on("click", function (e) {
    var container = $(".modal");

    // if the target of the click isn't the container nor a descendant of the container
    if (!container.is(e.target) && container.has(e.target).length === 0) {
        UnFatten();
    }
});

function DownloadAnim() {
    let md = $(".modal-download");
    md.text("");
    md.append($(`<embed src="assets/svg/thumbs_up.svg">`));
    //md.fadeTo(200, 1);
}

async function AddBadges() {
    let songcards = $(".search-result");
    for (let song of songcards) {
        let id = $(song).attr("id");
        if (typeof id == "undefined") {
            fatal(`Search result doesn't have an id`);
            return;
        }
        let isdownloaded = await SongDownloaded(id);
        if (isdownloaded) {
            let icondiv = $("<div>")
                .addClass("downloaded-icon")
                .append($("<embed>").attr("src", "assets/svg/save.svg"));
            $(`#${id}`).append(icondiv);
        }
    }
}
