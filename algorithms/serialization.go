package algorithms

import (
	"encoding/binary"
	"errors"
	"io"
)

// Serialization - Encoding/decoding for blockchain wire protocols
// Efficient serialization is critical for network bandwidth and storage

// ------------------------------------------------------------
// RLP (Recursive Length Prefix) - Ethereum's encoding
// Used for: Transactions, blocks, state trie nodes, wire protocol
// Spec: https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/
// ------------------------------------------------------------

// RLP can encode:
// - Strings (byte arrays) of length 0-55: 0x80 + len, then data
// - Strings > 55 bytes: 0xb7 + len(len), then len, then data
// - Lists of length 0-55: 0xc0 + len, then items
// - Lists > 55 bytes: 0xf7 + len(len), then len, then items
// - Single byte 0x00-0x7f: byte itself (no prefix)

// RLPEncode encodes a value to RLP format
// value can be: []byte, string, []interface{} (for lists), or uint64
func RLPEncode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return RLPEncodeString(v), nil
	case string:
		return RLPEncodeString([]byte(v)), nil
	case uint64:
		return RLPEncodeUint64(v), nil
	case []interface{}:
		// Recursively encode each item, then encode as list
		encodedItems := make([][]byte, len(v))
		for i, item := range v {
			encoded, err := RLPEncode(item)
			if err != nil {
				return nil, err
			}
			encodedItems[i] = encoded
		}
		return RLPEncodeList(encodedItems), nil
	case nil:
		return []byte{0x80}, nil // empty string
	default:
		return nil, errors.New("unsupported type for RLP encoding")
	}
}

// RLPEncodeString encodes a byte slice as RLP string
func RLPEncodeString(data []byte) []byte {
	dataLen := len(data)
	if dataLen == 1 && data[0] <= 0x7f {
		// Single byte in range[0x00, 0x7f]: return as-is
		return data
	} else if dataLen <= 55 {
		// 0-55 bytes: 0x80 + len, then data
		prefix := []byte{0x80 + byte(dataLen)}
		return append(prefix, data...)
	} else {
		// > 55 bytes: 0xb7 + len(lenBytes), then lenBytes, then data
		// Encode length as minimal big-endian bytes without leading zeros
		lenBytes := encodeBigEndian(uint64(dataLen)) // e.g., 256 -> [0x01, 0x00]
		prefix := append([]byte{0xb7 + byte(len(lenBytes))}, lenBytes...)
		return append(prefix, data...) // data stays intact
	}
}

// RLPEncodeList encodes a list of RLP-encoded items
func RLPEncodeList(items [][]byte) []byte {
	// Concatenate all items as flat bytes
	result := []byte{}
	for _, arr := range items {
		for _, b := range arr {
			result = append(result, b)
		}
	}

	// Then prefix with list header
	payloadLen := len(result)
	var prefix []byte
	if payloadLen <= 55 {
		// 0-55 bytes total: 0xc0 + len, then items
		prefix = []byte{0xc0 + byte(payloadLen)}
	} else {
		// > 55 bytes: 0xf7 + len(lenBytes), then lenBytes, then items
		lenBytes := encodeBigEndian(uint64(payloadLen))
		prefix = []byte{0xf7 + byte(len(lenBytes))}
		prefix = append(prefix, lenBytes...)
	}
	return append(prefix, result...)
}

// RLPEncodeUint64 encodes a uint64 as RLP
func RLPEncodeUint64(n uint64) []byte {
	if n == 0 {
		// 0: empty string (0x80)
		return []byte{0x80}
	} else if n <= 127 {
		// 1-127: single byte
		return []byte{byte(n)}
	} else {
		// > 127: encode as big-endian bytes without leading zeros
		// 256 should encode as [0x82, 0x01, 0x00] (prefix + 2 bytes)
		nBytes := encodeBigEndian(n)
		return RLPEncodeString(nBytes)
	}
}

// RLPDecode decodes RLP data into a value
// Returns one of: []byte (string), []interface{} (list)
func RLPDecode(data []byte) (interface{}, error) {
	result, _, err := rlpDecodeItem(data, 0)
	return result, err
}

// Converts an integer to the minimum bytes needed in big-endian order.
// Input: 56    → Output: [0x38]           (1 byte)
// Input: 255   → Output: [0xff]           (1 byte)
// Input: 256   → Output: [0x01, 0x00]     (2 bytes)
// Input: 1024  → Output: [0x04, 0x00]     (2 bytes)
// Input: 65536 → Output: [0x01, 0x00, 0x00] (3 bytes)
func encodeBigEndian(n uint64) []byte {
	// Put into 8-byte buffer
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)

	// Find first non-zero byte (strip leading zeros)
	i := 0
	for i < 7 && buf[i] == 0 {
		i++
	}
	return buf[i:]
}

// Takes a variable-length byte slice and returns the integer
// Input: [0x38]               → Output: 56
// Input: [0x01, 0x00]         → Output: 256
// Input: [0x01, 0x00, 0x00]   → Output: 65536
func decodeBigEndian(lenBytes []byte) int {
	// Pad to 8 bytes on the left
	buf := make([]byte, 8)
	copy(buf[8-len(lenBytes):], lenBytes)
	return int(binary.BigEndian.Uint64(buf))
}

// rlpDecodeItem decodes one RLP item starting at offset
// Returns: decoded value, bytes consumed, error
func rlpDecodeItem(data []byte, offset int) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("Empty data")
	} else if len(data) <= offset {
		return nil, 0, errors.New("Index out of bounds")
	}
	// Check first byte to determine type:
	key := data[offset]
	// Empty string: [ 0x80 ] -> null
	if key == 0x80 {
		return "", 1, nil
	}
	if key == 0xc0 {
		// Empty list
		return []byte{}, 1, nil
	}
	// Range [0x00, 0x7f]: single byte
	if key <= 0x7f {
		return key, 1, nil
	}
	// Range [0x80, 0xb7]: short string
	if key <= 0xb7 {
		dataLen := int(key - 0x80)
		return data[offset+1 : offset+1+dataLen], dataLen + 1, nil
	}
	// Range [0xb8, 0xbf]: long string
	if key <= 0xbf {
		// Convert the length from byte[] to int
		numLenBytes := int(key - 0xb7)
		dataLen := decodeBigEndian(data[offset+1 : offset+1+numLenBytes])
		dataStart := offset + 1 + numLenBytes
		totalConsumed := 1 + numLenBytes + dataLen
		return data[dataStart : dataStart+dataLen], totalConsumed, nil
	}
	// Range [0xc0, 0xf7]: short list
	if key <= 0xf7 {
		payloadLen := int(key - 0xc0)
		items := []interface{}{}
		pos := offset + 1
		endPos := offset + 1 + payloadLen
		for pos < endPos {
			item, consumed, err := rlpDecodeItem(data, pos)
			if err != nil {
				return nil, 0, err
			}
			items = append(items, item)
			pos += consumed
		}
		return items, 1 + payloadLen, nil
	}
	// Range [0xf8, 0xff]: long list
	numLenBytes := int(key - 0xf7)
	payloadLen := decodeBigEndian(data[offset+1 : offset+1+numLenBytes])
	items := []interface{}{}
	pos := offset + 1 + numLenBytes // Start of payload
	endPos := pos + payloadLen      // End of payload
	for pos < endPos {
		item, consumed, err := rlpDecodeItem(data, pos)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
		pos += consumed
	}
	return items, 1 + numLenBytes + payloadLen, nil
}

// RLPReader provides streaming RLP decoding
type RLPReader struct {
	r io.Reader
}

// NewRLPReader creates an RLP reader
func NewRLPReader(r io.Reader) *RLPReader {
	return &RLPReader{r: r}
}

// ReadString reads the next RLP string
func (r *RLPReader) ReadString() ([]byte, error) {
	// Read prefix byte
	prefix := make([]byte, 1)
	if _, err := io.ReadFull(r.r, prefix); err != nil {
		return nil, err
	}
	key := prefix[0]

	// Single byte [0x00, 0x7f]
	if key <= 0x7f {
		return prefix, nil
	}
	// Empty string
	if key == 0x80 {
		return []byte{}, nil
	}
	// Short string [0x81, 0xb7]
	if key <= 0xb7 {
		dataLen := int(key - 0x80)
		data := make([]byte, dataLen)
		if _, err := io.ReadFull(r.r, data); err != nil {
			return nil, err
		}
		return data, nil
	}
	// Long string [0xb8, 0xbf]
	if key <= 0xbf {
		numLenBytes := int(key - 0xb7)
		lenBytes := make([]byte, numLenBytes)
		if _, err := io.ReadFull(r.r, lenBytes); err != nil {
			return nil, err
		}
		dataLen := decodeBigEndian(lenBytes)
		data := make([]byte, dataLen)
		if _, err := io.ReadFull(r.r, data); err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("expected string, got list prefix")
}

// ReadList reads the next RLP list, returning a reader for the list contents
func (r *RLPReader) ReadList() (*RLPReader, error) {
	// Read prefix byte
	prefix := make([]byte, 1)
	if _, err := io.ReadFull(r.r, prefix); err != nil {
		return nil, err
	}
	key := prefix[0]

	var payloadLen int

	if key == 0xc0 {
		// Empty list
		payloadLen = 0
	} else if key > 0xc0 && key <= 0xf7 {
		// Short list [0xc1, 0xf7]
		payloadLen = int(key - 0xc0)
	} else if key > 0xf7 {
		// Long list [0xf8, 0xff]
		numLenBytes := int(key - 0xf7)
		lenBytes := make([]byte, numLenBytes)
		if _, err := io.ReadFull(r.r, lenBytes); err != nil {
			return nil, err
		}
		payloadLen = decodeBigEndian(lenBytes)
	} else {
		return nil, errors.New("expected list, got string prefix")
	}

	// Return a reader limited to the list payload
	return &RLPReader{r: io.LimitReader(r.r, int64(payloadLen))}, nil
}

// ------------------------------------------------------------
// SSZ (Simple Serialize) - Ethereum 2.0 encoding
// Used for: Beacon chain blocks, attestations, state
// Spec: https://ethereum.org/developers/docs/data-structures-and-encoding/ssz/
// 		 https://www.ssz.dev/overview
// ------------------------------------------------------------

// SSZ features:
// - Fixed-size types: encoded as-is (little-endian for integers)
// - Variable-size types: offset + data pattern
// - Merkleization: SSZ data can be efficiently hashed into Merkle roots

// SSZEncoder handles SSZ encoding
type SSZEncoder struct {
	buf []byte
}

// NewSSZEncoder creates a new SSZ encoder
func NewSSZEncoder() *SSZEncoder {
	return &SSZEncoder{buf: make([]byte, 0)}
}

// EncodeUint8 encodes a uint8
func (e *SSZEncoder) EncodeUint8(v uint8) {
	// TODO: Implement (1 byte)
}

// EncodeUint16 encodes a uint16 (little-endian)
func (e *SSZEncoder) EncodeUint16(v uint16) {
	// TODO: Implement (2 bytes, little-endian)
}

// EncodeUint32 encodes a uint32 (little-endian)
func (e *SSZEncoder) EncodeUint32(v uint32) {
	// TODO: Implement (4 bytes, little-endian)
}

// EncodeUint64 encodes a uint64 (little-endian)
func (e *SSZEncoder) EncodeUint64(v uint64) {
	// TODO: Implement (8 bytes, little-endian)
}

// EncodeBytes encodes a fixed-size byte array
func (e *SSZEncoder) EncodeBytes(data []byte) {
	// TODO: Implement
}

// EncodeVariableBytes encodes variable-length bytes (with offset)
func (e *SSZEncoder) EncodeVariableBytes(data []byte) {
	// TODO: Implement
	// Write 4-byte offset, then data goes in variable part
}

// Bytes returns the encoded bytes
func (e *SSZEncoder) Bytes() []byte {
	return e.buf
}

// SSZDecoder handles SSZ decoding
type SSZDecoder struct {
	data   []byte
	offset int
}

// NewSSZDecoder creates a new SSZ decoder
func NewSSZDecoder(data []byte) *SSZDecoder {
	return &SSZDecoder{data: data, offset: 0}
}

// DecodeUint8 decodes a uint8
func (d *SSZDecoder) DecodeUint8() (uint8, error) {
	// TODO: Implement
	return 0, errors.New("not implemented")
}

// DecodeUint64 decodes a uint64 (little-endian)
func (d *SSZDecoder) DecodeUint64() (uint64, error) {
	// TODO: Implement
	return 0, errors.New("not implemented")
}

// DecodeBytes decodes a fixed-size byte array
func (d *SSZDecoder) DecodeBytes(length int) ([]byte, error) {
	// TODO: Implement
	return nil, errors.New("not implemented")
}

// ------------------------------------------------------------
// SSZ Merkleization
// Compute Merkle roots from SSZ data for light client proofs
// ------------------------------------------------------------

// SSZHashTreeRoot computes the hash tree root of SSZ data
// Used for: Block roots, state roots, signing roots
func SSZHashTreeRoot(data []byte) [32]byte {
	// TODO: Implement
	// 1. Split data into 32-byte chunks (pad last chunk if needed)
	// 2. Build Merkle tree from chunks
	// 3. Return root
	return [32]byte{}
}

// SSZMerkleize computes merkle root from a list of 32-byte leaves
func SSZMerkleize(leaves [][32]byte) [32]byte {
	// TODO: Implement
	// 1. If odd number of leaves, add zero leaf
	// 2. Hash pairs: parent = hash(left || right)
	// 3. Repeat until one root remains
	return [32]byte{}
}

// SSZMerkleProof generates a Merkle proof for a leaf at index
func SSZMerkleProof(leaves [][32]byte, index int) ([][32]byte, error) {
	// TODO: Implement
	// Return sibling hashes from leaf to root
	return nil, errors.New("not implemented")
}

// SSZVerifyProof verifies a Merkle proof
func SSZVerifyProof(root [32]byte, leaf [32]byte, proof [][32]byte, index int) bool {
	// TODO: Implement
	// Recompute root from leaf and proof, compare to expected root
	return false
}

// ------------------------------------------------------------
// Protocol Buffers (simplified concepts)
// Used by: gRPC, libp2p, many blockchain projects
// ------------------------------------------------------------

// Protobuf uses varints and wire types:
// - Varint: variable-length integer encoding
// - Wire types: 0=varint, 1=64-bit, 2=length-delimited, 5=32-bit

// ProtobufEncodeVarint encodes a uint64 as a varint
func ProtobufEncodeVarint(n uint64) []byte {
	// TODO: Implement
	// Each byte: 7 bits of data + 1 continuation bit
	// MSB = 1 means more bytes follow
	// Example: 300 = 0b100101100 -> 0xAC 0x02
	return nil
}

// ProtobufDecodeVarint decodes a varint from bytes
func ProtobufDecodeVarint(data []byte) (uint64, int, error) {
	// TODO: Implement
	// Returns: value, bytes consumed, error
	return 0, 0, errors.New("not implemented")
}

// ProtobufEncodeField encodes a protobuf field
// fieldNum: field number (1-based)
// wireType: 0=varint, 2=length-delimited
// data: field data
func ProtobufEncodeField(fieldNum int, wireType int, data []byte) []byte {
	// TODO: Implement
	// Tag = (fieldNum << 3) | wireType
	// For length-delimited: tag, length (varint), data
	return nil
}

// ------------------------------------------------------------
// Benchmarking Utilities
// Compare serialization performance
// ------------------------------------------------------------

// BenchmarkData generates test data of specified size
func BenchmarkData(size int) []byte {
	// TODO: Implement
	return make([]byte, size)
}

// MeasureEncodingSize compares encoding sizes
// Returns: RLP size, SSZ size, JSON size (for reference)
func MeasureEncodingSize(data interface{}) (rlpSize, sszSize, jsonSize int) {
	// TODO: Implement
	return 0, 0, 0
}

// Serialization concepts to understand:
//
// 1. RLP (Ethereum)
//    - Self-describing (can decode without schema)
//    - Compact for small data
//    - Used in MPT nodes, transactions, blocks
//
// 2. SSZ (Ethereum 2.0)
//    - Fixed offsets enable O(1) field access
//    - Native Merkleization for proofs
//    - Required for consensus layer
//
// 3. Protobuf
//    - Requires schema (.proto files)
//    - Efficient and widely supported
//    - Used by Cosmos SDK, libp2p
//
// 4. Performance considerations:
//    - Encoding/decoding speed
//    - Encoded size (bandwidth)
//    - Random access capability
//    - Proof generation support
