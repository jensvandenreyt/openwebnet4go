package communication

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// DigitToHex converts a [0-9] digit string to hex string.
func DigitToHex(digits string) string {
	out := ""
	chars := []rune(digits)
	for i := 0; i+3 < len(chars); i += 4 {
		v1 := charToInt(chars[i])*10 + charToInt(chars[i+1])
		v2 := charToInt(chars[i+2])*10 + charToInt(chars[i+3])
		out += fmt.Sprintf("%x%x", v1, v2)
	}
	return out
}

// HexToDigit converts a hex string to [0-9] digit string.
func HexToDigit(hexString string) string {
	out := ""
	for _, c := range hexString {
		v, _ := strconv.ParseInt(string(c), 16, 64)
		out += fmt.Sprintf("%d%d", v/10, v%10)
	}
	return out
}

// BytesToHex converts a byte array to hex string.
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// CalcHmacRb generates Rb HMAC random hex string using SHA-256.
func CalcHmacRb() string {
	return CalcSHA256(fmt.Sprintf("time%d", time.Now().UnixMilli()))
}

// CalcSHA256 returns the SHA-256 hash of the input string.
func CalcSHA256(message string) string {
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:])
}

// CalcOpenPass encodes a numeric OPEN password.
func CalcOpenPass(pass string, nonce string) (string, error) {
	password, err := strconv.ParseUint(pass, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid password: %v", err)
	}

	flag := true
	var num1 uint32
	var num2 uint32

	for _, c := range nonce {
		if c != '0' {
			if flag {
				num2 = uint32(password)
			}
			flag = false
		}
		switch c {
		case '1':
			num1 = num2 & 0xFFFFFF80
			num1 = num1 >> 7
			num2 = num2 << 25
			num1 = num1 + num2
		case '2':
			num1 = num2 & 0xFFFFFFF0
			num1 = num1 >> 4
			num2 = num2 << 28
			num1 = num1 + num2
		case '3':
			num1 = num2 & 0xFFFFFFF8
			num1 = num1 >> 3
			num2 = num2 << 29
			num1 = num1 + num2
		case '4':
			num1 = num2 << 1
			num2 = num2 >> 31
			num1 = num1 + num2
		case '5':
			num1 = num2 << 5
			num2 = num2 >> 27
			num1 = num1 + num2
		case '6':
			num1 = num2 << 12
			num2 = num2 >> 20
			num1 = num1 + num2
		case '7':
			num1 = num2 & 0x0000FF00
			num1 = num1 + ((num2 & 0x000000FF) << 24)
			num1 = num1 + ((num2 & 0x00FF0000) >> 16)
			num2 = (num2 & 0xFF000000) >> 8
			num1 = num1 + num2
		case '8':
			num1 = num2 & 0x0000FFFF
			num1 = num1 << 16
			num1 = num1 + (num2 >> 24)
			num2 = num2 & 0x00FF0000
			num2 = num2 >> 8
			num1 = num1 + num2
		case '9':
			num1 = ^num2
		case '0':
			num1 = num2
		}
		num2 = num1
	}
	return strconv.FormatUint(uint64(num1), 10), nil
}

func charToInt(c rune) int {
	return int(c - '0')
}
