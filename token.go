package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"time"
)

type Identity struct {
	username string
	password string
	created  time.Time
}

type TokenPool struct {
	tokens map[string]Identity
}

var expires = 86400

var tokenPool = &TokenPool{
	tokens: make(map[string]Identity),
}

func (t *TokenPool) Add(username string, password string) string {
	identity := Identity{
		username: username,
		password: password,
		created:  time.Now(),
	}

	coder := sha1.New()
	io.WriteString(coder, username+":"+password+":"+fmt.Sprintf("%d", identity.created.Unix()))
	token := fmt.Sprintf("%x", coder.Sum(nil))

	t.tokens[token] = identity
	return token
}

func (t *TokenPool) Get(token string) *Identity {
	id, ok := t.tokens[token]
	if !ok {
		return nil
	}
	return &id
}

func (t *TokenPool) GC() {
	total := 0
	deleted := 0
	for k, v := range t.tokens {
		total++
		if v.created.Unix()+43200 < time.Now().Unix() {
			delete(t.tokens, k)
			deleted++
		}
	}
	log.Println("Token Poll GC: deleted", deleted, "of", total)
}

func tokenPoolGC() {
	go func() {
		duration := 30 * time.Second
		timer := time.NewTimer(duration)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				tokenPool.GC()
				timer.Reset(duration)
			}
		}
	}()
}
