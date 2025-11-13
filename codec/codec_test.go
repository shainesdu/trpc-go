//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 Tencent.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the  Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

package codec_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"trpc.group/trpc-go/trpc-go/codec"
)

// go test -v -coverprofile=cover.out
// go tool cover -func=cover.out

// Fake is a fake codec for test
type Fake struct {
}

func (c *Fake) Encode(message codec.Msg, inbody []byte) (outbuf []byte, err error) {
	return nil, nil
}

func (c *Fake) Decode(message codec.Msg, inbuf []byte) (outbody []byte, err error) {
	return nil, nil
}

// TestCodec is unit test for the register logic of codec.
func TestCodec(t *testing.T) {
	f := &Fake{}

	err := codec.Register("fake", f, f)
	assert.NoError(t, err)

	serverCodec, err := codec.GetServer("NoExists")
	assert.Error(t, err)
	assert.Nil(t, serverCodec)

	clientCodec, err := codec.GetClient("NoExists")
	assert.Error(t, err)
	assert.Nil(t, clientCodec)

	serverCodec, err = codec.GetServer("fake")
	assert.NoError(t, err)
	assert.Equal(t, f, serverCodec)

	clientCodec, err = codec.GetClient("fake")
	assert.NoError(t, err)
	assert.Equal(t, f, clientCodec)
}

func TestRegisterWithValidation(t *testing.T) {
	f := &Fake{}

	t.Run("successful registration", func(t *testing.T) {
		err := codec.RegisterWithValidation("test-codec-1", f, f)
		assert.NoError(t, err)
	})

	t.Run("duplicate registration error", func(t *testing.T) {
		err := codec.RegisterWithValidation("test-codec-2", f, f)
		assert.NoError(t, err)
		
		err = codec.RegisterWithValidation("test-codec-2", f, f)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")
	})

	t.Run("empty codec name error", func(t *testing.T) {
		err := codec.RegisterWithValidation("", f, f)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestGetServerWithError(t *testing.T) {
	f := &Fake{}

	t.Run("get non-existent codec", func(t *testing.T) {
		c, err := codec.GetServerWithError("non-existent-codec")
		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("get existing codec", func(t *testing.T) {
		err := codec.Register("test-server-codec", f, f)
		assert.NoError(t, err)
		
		c, err := codec.GetServerWithError("test-server-codec")
		assert.NoError(t, err)
		assert.Equal(t, f, c)
	})
}

func TestGetClientWithError(t *testing.T) {
	f := &Fake{}

	t.Run("get non-existent codec", func(t *testing.T) {
		c, err := codec.GetClientWithError("non-existent-client-codec")
		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("get existing codec", func(t *testing.T) {
		err := codec.Register("test-client-codec", f, f)
		assert.NoError(t, err)
		
		c, err := codec.GetClientWithError("test-client-codec")
		assert.NoError(t, err)
		assert.Equal(t, f, c)
	})
}

// GOMAXPROCS=1 go test -bench=WithNewMessage -benchmem -benchtime=10s
// -memprofile mem.out -cpuprofile cpu.out codec_test.go

// BenchmarkWithNewMessage is the benchmark test of codec
func BenchmarkWithNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		codec.WithNewMessage(context.Background())
	}
}
