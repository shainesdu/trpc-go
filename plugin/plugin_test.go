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

package plugin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"trpc.group/trpc-go/trpc-go/plugin"
)

type mockPlugin struct{}

func (p *mockPlugin) Type() string {
	return pluginType
}

func (p *mockPlugin) Setup(name string, decoder plugin.Decoder) error {
	return nil
}

func TestGet(t *testing.T) {
	plugin.Register(pluginName, &mockPlugin{})
	// test duplicate registration
	plugin.Register(pluginName, &mockPlugin{})
	p := plugin.Get(pluginType, pluginName)
	assert.NotNil(t, p)

	pNo := plugin.Get("notexist", pluginName)
	assert.Nil(t, pNo)
}

func TestRegisterWithValidation(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		err := plugin.RegisterWithValidation("test-plugin-1", &mockPlugin{})
		assert.NoError(t, err)
	})

	t.Run("duplicate registration error", func(t *testing.T) {
		err := plugin.RegisterWithValidation("test-plugin-2", &mockPlugin{})
		assert.NoError(t, err)
		
		err = plugin.RegisterWithValidation("test-plugin-2", &mockPlugin{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")
	})

	t.Run("empty plugin name error", func(t *testing.T) {
		err := plugin.RegisterWithValidation("", &mockPlugin{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("nil factory error", func(t *testing.T) {
		err := plugin.RegisterWithValidation("test-plugin-nil", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

type mockPluginEmptyType struct{}

func (p *mockPluginEmptyType) Type() string {
	return ""
}

func (p *mockPluginEmptyType) Setup(name string, decoder plugin.Decoder) error {
	return nil
}

func TestRegisterWithValidationEmptyType(t *testing.T) {
	t.Run("empty plugin type error", func(t *testing.T) {
		err := plugin.RegisterWithValidation("test-plugin-empty-type", &mockPluginEmptyType{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type cannot be empty")
	})
}

func TestGetWithError(t *testing.T) {
	t.Run("get non-existent plugin type", func(t *testing.T) {
		p, err := plugin.GetWithError("non-existent-type", "some-name")
		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Contains(t, err.Error(), "no plugins registered")
	})

	t.Run("get non-existent plugin name", func(t *testing.T) {
		plugin.Register("existing-plugin", &mockPlugin{})
		
		p, err := plugin.GetWithError(pluginType, "non-existent-name")
		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("get existing plugin", func(t *testing.T) {
		plugin.Register("test-get-plugin", &mockPlugin{})
		
		p, err := plugin.GetWithError(pluginType, "test-get-plugin")
		assert.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("empty plugin type error", func(t *testing.T) {
		p, err := plugin.GetWithError("", "some-name")
		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Contains(t, err.Error(), "type cannot be empty")
	})

	t.Run("empty plugin name error", func(t *testing.T) {
		p, err := plugin.GetWithError("some-type", "")
		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})
}
