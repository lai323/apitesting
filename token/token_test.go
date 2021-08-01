package token

import (
	"reflect"
	"testing"
)

func TestTokenCodecBytes(t *testing.T) {
	codec := NewTokenCodec([]byte("AQdJz7ZGbv8GepLS"), []byte{'_'})
	token := Token{
		Id:   1,
		Name: "name",
	}
	rightBytes := []byte("1_name_2413a4da604906c6c9b4145de7449071")

	bytes := codec.Bytes(token.Id, token.Name)
	if !reflect.DeepEqual(bytes, rightBytes) {
		panic("wrong token")
	}
}

func TestTokenCodecToken(t *testing.T) {
	codec := NewTokenCodec([]byte("AQdJz7ZGbv8GepLS"), []byte{'_'})
	rightToken := Token{
		Id:   1,
		Name: "name",
	}
	tokenbytes := []byte("1_name_2413a4da604906c6c9b4145de7449071")
	token, err := codec.Token(tokenbytes)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(rightToken, token) {
		panic("wrong token")
	}
}
