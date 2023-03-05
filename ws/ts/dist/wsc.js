"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.wsc = void 0;
const $msg_pb = require("./msg_pb");
var wsc;
(function (wsc) {
    class Connection {
        constructor() {
            this.ws = null;
            this.connected = false;
            this.msgHandler = new Map();
            this.establishHandler = null;
            this.errHandler = null;
            this.closeHandler = null;
            this.displacedHandler = null;
            this.snCounter = 0;
            this.snChanMap = new Map();
            this.msgHandler.set($msg_pb.ws.P_BASE.s2c_err_displace, (ws, buffer) => {
                this.onDisplaced(ws, buffer);
            });
        }
        registerMsgHandler(protocolId, handler) {
            this.msgHandler.set(protocolId, handler);
        }
        setEstablishHandler(establishHandler) {
            this.establishHandler = establishHandler;
        }
        setErrHandler(errHandler) {
            this.errHandler = errHandler;
        }
        setCloseHandler(closeHandler) {
            this.closeHandler = closeHandler;
        }
        setDisplacedHandler(displacedHandler) {
            this.displacedHandler = displacedHandler;
        }
        connect(url, retryInterval) {
            if (this.ws !== null) {
                this.ws.close();
                this.connected = false;
            }
            this.ws = new WebSocket(url);
            this.ws.onopen = (e) => {
                var _a;
                this.ws.binaryType = 'arraybuffer';
                this.connected = true;
                (_a = this.establishHandler) === null || _a === void 0 ? void 0 : _a.call(this, this.ws, e);
            };
            this.ws.onerror = (error) => {
                var _a;
                (_a = this.errHandler) === null || _a === void 0 ? void 0 : _a.call(this, this.ws, error);
            };
            this.ws.onmessage = (e) => {
                let msgPack = this.unpackMsg(e.data);
                if (msgPack.sn > 0) {
                    let respCallback = this.snChanMap.get(msgPack.sn);
                    this.snChanMap.delete(msgPack.sn);
                    respCallback === null || respCallback === void 0 ? void 0 : respCallback(this.ws, new Uint8Array(msgPack.dataBuffer));
                }
                let handler = this.msgHandler.get(msgPack.protocolId);
                handler === null || handler === void 0 ? void 0 : handler(this.ws, new Uint8Array(msgPack.dataBuffer));
            };
            this.ws.onclose = (e) => {
                var _a;
                (_a = this.closeHandler) === null || _a === void 0 ? void 0 : _a.call(this, this.ws, e);
                this.connected = false;
                this.ws = null;
                if (retryInterval === undefined)
                    return;
                setTimeout(() => this.connect(url, retryInterval), retryInterval);
            };
        }
        sendMsg(protocolId, data) {
            this.ws.send(this.packMsg(protocolId, 0, data));
        }
        sendRequestMsg(protocolId, data, respCallback) {
            let sn = ++this.snCounter;
            if (sn <= 0) {
                sn = this.snCounter = 1;
            }
            this.snChanMap.set(sn, respCallback);
            this.ws.send(this.packMsg(protocolId, sn, data));
        }
        unpackMsg(buffer) {
            let packetHeadFlag = buffer.slice(0, 2);
            const packetLenDv = new DataView(buffer.slice(2, 6));
            const packetLength = packetLenDv.getUint32(0, true);
            const protocolIdDv = new DataView(buffer.slice(6, 10));
            const protocolId = protocolIdDv.getUint32(0, true);
            const snDv = new DataView(buffer.slice(10, 14));
            const sn = snDv.getUint32(0, true);
            let dataBuffer = buffer.slice(14);
            return {
                packetHeadFlag,
                packetLength,
                protocolId,
                sn,
                dataBuffer
            };
        }
        packMsg(protocolId, sn, dataArray) {
            let packetLength = new Uint8Array(4);
            new DataView(packetLength.buffer).setUint32(0, dataArray.byteLength + 8, true);
            let protocolIdArray = new Uint8Array(4);
            new DataView(protocolIdArray.buffer).setUint32(0, protocolId, true);
            let snArray = new Uint8Array(4);
            new DataView(snArray.buffer).setUint32(0, sn, true);
            let packet = new Uint8Array([...Connection.packetHeadFlag, ...packetLength, ...protocolIdArray, ...snArray, ...dataArray]);
            return packet.buffer;
        }
        onDisplaced(ws, buffer) {
            var _a;
            let displacedMsg = $msg_pb.ws.P_DISPLACE.decode(buffer);
            let oldIp = new TextDecoder().decode(displacedMsg.oldIp);
            let newIp = new TextDecoder().decode(displacedMsg.newIp);
            console.log(oldIp, " displaced by ", newIp, " at ", displacedMsg.ts);
            (_a = this.displacedHandler) === null || _a === void 0 ? void 0 : _a.call(this, this.ws, oldIp, newIp, displacedMsg.ts.toString());
        }
    }
    Connection.packetHeadFlag = new Uint8Array([254, 239]);
    wsc.Connection = Connection;
})(wsc = exports.wsc || (exports.wsc = {}));
