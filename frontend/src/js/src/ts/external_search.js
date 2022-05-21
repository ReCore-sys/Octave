import { fatal } from "./utils.js";
$(".search-bar").on("keypress", async function (e) {
    let utils = await utilspromise;
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
            utils.debug(`Searching for ${searchq}`);
            let resprom = window.go.main.App.Search(searchq);
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
        }
        else {
            window.go.main.App.JSLog("warn", `Search was not a string\n${searchq}`);
        }
    }
});
function setupSearch() {
    $(".search-result").tilt({
        maxTilt: 5,
    });
    $(".search-result").on("click", async function () {
        await ClickObesity($(this));
    });
}
async function ClickObesity(e) {
    console.log(0);
    let utils = await utilspromise;
    console.log(1);
    let id = e.attr("id");
    if (typeof id == "undefined") {
        console.log("Well fuck");
        utils.fatal(`Clicked on a search result that doesn't have an id`);
        return;
    }
    console.log(2);
    let song = await window.go.main.App.InCache(id);
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
    if (!(await window.go.main.App.SongDownloaded(id))) {
        downloadbutton.text("Downloaded");
    }
    else {
        downloadbutton.append($(`<embed src="assets/svg/thumbs_up.svg">`));
    }
    downloadbutton.attr("song-id", song.id);
    downloadbutton.on("click", async function () {
        let id = $(this).attr("song-id");
        if (typeof id == "undefined") {
            utils.fatal(`Clicked on a download button that doesn't have an id`);
            return;
        }
        let worked = await window.go.main.App.Save(id);
        if (worked) {
            DownloadAnim();
        }
        AddBadges();
    });
    newmodal.append(downloadbutton);
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
    if (!container.is(e.target) && container.has(e.target).length === 0) {
        UnFatten();
    }
});
function DownloadAnim() {
    let md = $(".modal-download");
    md.text("");
    md.append($(`<embed src="assets/svg/thumbs_up.svg">`));
}
async function AddBadges() {
    let songcards = $(".search-result");
    for (let song of songcards) {
        let id = $(song).attr("id");
        if (typeof id == "undefined") {
            fatal(`Search result doesn't have an id`);
            return;
        }
        let isdownloaded = await window.go.main.App.SongDownloaded(id);
        if (isdownloaded) {
            let icondiv = $("<div>")
                .addClass("downloaded-icon")
                .append($("<embed>").attr("src", "assets/svg/save.svg"));
            $(`#${id}`).append(icondiv);
        }
    }
}