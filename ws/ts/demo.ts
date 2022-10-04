import {ws} from "./ws"
//import {TextEncoder} from "text-encoder"

function initConn(){
    let conn = new ws.Connection();

    conn.setEstablishHandler((ws) => {
        console.log("connected");
        conn.sendMsg(1, new TextEncoder().encode("js request"));
    });

    conn.setErrHandler((ws, error) => {
        console.log("err" + error);
    });

    conn.setCloseHandler((ws) => {
        console.log("closed");
    });

    conn.registerMsgHandler(2, (ws, data) => {
        console.log(new TextDecoder().decode(data));
    });
    conn.connect("ws://127.0.0.1:8003/join?uid=x10000");
}

initConn();