import { wsc } from "./wsc"

let wscConn: wsc.Connection;

function initConn() {
    wscConn = new wsc.Connection();

    wscConn.setEstablishHandler((ws) => {
        console.log("connected");
        wscConn.sendMsg(2, new TextEncoder().encode("js request"));
        for (let i = 0; i < 10; i++) {
            wscConn.sendRequestMsg(2,new TextEncoder().encode("js request rpc"+i), (ws, resp) => {
                const bytesString = new TextDecoder().decode(resp)
                console.log(i, bytesString);
            });
        }
    });

    wscConn.setErrHandler((ws, error) => {
        console.log("err" + error);
    });

    wscConn.setCloseHandler((ws, e) => {
        console.log("closed", e);
    });
    wscConn.setDisplacedHandler((ws, oldIp, newIp, ts) => {
        console.log(oldIp, " displaced by ", newIp, " at ", ts);
    });
    wscConn.registerMsgHandler(3, (ws, data) => {
        console.log(new TextDecoder().decode(data));
    });
    wscConn.connect("ws://127.0.0.1:8003/join?uid=x10000", 2000);
}

initConn();