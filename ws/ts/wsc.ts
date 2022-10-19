import * as $msg_pb from "./msg_pb"

export namespace wsc {
    type MsgHandler = (ws: WebSocket, buffer: Uint8Array) => void;
    type EvtHandler = (ws: WebSocket, evt?: Event) => void;
    type DisplacedHandler = (ws: WebSocket, oldIp: string, newIp: string, ts: string) => void;

    export class Connection {
        private ws: WebSocket = null;
        private connected: boolean = false;
        private msgHandler: Map<number, MsgHandler> = new Map();
        private establishHandler: EvtHandler = null;
        private errHandler: EvtHandler = null;
        private closeHandler: EvtHandler = null;
        private displacedHandler: DisplacedHandler = null;
        private static packetHeadFlag: Uint8Array = new Uint8Array([254, 238]);

        constructor() {
            this.msgHandler.set($msg_pb.ws.P_BASE.s2c_err_displace, (ws: WebSocket, buffer: Uint8Array) => { this.onDisplaced(ws, buffer); });
        }

        registerMsgHandler(protocolId: number, handler: MsgHandler) {
            this.msgHandler.set(protocolId, handler);
        }

        setEstablishHandler(establishHandler: EvtHandler) {
            this.establishHandler = establishHandler;
        }

        setErrHandler(errHandler: EvtHandler) {
            this.errHandler = errHandler;
        }

        setCloseHandler(closeHandler: EvtHandler) {
            this.closeHandler = closeHandler;
        }

        setDisplacedHandler(displacedHandler: DisplacedHandler) {
            this.displacedHandler = displacedHandler;
        }

        connect(url: string, retryInterval?: number) {
            if (this.ws !== null) {
                this.ws.close();
                this.connected = false;
            }

            this.ws = new WebSocket(url);
            this.ws.onopen = (e) => {
                this.ws.binaryType = 'arraybuffer';
                this.connected = true;
                
                this.establishHandler?.(this.ws, e);
            };
            this.ws.onerror = (error) => {
                this.errHandler?.(this.ws, error);
            };

            this.ws.onmessage = (e) => {
                let msgPack = this.unpackMsg(e.data);
                //console.log(msgPack);
                let handler = this.msgHandler.get(msgPack.protocolId);
                handler?.(this.ws, new Uint8Array(msgPack.dataBuffer));
            };

            this.ws.onclose = (e) => {
                this.closeHandler?.(this.ws, e);
                this.connected = false;
                this.ws = null;

                if (retryInterval === undefined) return;

                setTimeout(() => this.connect(url, retryInterval), retryInterval);
            };
        }

        sendMsg(protocolId: number, data: Uint8Array) {
            this.ws.send(this.packMsg(protocolId, data));
        }

        unpackMsg(buffer: ArrayBuffer) {
            let packetHeadFlag = buffer.slice(0, 2);
            const packetLenDv = new DataView(buffer.slice(2, 6));
            const packetLength = packetLenDv.getUint32(0, /* little endian data */ true);
            const protocolIdDv = new DataView(buffer.slice(6, 10));
            const protocolId = protocolIdDv.getUint32(0, /* little endian data */ true);
            let dataBuffer = buffer.slice(10);

            return {
                packetHeadFlag,
                packetLength,
                protocolId,
                dataBuffer
            };
        }

        packMsg(protocolId: number, dataArray: Uint8Array): ArrayBuffer {
            let packetLength = new Uint8Array(4);
            new DataView(packetLength.buffer).setUint32(0, dataArray.byteLength+4, true /* littleEndian */);

            let protocolIdArray = new Uint8Array(4);
            new DataView(protocolIdArray.buffer).setUint32(0, protocolId, true /* littleEndian */);

            let packet = new Uint8Array([...Connection.packetHeadFlag, ...packetLength, ...protocolIdArray, ...dataArray]);
            return packet.buffer;
        }

        onDisplaced(ws: WebSocket, buffer: Uint8Array) {
            let displacedMsg = $msg_pb.ws.P_DISPLACE.decode(buffer);
            let oldIp = new TextDecoder().decode(displacedMsg.oldIp);
            let newIp = new TextDecoder().decode(displacedMsg.newIp);
            console.log(oldIp, " displaced by ", newIp, " at ", displacedMsg.ts);
            this.displacedHandler?.(this.ws, oldIp, newIp, displacedMsg.ts.toString());
        }
    }
}