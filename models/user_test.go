package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	public, err := x509.MarshalPKIXPublicKey(priv.Public())
	assert.Nil(err)
	user, err := CreateUser(ctx, "im.yuqlee@gmailabcefgh.com", "username", "nickname", "password", hex.EncodeToString(public))
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "pass", hex.EncodeToString(public))
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "    pass     ", hex.EncodeToString(public))
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "password", hex.EncodeToString(public))
	assert.Nil(err)
	assert.NotNil(user)
	assert.NotEqual("", user.SessionId)
	new, err := FindUser(ctx, user.UserId)
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal(user.Username, new.Username)
	assert.Equal(user.Nickname, new.Nickname)
	err = bcrypt.CompareHashAndPassword([]byte(new.EncryptedPassword), []byte("password"))
	assert.Nil(err)
	new, err = FindUserByUsernameOrEmail(ctx, "None")
	assert.Nil(err)
	assert.Nil(new)
	new, err = FindUserByUsernameOrEmail(ctx, "im.yuqlee@Gmail.com")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = FindUserByUsernameOrEmail(ctx, "UserName")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = FindUserByUsernameOrEmail(ctx, "im.yuqlee@Gmail.com")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = CreateSession(ctx, "im.yuqlee@Gmail.com", "password", hex.EncodeToString(public))
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal("username", user.Username)

	sess, err := readSession(ctx, new.UserId, new.SessionId)
	assert.Nil(err)
	assert.NotNil(sess)
	sess, err = readSession(ctx, uuid.NewV4().String(), new.SessionId)
	assert.Nil(err)
	assert.Nil(sess)

	claims := &jwt.MapClaims{
		"uid": new.UserId,
		"sid": new.SessionId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(priv)
	assert.Nil(err)
	new, err = AuthenticateUser(ctx, ss)
	assert.Nil(err)
	assert.NotNil(new)
}