// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configtypes "github.com/vmware-tanzu/tanzu-plugin-runtime/config/types"
)

func TestGetAllEnvs(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tests := []struct {
		name   string
		in     *configtypes.ClientConfig
		out    map[string]string
		errStr string
	}{
		{
			name: "success k8s",
			in: &configtypes.ClientConfig{
				ClientOptions: &configtypes.ClientOptions{
					Env: map[string]string{
						"test": "test",
					},
				},
			},
			out: map[string]string{
				"test": "test",
			},
		},
	}
	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			err := StoreClientConfig(spec.in)
			assert.NoError(t, err)
			c, err := GetAllEnvs()
			assert.NoError(t, err)
			assert.Equal(t, spec.out, c)
			assert.NoError(t, err)
		})
	}
}

func TestGetEnv(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()

	tests := []struct {
		name           string
		in             *configtypes.ClientConfig
		out            string
		errStr         string
		errStrForInput string
	}{
		{
			name: "success k8s",
			in: &configtypes.ClientConfig{
				ClientOptions: &configtypes.ClientOptions{
					Env: map[string]string{
						"test": "test",
					},
				},
			},
			out: "test",
		},
		{
			name: "get options with empty key",
			in: &configtypes.ClientConfig{
				ClientOptions: &configtypes.ClientOptions{
					Env: map[string]string{
						"test": "test",
					},
				},
			},
			out:    "",
			errStr: "key cannot be empty",
		},
		{
			name: "store options with empty key",
			in: &configtypes.ClientConfig{
				ClientOptions: &configtypes.ClientOptions{
					Env: map[string]string{
						"": "test",
					},
				},
			},
			out:            "",
			errStr:         "key cannot be empty",
			errStrForInput: "key cannot be empty",
		},
		{
			name: "store options with empty val",
			in: &configtypes.ClientConfig{
				ClientOptions: &configtypes.ClientOptions{
					Env: map[string]string{
						"test-empty-val": "",
					},
				},
			},
			out:            "test-empty-val",
			errStr:         "not found",
			errStrForInput: "value cannot be empty",
		},
	}
	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			err := StoreClientConfig(spec.in)
			if spec.errStrForInput != "" {
				assert.Equal(t, spec.errStrForInput, err.Error())
			} else {
				assert.NoError(t, err)
			}

			c, err := GetEnv(spec.out)
			if spec.errStr != "" {
				assert.Equal(t, spec.errStr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, spec.out, c)
			}
		})
	}
}

func TestSetEnv(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tests := []struct {
		name    string
		key     string
		val     string
		persist bool
	}{
		{
			name: "should add new env to empty envs",
			key:  "test",
			val:  "test-test",
		},
		{
			name: "should add new env to existing envs",
			key:  "test2",
			val:  "test2",
		},
		{
			name: "should update existing env",
			key:  "test",
			val:  "updated-test",
		},
		{
			name: "should not update same env",
			key:  "test2",
			val:  "test2",
		},
	}

	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			err := SetEnv(spec.key, spec.val)
			assert.NoError(t, err)
			val, err := GetEnv(spec.key)
			assert.Equal(t, spec.val, val)
			assert.NoError(t, err)
		})
	}
}
func TestDeleteEnv(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
		cfg := &configtypes.ClientConfig{
			ClientOptions: &configtypes.ClientOptions{
				Env: map[string]string{
					"test":  "test",
					"test2": "test2",
					"test4": "test2",
				},
			},
		}
		err := StoreClientConfig(cfg)
		assert.NoError(t, err)
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "success delete test",
			in:   "test",
			out:  true,
		},
		{
			name: "success delete test2",
			in:   "test2",
			out:  true,
		},

		{
			name: "success delete test3",
			in:   "test3",
			out:  true,
		},
	}

	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			err := DeleteEnv(spec.in)
			assert.NoError(t, err)
			c, err := GetEnv(spec.in)
			assert.Equal(t, "not found", err.Error())
			assert.Equal(t, spec.out, c == "")
		})
	}
}
