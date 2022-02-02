package utilities

import (
	"errors"
	"fmt"

	"time"

	"github.com/aead/chacha20poly1305"
	
	"github.com/o1egl/paseto"
)


type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}
type Payload struct {
	Username  string    `json:"username"`
	AccoutType string   `json:"accountType"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
type Maker interface {
    CreateToken(username string, accountType string, duration time.Duration) (string, error)
    VerifyToken(token string) (*Payload, error)
}
var TokenMaker,_ = NewPasetoMaker("rthfrutg7tjgmbh4ryfkg7rhf7rjf74h") // secrete must be 32 bit char

var  (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid"))

func NewPayload(username string, accountType string) (*Payload, error) {
    payload := &Payload{
        Username:  username,
		AccoutType: accountType,
        IssuedAt:  time.Now(),
        ExpiredAt: time.Now().Add(time.Minute * 60),
    }
    return payload, nil
}
func (payload *Payload) Valid()  error{
	if time.Now().After(payload.ExpiredAt){
		return ErrExpiredToken
	}
	return nil
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
    if len(symmetricKey) != chacha20poly1305.KeySize {
        return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
    }
	maker := &PasetoMaker{
        paseto:       paseto.NewV2(),
        symmetricKey: []byte(symmetricKey),
    }

    return maker, nil
}
func (maker *PasetoMaker) CreateToken(username string, accountType string, duration time.Duration) (string, error) {
    payload, err := NewPayload(username, accountType)
    if err != nil {
        return "", err
    }

    return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
    payload := &Payload{}

    err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
    if err != nil {
        return nil, ErrInvalidToken
    }

    err = payload.Valid()
    if err != nil {
        return nil, err
    }

    return payload, nil
}

  
