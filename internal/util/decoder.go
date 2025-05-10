package util

import (
    "encoding/base64"
    "encoding/binary"
    "errors"
)

func ParseSPLInstruction(data string) (uint64, error) {
    raw, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return 0, err
    }
    if len(raw) < 9 || (raw[0] != 3 && raw[0] != 12) {
        return 0, errors.New("unsupported instruction")
    }
    return binary.LittleEndian.Uint64(raw[1:9]), nil
}
