package algorithms

import (
	"bytes"
	"testing"
)

// ========== encodeBigEndian / decodeBigEndian Tests ==========

func TestEncodeBigEndian(t *testing.T) {
	tests := []struct {
		input    uint64
		expected []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{56, []byte{0x38}},
		{127, []byte{0x7f}},
		{128, []byte{0x80}},
		{255, []byte{0xff}},
		{256, []byte{0x01, 0x00}},
		{1024, []byte{0x04, 0x00}},
		{65535, []byte{0xff, 0xff}},
		{65536, []byte{0x01, 0x00, 0x00}},
	}

	for _, tc := range tests {
		got := encodeBigEndian(tc.input)
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("encodeBigEndian(%d) = %x, want %x", tc.input, got, tc.expected)
		}
	}
}

func TestDecodeBigEndian(t *testing.T) {
	tests := []struct {
		input    []byte
		expected int
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x38}, 56},
		{[]byte{0x7f}, 127},
		{[]byte{0x80}, 128},
		{[]byte{0xff}, 255},
		{[]byte{0x01, 0x00}, 256},
		{[]byte{0x04, 0x00}, 1024},
		{[]byte{0xff, 0xff}, 65535},
		{[]byte{0x01, 0x00, 0x00}, 65536},
	}

	for _, tc := range tests {
		got := decodeBigEndian(tc.input)
		if got != tc.expected {
			t.Errorf("decodeBigEndian(%x) = %d, want %d", tc.input, got, tc.expected)
		}
	}
}

func TestBigEndianRoundTrip(t *testing.T) {
	values := []uint64{0, 1, 127, 128, 255, 256, 1000, 65535, 65536, 1000000}
	for _, v := range values {
		encoded := encodeBigEndian(v)
		decoded := decodeBigEndian(encoded)
		if uint64(decoded) != v {
			t.Errorf("Round trip failed for %d: encoded=%x, decoded=%d", v, encoded, decoded)
		}
	}
}

// ========== RLPEncodeString Tests ==========

func TestRLPEncodeString_SingleByte(t *testing.T) {
	// Single bytes 0x00-0x7f encode as themselves
	tests := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{0x00}, []byte{0x00}},
		{[]byte{0x01}, []byte{0x01}},
		{[]byte{0x7f}, []byte{0x7f}},
	}

	for _, tc := range tests {
		got := RLPEncodeString(tc.input)
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("RLPEncodeString(%x) = %x, want %x", tc.input, got, tc.expected)
		}
	}
}

func TestRLPEncodeString_SingleByteHighValue(t *testing.T) {
	// Single byte >= 0x80 needs a prefix
	input := []byte{0x80}
	got := RLPEncodeString(input)
	expected := []byte{0x81, 0x80} // 0x80 + 1 = 0x81, then data

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeString(%x) = %x, want %x", input, got, expected)
	}
}

func TestRLPEncodeString_Empty(t *testing.T) {
	// Empty string encodes as 0x80
	got := RLPEncodeString([]byte{})
	expected := []byte{0x80}

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeString(empty) = %x, want %x", got, expected)
	}
}

func TestRLPEncodeString_Short(t *testing.T) {
	// "dog" = [0x64, 0x6f, 0x67] -> 0x83 + data
	input := []byte("dog")
	got := RLPEncodeString(input)
	expected := []byte{0x83, 0x64, 0x6f, 0x67}

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeString(%q) = %x, want %x", input, got, expected)
	}
}

func TestRLPEncodeString_55Bytes(t *testing.T) {
	// 55 bytes: prefix = 0x80 + 55 = 0xb7
	input := bytes.Repeat([]byte{0x41}, 55)
	got := RLPEncodeString(input)

	if got[0] != 0xb7 {
		t.Errorf("55-byte string prefix = %x, want 0xb7", got[0])
	}
	if len(got) != 56 {
		t.Errorf("55-byte string encoded length = %d, want 56", len(got))
	}
}

func TestRLPEncodeString_Long(t *testing.T) {
	// 56 bytes: prefix = 0xb8 (0xb7 + 1), then length byte 0x38 (56)
	input := bytes.Repeat([]byte{0x41}, 56)
	got := RLPEncodeString(input)

	if got[0] != 0xb8 {
		t.Errorf("56-byte string prefix = %x, want 0xb8", got[0])
	}
	if got[1] != 0x38 {
		t.Errorf("56-byte string length byte = %x, want 0x38", got[1])
	}
	if len(got) != 58 { // 1 prefix + 1 length + 56 data
		t.Errorf("56-byte string encoded length = %d, want 58", len(got))
	}
}

func TestRLPEncodeString_VeryLong(t *testing.T) {
	// 1024 bytes: length = 0x0400, needs 2 length bytes
	// prefix = 0xb7 + 2 = 0xb9
	input := bytes.Repeat([]byte{0x41}, 1024)
	got := RLPEncodeString(input)

	if got[0] != 0xb9 {
		t.Errorf("1024-byte string prefix = %x, want 0xb9", got[0])
	}
	if got[1] != 0x04 || got[2] != 0x00 {
		t.Errorf("1024-byte string length bytes = %x %x, want 04 00", got[1], got[2])
	}
	if len(got) != 1027 { // 1 prefix + 2 length + 1024 data
		t.Errorf("1024-byte string encoded length = %d, want 1027", len(got))
	}
}

// ========== RLPEncodeUint64 Tests ==========

func TestRLPEncodeUint64_Zero(t *testing.T) {
	got := RLPEncodeUint64(0)
	expected := []byte{0x80} // Empty string

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeUint64(0) = %x, want %x", got, expected)
	}
}

func TestRLPEncodeUint64_SmallValues(t *testing.T) {
	tests := []struct {
		input    uint64
		expected []byte
	}{
		{1, []byte{0x01}},
		{127, []byte{0x7f}},
	}

	for _, tc := range tests {
		got := RLPEncodeUint64(tc.input)
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("RLPEncodeUint64(%d) = %x, want %x", tc.input, got, tc.expected)
		}
	}
}

func TestRLPEncodeUint64_LargerValues(t *testing.T) {
	tests := []struct {
		input    uint64
		expected []byte
	}{
		{128, []byte{0x81, 0x80}},               // 1 byte
		{255, []byte{0x81, 0xff}},               // 1 byte
		{256, []byte{0x82, 0x01, 0x00}},         // 2 bytes
		{1024, []byte{0x82, 0x04, 0x00}},        // 2 bytes
		{65535, []byte{0x82, 0xff, 0xff}},       // 2 bytes
		{65536, []byte{0x83, 0x01, 0x00, 0x00}}, // 3 bytes
	}

	for _, tc := range tests {
		got := RLPEncodeUint64(tc.input)
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("RLPEncodeUint64(%d) = %x, want %x", tc.input, got, tc.expected)
		}
	}
}

// ========== RLPEncodeList Tests ==========

func TestRLPEncodeList_Empty(t *testing.T) {
	got := RLPEncodeList([][]byte{})
	expected := []byte{0xc0}

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeList(empty) = %x, want %x", got, expected)
	}
}

func TestRLPEncodeList_SingleItem(t *testing.T) {
	// List containing "dog"
	items := [][]byte{RLPEncodeString([]byte("dog"))}
	got := RLPEncodeList(items)
	// "dog" encodes to [0x83, 0x64, 0x6f, 0x67] = 4 bytes
	// list prefix = 0xc0 + 4 = 0xc4
	expected := []byte{0xc4, 0x83, 0x64, 0x6f, 0x67}

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeList([dog]) = %x, want %x", got, expected)
	}
}

func TestRLPEncodeList_MultipleItems(t *testing.T) {
	// List containing "cat" and "dog"
	items := [][]byte{
		RLPEncodeString([]byte("cat")),
		RLPEncodeString([]byte("dog")),
	}
	got := RLPEncodeList(items)
	// Each word encodes to 4 bytes, total payload = 8 bytes
	// list prefix = 0xc0 + 8 = 0xc8
	expected := []byte{0xc8, 0x83, 0x63, 0x61, 0x74, 0x83, 0x64, 0x6f, 0x67}

	if !bytes.Equal(got, expected) {
		t.Errorf("RLPEncodeList([cat, dog]) = %x, want %x", got, expected)
	}
}

func TestRLPEncodeList_NestedList(t *testing.T) {
	// [ "cat", [ "dog" ] ]
	innerList := RLPEncodeList([][]byte{RLPEncodeString([]byte("dog"))})
	items := [][]byte{
		RLPEncodeString([]byte("cat")),
		innerList,
	}
	got := RLPEncodeList(items)

	// "cat" = [0x83, 0x63, 0x61, 0x74] = 4 bytes
	// inner list = [0xc4, 0x83, 0x64, 0x6f, 0x67] = 5 bytes
	// total payload = 9 bytes
	// outer list prefix = 0xc0 + 9 = 0xc9
	if got[0] != 0xc9 {
		t.Errorf("Nested list prefix = %x, want 0xc9", got[0])
	}
	if len(got) != 10 {
		t.Errorf("Nested list length = %d, want 10", len(got))
	}
}

// ========== RLPDecode Tests ==========

func TestRLPDecode_EmptyString(t *testing.T) {
	input := []byte{0x80}
	got, err := RLPDecode(input)
	if err != nil {
		t.Fatalf("RLPDecode failed: %v", err)
	}

	if got != "" {
		t.Errorf("RLPDecode(0x80) = %v, want empty string", got)
	}
}

func TestRLPDecode_SingleByte(t *testing.T) {
	tests := []struct {
		input    []byte
		expected byte
	}{
		{[]byte{0x00}, 0x00},
		{[]byte{0x01}, 0x01},
		{[]byte{0x7f}, 0x7f},
	}

	for _, tc := range tests {
		got, err := RLPDecode(tc.input)
		if err != nil {
			t.Fatalf("RLPDecode(%x) failed: %v", tc.input, err)
		}

		if got != tc.expected {
			t.Errorf("RLPDecode(%x) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestRLPDecode_ShortString(t *testing.T) {
	// "dog" encoded
	input := []byte{0x83, 0x64, 0x6f, 0x67}
	got, err := RLPDecode(input)
	if err != nil {
		t.Fatalf("RLPDecode failed: %v", err)
	}

	gotBytes, ok := got.([]byte)
	if !ok {
		t.Fatalf("RLPDecode returned %T, want []byte", got)
	}

	if string(gotBytes) != "dog" {
		t.Errorf("RLPDecode = %q, want %q", gotBytes, "dog")
	}
}

func TestRLPDecode_EmptyList(t *testing.T) {
	input := []byte{0xc0}
	got, err := RLPDecode(input)
	if err != nil {
		t.Fatalf("RLPDecode failed: %v", err)
	}

	gotSlice, ok := got.([]byte)
	if !ok {
		t.Fatalf("RLPDecode returned %T, want []byte for empty list", got)
	}

	if len(gotSlice) != 0 {
		t.Errorf("RLPDecode(0xc0) = %v, want empty slice", gotSlice)
	}
}

func TestRLPDecode_ShortList(t *testing.T) {
	// List containing "cat" and "dog"
	input := []byte{0xc8, 0x83, 0x63, 0x61, 0x74, 0x83, 0x64, 0x6f, 0x67}
	got, err := RLPDecode(input)
	if err != nil {
		t.Fatalf("RLPDecode failed: %v", err)
	}

	gotList, ok := got.([]interface{})
	if !ok {
		t.Fatalf("RLPDecode returned %T, want []interface{}", got)
	}

	if len(gotList) != 2 {
		t.Fatalf("List length = %d, want 2", len(gotList))
	}

	cat, ok := gotList[0].([]byte)
	if !ok || string(cat) != "cat" {
		t.Errorf("First item = %v, want 'cat'", gotList[0])
	}

	dog, ok := gotList[1].([]byte)
	if !ok || string(dog) != "dog" {
		t.Errorf("Second item = %v, want 'dog'", gotList[1])
	}
}

func TestRLPDecode_Error_Empty(t *testing.T) {
	_, err := RLPDecode([]byte{})
	if err == nil {
		t.Error("RLPDecode(empty) should return error")
	}
}

// ========== Round-trip Tests ==========

func TestRLPRoundTrip_String(t *testing.T) {
	testCases := [][]byte{
		{},
		{0x00},
		{0x7f},
		{0x80},
		[]byte("hello"),
		[]byte("hello world this is a longer string"),
		bytes.Repeat([]byte{0x41}, 56),
		bytes.Repeat([]byte{0x42}, 1024),
	}

	for _, original := range testCases {
		encoded := RLPEncodeString(original)
		decoded, err := RLPDecode(encoded)
		if err != nil {
			t.Errorf("Round trip decode failed for %x: %v", original, err)
			continue
		}

		// Handle the single-byte case which returns byte, not []byte
		var decodedBytes []byte
		switch v := decoded.(type) {
		case []byte:
			decodedBytes = v
		case byte:
			decodedBytes = []byte{v}
		case string:
			decodedBytes = []byte(v)
		default:
			t.Errorf("Unexpected decoded type %T for input %x", decoded, original)
			continue
		}

		if !bytes.Equal(decodedBytes, original) {
			t.Errorf("Round trip failed: original=%x, encoded=%x, decoded=%x",
				original, encoded, decodedBytes)
		}
	}
}

func TestRLPRoundTrip_Uint64(t *testing.T) {
	values := []uint64{0, 1, 127, 128, 255, 256, 1024, 65535, 65536, 1000000}

	for _, original := range values {
		encoded := RLPEncodeUint64(original)
		decoded, err := RLPDecode(encoded)
		if err != nil {
			t.Errorf("Round trip decode failed for %d: %v", original, err)
			continue
		}

		// Convert decoded value back to uint64
		var decodedVal uint64
		switch v := decoded.(type) {
		case byte:
			decodedVal = uint64(v)
		case []byte:
			decodedVal = uint64(decodeBigEndian(v))
		case string:
			if v == "" {
				decodedVal = 0
			} else {
				t.Errorf("Unexpected string value for %d", original)
				continue
			}
		default:
			t.Errorf("Unexpected decoded type %T for %d", decoded, original)
			continue
		}

		if decodedVal != original {
			t.Errorf("Round trip failed: original=%d, decoded=%d", original, decodedVal)
		}
	}
}
