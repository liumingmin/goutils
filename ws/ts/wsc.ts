import * as $msg_pb from "./msg_pb"

export namespace wsc {
    type MsgHandler = (ws: WebSocket, buffer: Uint8Array) => void;
    type EvtHanlder = (ws: WebSocket, evt?: Event) => void;

    export class Connection {
        private ws: WebSocket = null;
        private connected: boolean = false;
        private msgHandler: Map<number, MsgHandler> = new Map();
        private establishHandler: EvtHanlder = null;
        private errHandler: EvtHanlder = null;
        private closeHandler: EvtHanlder = null;
        private static packetHeadFlag: Uint8Array = new Uint8Array([254, 238]);

        constructor() {
        }

        registerMsgHandler(protocolId: number, handler: MsgHandler) {
            this.msgHandler.set(protocolId, handler);
        }

        setEstablishHandler(establishHandler: EvtHanlder) {
            this.establishHandler = establishHandler;
        }

        setErrHandler(errHandler: EvtHanlder) {
            this.errHandler = errHandler;
        }

        setCloseHandler(closeHandler: EvtHanlder) {
            this.closeHandler = closeHandler;
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
                if (this.establishHandler !== null) {
                    this.establishHandler(this.ws, e);
                }
            };
            this.ws.onerror = (error) => {
                if (this.errHandler !== null) {
                    this.errHandler(this.ws, error);
                }
            };

            this.ws.onmessage = (e) => {
                let msgPack = this.unpackMsg(e.data);
                //console.log(msgPack);
                let wsMessage = $msg_pb.ws.P_MESSAGE.decode(new Uint8Array(msgPack.dataBuffer));
                let handler = this.msgHandler.get(wsMessage.protocolId);
                if (handler !== null) {
                    handler(this.ws, wsMessage.data);
                }
            };
            this.ws.onclose = (e) => {
                this.connected = false;
                if (this.closeHandler !== null) {
                    this.closeHandler(this.ws, e);
                }
                this.ws = null;

                if (retryInterval === undefined) return;

                setTimeout(() => this.connect(url, retryInterval), retryInterval);
            };
        }

        sendMsg(protocolId: number, data: Uint8Array) {
            let wsMessage = new $msg_pb.ws.P_MESSAGE();
            wsMessage.protocolId = protocolId;
            wsMessage.data = data;
            this.ws.send(this.packMsg($msg_pb.ws.P_MESSAGE.encode(wsMessage).finish()));
        }

        unpackMsg(buffer: ArrayBuffer) {
            let packetHeadFlag = buffer.slice(0, 2);
            const dv = new DataView(buffer.slice(2, 6));
            const packetLength = dv.getUint32(0, /* little endian data */ true);
            let dataBuffer = buffer.slice(6);
            return {
                packetHeadFlag,
                packetLength,
                dataBuffer
            };
        }

        packMsg(dataArray: Uint8Array): ArrayBuffer {
            let packetLength = new Uint8Array(4);
            new DataView(packetLength.buffer).setUint32(0, dataArray.byteLength, true /* littleEndian */);
            let packet = new Uint8Array([...Connection.packetHeadFlag, ...packetLength, ...dataArray]);
            return packet.buffer;
        }
    }
}