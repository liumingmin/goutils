import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace ws. */
export namespace ws {

    /** Properties of a P_MESSAGE. */
    interface IP_MESSAGE {

        /** P_MESSAGE protocolId */
        protocolId?: (number|null);

        /** P_MESSAGE data */
        data?: (Uint8Array|null);
    }

    /** Represents a P_MESSAGE. */
    class P_MESSAGE implements IP_MESSAGE {

        /**
         * Constructs a new P_MESSAGE.
         * @param [properties] Properties to set
         */
        constructor(properties?: ws.IP_MESSAGE);

        /** P_MESSAGE protocolId. */
        public protocolId: number;

        /** P_MESSAGE data. */
        public data: Uint8Array;

        /**
         * Creates a new P_MESSAGE instance using the specified properties.
         * @param [properties] Properties to set
         * @returns P_MESSAGE instance
         */
        public static create(properties?: ws.IP_MESSAGE): ws.P_MESSAGE;

        /**
         * Encodes the specified P_MESSAGE message. Does not implicitly {@link ws.P_MESSAGE.verify|verify} messages.
         * @param message P_MESSAGE message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: ws.IP_MESSAGE, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified P_MESSAGE message, length delimited. Does not implicitly {@link ws.P_MESSAGE.verify|verify} messages.
         * @param message P_MESSAGE message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: ws.IP_MESSAGE, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a P_MESSAGE message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns P_MESSAGE
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ws.P_MESSAGE;

        /**
         * Decodes a P_MESSAGE message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns P_MESSAGE
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ws.P_MESSAGE;

        /**
         * Verifies a P_MESSAGE message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a P_MESSAGE message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns P_MESSAGE
         */
        public static fromObject(object: { [k: string]: any }): ws.P_MESSAGE;

        /**
         * Creates a plain object from a P_MESSAGE message. Also converts values to other types if specified.
         * @param message P_MESSAGE
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: ws.P_MESSAGE, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this P_MESSAGE to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for P_MESSAGE
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a P_DISPLACE. */
    interface IP_DISPLACE {

        /** P_DISPLACE oldIp */
        oldIp?: (Uint8Array|null);

        /** P_DISPLACE newIp */
        newIp?: (Uint8Array|null);

        /** P_DISPLACE ts */
        ts?: (number|Long|null);
    }

    /** Represents a P_DISPLACE. */
    class P_DISPLACE implements IP_DISPLACE {

        /**
         * Constructs a new P_DISPLACE.
         * @param [properties] Properties to set
         */
        constructor(properties?: ws.IP_DISPLACE);

        /** P_DISPLACE oldIp. */
        public oldIp: Uint8Array;

        /** P_DISPLACE newIp. */
        public newIp: Uint8Array;

        /** P_DISPLACE ts. */
        public ts: (number|Long);

        /**
         * Creates a new P_DISPLACE instance using the specified properties.
         * @param [properties] Properties to set
         * @returns P_DISPLACE instance
         */
        public static create(properties?: ws.IP_DISPLACE): ws.P_DISPLACE;

        /**
         * Encodes the specified P_DISPLACE message. Does not implicitly {@link ws.P_DISPLACE.verify|verify} messages.
         * @param message P_DISPLACE message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: ws.IP_DISPLACE, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified P_DISPLACE message, length delimited. Does not implicitly {@link ws.P_DISPLACE.verify|verify} messages.
         * @param message P_DISPLACE message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: ws.IP_DISPLACE, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a P_DISPLACE message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns P_DISPLACE
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): ws.P_DISPLACE;

        /**
         * Decodes a P_DISPLACE message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns P_DISPLACE
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): ws.P_DISPLACE;

        /**
         * Verifies a P_DISPLACE message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a P_DISPLACE message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns P_DISPLACE
         */
        public static fromObject(object: { [k: string]: any }): ws.P_DISPLACE;

        /**
         * Creates a plain object from a P_DISPLACE message. Also converts values to other types if specified.
         * @param message P_DISPLACE
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: ws.P_DISPLACE, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this P_DISPLACE to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for P_DISPLACE
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** P_BASE enum. */
    enum P_BASE {
        raw_bytes_msg = 0,
        s2c_err_displace = -1
    }
}
