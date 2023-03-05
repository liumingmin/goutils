/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

$root.ws = (function() {

    /**
     * Namespace ws.
     * @exports ws
     * @namespace
     */
    var ws = {};

    ws.P_DISPLACE = (function() {

        /**
         * Properties of a P_DISPLACE.
         * @memberof ws
         * @interface IP_DISPLACE
         * @property {Uint8Array|null} [oldIp] P_DISPLACE oldIp
         * @property {Uint8Array|null} [newIp] P_DISPLACE newIp
         * @property {number|Long|null} [ts] P_DISPLACE ts
         */

        /**
         * Constructs a new P_DISPLACE.
         * @memberof ws
         * @classdesc Represents a P_DISPLACE.
         * @implements IP_DISPLACE
         * @constructor
         * @param {ws.IP_DISPLACE=} [properties] Properties to set
         */
        function P_DISPLACE(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * P_DISPLACE oldIp.
         * @member {Uint8Array} oldIp
         * @memberof ws.P_DISPLACE
         * @instance
         */
        P_DISPLACE.prototype.oldIp = $util.newBuffer([]);

        /**
         * P_DISPLACE newIp.
         * @member {Uint8Array} newIp
         * @memberof ws.P_DISPLACE
         * @instance
         */
        P_DISPLACE.prototype.newIp = $util.newBuffer([]);

        /**
         * P_DISPLACE ts.
         * @member {number|Long} ts
         * @memberof ws.P_DISPLACE
         * @instance
         */
        P_DISPLACE.prototype.ts = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new P_DISPLACE instance using the specified properties.
         * @function create
         * @memberof ws.P_DISPLACE
         * @static
         * @param {ws.IP_DISPLACE=} [properties] Properties to set
         * @returns {ws.P_DISPLACE} P_DISPLACE instance
         */
        P_DISPLACE.create = function create(properties) {
            return new P_DISPLACE(properties);
        };

        /**
         * Encodes the specified P_DISPLACE message. Does not implicitly {@link ws.P_DISPLACE.verify|verify} messages.
         * @function encode
         * @memberof ws.P_DISPLACE
         * @static
         * @param {ws.IP_DISPLACE} message P_DISPLACE message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        P_DISPLACE.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.oldIp != null && Object.hasOwnProperty.call(message, "oldIp"))
                writer.uint32(/* id 1, wireType 2 =*/10).bytes(message.oldIp);
            if (message.newIp != null && Object.hasOwnProperty.call(message, "newIp"))
                writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.newIp);
            if (message.ts != null && Object.hasOwnProperty.call(message, "ts"))
                writer.uint32(/* id 3, wireType 0 =*/24).int64(message.ts);
            return writer;
        };

        /**
         * Encodes the specified P_DISPLACE message, length delimited. Does not implicitly {@link ws.P_DISPLACE.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.P_DISPLACE
         * @static
         * @param {ws.IP_DISPLACE} message P_DISPLACE message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        P_DISPLACE.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a P_DISPLACE message from the specified reader or buffer.
         * @function decode
         * @memberof ws.P_DISPLACE
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.P_DISPLACE} P_DISPLACE
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        P_DISPLACE.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.P_DISPLACE();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.oldIp = reader.bytes();
                        break;
                    }
                case 2: {
                        message.newIp = reader.bytes();
                        break;
                    }
                case 3: {
                        message.ts = reader.int64();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a P_DISPLACE message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.P_DISPLACE
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.P_DISPLACE} P_DISPLACE
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        P_DISPLACE.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a P_DISPLACE message.
         * @function verify
         * @memberof ws.P_DISPLACE
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        P_DISPLACE.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.oldIp != null && message.hasOwnProperty("oldIp"))
                if (!(message.oldIp && typeof message.oldIp.length === "number" || $util.isString(message.oldIp)))
                    return "oldIp: buffer expected";
            if (message.newIp != null && message.hasOwnProperty("newIp"))
                if (!(message.newIp && typeof message.newIp.length === "number" || $util.isString(message.newIp)))
                    return "newIp: buffer expected";
            if (message.ts != null && message.hasOwnProperty("ts"))
                if (!$util.isInteger(message.ts) && !(message.ts && $util.isInteger(message.ts.low) && $util.isInteger(message.ts.high)))
                    return "ts: integer|Long expected";
            return null;
        };

        /**
         * Creates a P_DISPLACE message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.P_DISPLACE
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.P_DISPLACE} P_DISPLACE
         */
        P_DISPLACE.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.P_DISPLACE)
                return object;
            var message = new $root.ws.P_DISPLACE();
            if (object.oldIp != null)
                if (typeof object.oldIp === "string")
                    $util.base64.decode(object.oldIp, message.oldIp = $util.newBuffer($util.base64.length(object.oldIp)), 0);
                else if (object.oldIp.length >= 0)
                    message.oldIp = object.oldIp;
            if (object.newIp != null)
                if (typeof object.newIp === "string")
                    $util.base64.decode(object.newIp, message.newIp = $util.newBuffer($util.base64.length(object.newIp)), 0);
                else if (object.newIp.length >= 0)
                    message.newIp = object.newIp;
            if (object.ts != null)
                if ($util.Long)
                    (message.ts = $util.Long.fromValue(object.ts)).unsigned = false;
                else if (typeof object.ts === "string")
                    message.ts = parseInt(object.ts, 10);
                else if (typeof object.ts === "number")
                    message.ts = object.ts;
                else if (typeof object.ts === "object")
                    message.ts = new $util.LongBits(object.ts.low >>> 0, object.ts.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a P_DISPLACE message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.P_DISPLACE
         * @static
         * @param {ws.P_DISPLACE} message P_DISPLACE
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        P_DISPLACE.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                if (options.bytes === String)
                    object.oldIp = "";
                else {
                    object.oldIp = [];
                    if (options.bytes !== Array)
                        object.oldIp = $util.newBuffer(object.oldIp);
                }
                if (options.bytes === String)
                    object.newIp = "";
                else {
                    object.newIp = [];
                    if (options.bytes !== Array)
                        object.newIp = $util.newBuffer(object.newIp);
                }
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.ts = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.ts = options.longs === String ? "0" : 0;
            }
            if (message.oldIp != null && message.hasOwnProperty("oldIp"))
                object.oldIp = options.bytes === String ? $util.base64.encode(message.oldIp, 0, message.oldIp.length) : options.bytes === Array ? Array.prototype.slice.call(message.oldIp) : message.oldIp;
            if (message.newIp != null && message.hasOwnProperty("newIp"))
                object.newIp = options.bytes === String ? $util.base64.encode(message.newIp, 0, message.newIp.length) : options.bytes === Array ? Array.prototype.slice.call(message.newIp) : message.newIp;
            if (message.ts != null && message.hasOwnProperty("ts"))
                if (typeof message.ts === "number")
                    object.ts = options.longs === String ? String(message.ts) : message.ts;
                else
                    object.ts = options.longs === String ? $util.Long.prototype.toString.call(message.ts) : options.longs === Number ? new $util.LongBits(message.ts.low >>> 0, message.ts.high >>> 0).toNumber() : message.ts;
            return object;
        };

        /**
         * Converts this P_DISPLACE to JSON.
         * @function toJSON
         * @memberof ws.P_DISPLACE
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        P_DISPLACE.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        /**
         * Gets the default type url for P_DISPLACE
         * @function getTypeUrl
         * @memberof ws.P_DISPLACE
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        P_DISPLACE.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/ws.P_DISPLACE";
        };

        return P_DISPLACE;
    })();

    /**
     * P_BASE enum.
     * @name ws.P_BASE
     * @enum {number}
     * @property {number} none=0 none value
     * @property {number} s2c_err_displace=2147483647 s2c_err_displace value
     */
    ws.P_BASE = (function() {
        var valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "none"] = 0;
        values[valuesById[2147483647] = "s2c_err_displace"] = 2147483647;
        return values;
    })();

    return ws;
})();

module.exports = $root;
