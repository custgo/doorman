package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"github.com/heiing/logs"
)

type Identity struct {
	username string
	password string
	created  time.Time
}

type TokenPool struct {
	tokens map[string]Identity
}

var expires int64 = 86400
var gcDuration = 30 * time.Second
var tokenPool = &TokenPool{
	tokens: make(map[string]Identity),
}

func (i Identity) isExpired() bool {
	return i.created.Unix()+expires < time.Now().Unix()
}

func (t *TokenPool) Add(username string, password string) string {
	identity := Identity{
		username: username,
		password: password,
		created:  time.Now(),
	}

	coder := sha1.New()
	io.WriteString(coder, username+":"+password+":"+fmt.Sprintf("%d", identity.created.UnixNano()))
	token := fmt.Sprintf("%x", coder.Sum(nil))

	t.tokens[token] = identity
	return token
}

func (t *TokenPool) Get(token string) *Identity {
	id, ok := t.tokens[token]
	if !ok || id.isExpired() {
		return nil
	}
	return &id
}

func (t *TokenPool) GC() {
	total := 0
	deleted := 0
	for k, v := range t.tokens {
		total++
		if v.isExpired() {
			delete(t.tokens, k)
			deleted++
			logs.Debug("Deleted token: ", k, " for ", v.username)
		}
	}
	logs.Info("Token Poll GC: deleted ", deleted, " of ", total)
}

func tokenPoolGC() {
	go func() {
		timer := time.NewTimer(gcDuration)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				tokenPool.GC()
				timer.Reset(gcDuration)
			}
		}
	}()
}
