let utilspromise = import("./utils.js");

let ws = new WebSocket("ws://localhost:4678/ws");

ws.onmessage = async function (e) {
    let data = JSON.parse(e.data);
    let utils = await utilspromise;
    switch (data.type) {
        case "heartbeat":
            let out = {
                type: "heartbeat",
                id: data.id,
            };
            ws.send(JSON.stringify(out));
        case "update":
            let state = await utils.GetState();
            utils.Update();
            utils.SetPlaylist(state.activePlaylist);
    }
};
