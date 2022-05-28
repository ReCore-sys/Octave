import { State, Update, SetPlaylist } from "./utils";

let ws = new WebSocket("ws://localhost:4678/ws");

ws.onmessage = async function (e) {
    let data = JSON.parse(e.data);
    switch (data.type) {
        case "heartbeat":
            let out = {
                type: "heartbeat",
                id: data.id,
            };
            ws.send(JSON.stringify(out));
        case "update":
            let state = await State();
            Update();
            SetPlaylist(state.activePlaylist);
    }
};
