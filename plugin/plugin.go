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

// Package plugin implements a general plugin factory system which provides plugin registration and loading.
// It is mainly used when certain plugins must be loaded by configuration.
// This system is not supposed to register plugins that do not rely on configuration like codec. Instead, plugins
// that do not rely on configuration should be registered by calling methods in certain packages.
package plugin

import (
	"errors"
	"fmt"
)

var plugins = make(map[string]map[string]Factory) // plugin type => { plugin name => plugin factory }

// Factory is the interface for plugin factory abstraction.
// Custom Plugins need to implement this interface to be registered as a plugin with certain type.
type Factory interface {
	// Type returns type of the plugin, i.e. selector, log, config, tracing.
	Type() string
	// Setup loads plugin by configuration.
	// The data structure of the configuration of the plugin needs to be defined in advance。
	Setup(name string, dec Decoder) error
}

// Decoder is the interface used to decode plugin configuration.
type Decoder interface {
	Decode(cfg interface{}) error // the input param is the custom configuration of the plugin
}

// Register registers a plugin factory.
// Name of the plugin should be specified.
// It is supported to register instances which are the same implementation of plugin Factory
// but use different configuration.
func Register(name string, f Factory) {
	factories, ok := plugins[f.Type()]
	if !ok {
		factories = make(map[string]Factory)
		plugins[f.Type()] = factories
	}
	factories[name] = f
}

// RegisterWithValidation registers a plugin factory with validation.
func RegisterWithValidation(name string, f Factory) error {
	if name == "" {
		return errors.New("plugin name cannot be empty")
	}
	if f == nil {
		return errors.New("plugin factory cannot be nil")
	}
	
	pluginType := f.Type()
	if pluginType == "" {
		return errors.New("plugin type cannot be empty")
	}
	
	factories, ok := plugins[pluginType]
	if !ok {
		factories = make(map[string]Factory)
		plugins[pluginType] = factories
	}
	
	if _, exists := factories[name]; exists {
		return fmt.Errorf("plugin %q of type %q is already registered", name, pluginType)
	}
	
	factories[name] = f
	return nil
}

// Get returns a plugin Factory by its type and name.
func Get(typ string, name string) Factory {
	return plugins[typ][name]
}

// GetWithError returns a plugin Factory by its type and name.
func GetWithError(typ string, name string) (Factory, error) {
	if typ == "" {
		return nil, errors.New("plugin type cannot be empty")
	}
	if name == "" {
		return nil, errors.New("plugin name cannot be empty")
	}
	
	factories, ok := plugins[typ]
	if !ok {
		return nil, fmt.Errorf("no plugins registered for type %q", typ)
	}
	
	factory, ok := factories[name]
	if !ok {
		return nil, fmt.Errorf("plugin %q of type %q not found", name, typ)
	}
	
	return factory, nil
}
