package token

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
)

type Token struct {
	Id   uint
	Name string
}

type TokenCodec struct {
	key []byte
	sep []byte
}

func NewTokenCodec(key, sep []byte) TokenCodec {
	return TokenCodec{key, sep}
}

func (tc TokenCodec) Bytes(id uint, name string) []byte {
	var (
		token   []byte
		payload []byte
	)
	token = bytes.Join([][]byte{
		[]byte(strconv.FormatInt(int64(id), 10)),
		[]byte(name),
	}, tc.sep)

	payload = append(token, tc.key...)
	sum := sign(payload)

	token = append(token, tc.sep...)
	token = append(token, sum[:]...)
	return token
}

func sign(payload []byte) (sum [32]byte) {
	bsum := md5.Sum(payload)
	hex.Encode(sum[:], bsum[:])
	return sum
}

func (tc TokenCodec) Token(src []byte) (Token, error) {
	var t Token
	info := bytes.Split(src, tc.sep)
	if len(info) != 3 {
		return t, &ErrInvalidToken{src, "invalid token length"}
	}
	var sum [32]byte
	copy(sum[:], info[2])
	payload := append(bytes.Join(info[:2], tc.sep), tc.key...)

	if sum != sign(payload) {
		return t, &ErrInvalidToken{src, "invalid token sum"}
	}
	id, err := strconv.ParseInt(string(info[0]), 10, 64)
	if err != nil {
		return t, &ErrInvalidToken{src, "invalid token id"}
	}
	t.Id = uint(id)
	t.Name = string(info[1])
	return t, nil
}

type ErrInvalidToken struct {
	Token []byte
	Msg   string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token %#v %s", e.Token, e.Msg)
}
