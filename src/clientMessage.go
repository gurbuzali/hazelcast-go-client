/*
Client Message is the carrier framed data as defined below.
Any request parameter, response or event data will be carried in the payload.

0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|R|                      Frame Length                           |
+-------------+---------------+---------------------------------+
|  Version    |B|E|  Flags    |               Type              |
+-------------+---------------+---------------------------------+
|                                                               |
+                       CorrelationId                           +
|                                                               |
+---------------------------------------------------------------+
|                        PartitionId                            |
+-----------------------------+---------------------------------+
|        Data Offset          |                                 |
+-----------------------------+                                 |
|                      Message Payload Data                    ...
|
 */

package src

import (
	"encoding/binary"
)

const (
	VERSION = 0
	BEGIN_FLAG = 0x80
	END_FLAG = 0x40
	BEGIN_END_FLAG = BEGIN_FLAG | END_FLAG
	LISTENER_FLAG = 0x01

	PAYLOAD_OFFSET = 18
	SIZE_OFFSET = 0

	FRAME_LENGTH_FIELD_OFFSET = 0
	VERSION_FIELD_OFFSET = FRAME_LENGTH_FIELD_OFFSET + INT_SIZE_IN_BYTES
	FLAGS_FIELD_OFFSET = VERSION_FIELD_OFFSET + BYTE_SIZE_IN_BYTES
	TYPE_FIELD_OFFSET = FLAGS_FIELD_OFFSET + BYTE_SIZE_IN_BYTES
	CORRELATION_ID_FIELD_OFFSET = TYPE_FIELD_OFFSET + SHORT_SIZE_IN_BYTES
	PARTITION_ID_FIELD_OFFSET = CORRELATION_ID_FIELD_OFFSET + LONG_SIZE_IN_BYTES
	DATA_OFFSET_FIELD_OFFSET = PARTITION_ID_FIELD_OFFSET + INT_SIZE_IN_BYTES
	HEADER_SIZE = DATA_OFFSET_FIELD_OFFSET + SHORT_SIZE_IN_BYTES
)

type ClientMessage struct {
	Buffer      []byte
	WriteIndex  int
	ReadIndex   int
	IsRetryable bool
}

//Todo
func (msg *ClientMessage) NewClientMessage(buffer []byte, payloadSize int) *ClientMessage {
	if buffer {
		msg.Buffer = buffer
		msg.ReadIndex = 0
	}else {
		buffer = make([]byte, HEADER_SIZE + payloadSize)
		msg.SetDataOffset(HEADER_SIZE)
		msg.WriteIndex = 0
		return &ClientMessage{buffer}
	}
	msg.IsRetryable = false

	return msg
}

func (msg *ClientMessage) GetFrameLength() int32 {
	return binary.LittleEndian.Uint32(msg.Buffer[FRAME_LENGTH_FIELD_OFFSET:VERSION_FIELD_OFFSET])
}

func (msg *ClientMessage) SetFrameLength(v int32) {
	binary.LittleEndian.PutUint32(msg.Buffer[FRAME_LENGTH_FIELD_OFFSET:VERSION_FIELD_OFFSET],uint32(v))
}

func (msg *ClientMessage) SetVersion(v uint8) {
	msg.Buffer[VERSION_FIELD_OFFSET] = byte(v)
}

func (msg *ClientMessage) GetFlags() uint8{
	return msg.Buffer[FLAGS_FIELD_OFFSET]
}

func (msg *ClientMessage) SetFlags(v uint8) {
	msg.Buffer[FLAGS_FIELD_OFFSET] = byte(v)
}

func (msg *ClientMessage) HasFlags(flags uint8) uint8 {
	return msg.GetFlags() & flags
}

func (msg *ClientMessage) GetMessageType() uint16 {
	return binary.LittleEndian.Uint16(msg.Buffer[TYPE_FIELD_OFFSET:CORRELATION_ID_FIELD_OFFSET])
}

func (msg *ClientMessage) SetMessageType(v uint16) {
	binary.LittleEndian.PutUint16(msg.Buffer, v)
}

func (msg *ClientMessage) GetCorrelationId() int64 {
	return int64(binary.LittleEndian.Uint64(msg.Buffer[CORRELATION_ID_FIELD_OFFSET:PARTITION_ID_FIELD_OFFSET]))
}

func (msg *ClientMessage) SetCorrelationId(val uint64) {
	binary.LittleEndian.PutUint64(msg.Buffer[CORRELATION_ID_FIELD_OFFSET:PARTITION_ID_FIELD_OFFSET], uint64(val))
}

func (msg *ClientMessage) GetPartitionId() int32 {
	return int32(binary.LittleEndian.Uint32(msg.Buffer[PARTITION_ID_FIELD_OFFSET:DATA_OFFSET_FIELD_OFFSET]))
}

func (msg *ClientMessage) SetPartitionId(val int32) {
	binary.LittleEndian.PutUint32(msg.Buffer[PARTITION_ID_FIELD_OFFSET:DATA_OFFSET_FIELD_OFFSET], uint32(val))
}

func (msg *ClientMessage) GetDataOffset() uint16 {
	return binary.LittleEndian.Uint16(msg.Buffer[DATA_OFFSET_FIELD_OFFSET:HEADER_SIZE])
}

func (msg *ClientMessage) SetDataOffset(v uint16) {
	binary.LittleEndian.Uint16(msg.Buffer[DATA_OFFSET_FIELD_OFFSET:HEADER_SIZE])
}

func (msg *ClientMessage) SetIsRetryable(v bool) {
	msg.IsRetryable = v
}