// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	configtypes "github.com/vmware-tanzu/tanzu-plugin-runtime/config/types"
)

func TestSetGetDeleteContext(t *testing.T) {
	// Setup config data
	_, cleanUp := setupTestConfig(t, &CfgTestData{})

	defer func() {
		cleanUp()
	}()

	ctx1 := &configtypes.Context{
		Name:   "test1",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			Path:                "test-path",
			Context:             "test-context",
			IsManagementCluster: true,
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "updated-test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	ctx2 := &configtypes.Context{
		Name:   "test2",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			Path:                "test-path",
			Context:             "test-context",
			IsManagementCluster: true,
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "updated-test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	ctx, err := GetContext("test1")
	assert.Equal(t, "context test1 not found", err.Error())
	assert.Nil(t, ctx)

	err = SetContext(ctx1, true)
	assert.NoError(t, err)

	ctx, err = GetContext("test1")
	assert.Nil(t, err)
	assert.Equal(t, ctx1, ctx)

	ctx, err = GetCurrentContext(configtypes.TargetK8s)
	assert.Nil(t, err)
	assert.Equal(t, ctx1, ctx)

	err = SetContext(ctx2, false)
	assert.NoError(t, err)

	ctx, err = GetContext("test2")
	assert.Nil(t, err)
	assert.Equal(t, ctx2, ctx)

	ctx, err = GetCurrentContext(configtypes.TargetK8s)
	assert.Nil(t, err)
	assert.Equal(t, ctx1, ctx)

	err = DeleteContext("test")
	assert.Equal(t, "context test not found", err.Error())

	err = DeleteContext("test1")
	assert.Nil(t, err)

	ctx, err = GetContext("test1")
	assert.Nil(t, ctx)
	assert.Equal(t, "context test1 not found", err.Error())
}

func TestSetContextWithOldVersion(t *testing.T) {
	tanzuConfigBytes := `
currentContext:
    kubernetes: test-mc
contexts:
    - name: test-mc
      ctx-field: new-ctx-field
      optional: true
      target: kubernetes
      clusterOpts:
        isManagementCluster: true
        endpoint: old-test-endpoint
        annotation: one
        required: true
        annotationStruct:
            one: one
      discoverySources:
        - gcp:
            name: test
            bucket: test-ctx-bucket
            manifestPath: test-ctx-manifest-path
            annotation: one
            required: true
        - gcp:
            name: test2
            bucket: test2-bucket
            manifestPath: test2-manifest-path
            annotation: one
            required: true
`

	// Setup config data
	_, cleanUp := setupTestConfig(t, &CfgTestData{cfgNextGen: tanzuConfigBytes})

	defer func() {
		cleanUp()
	}()

	ctx := &configtypes.Context{
		Name:   "test-mc",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			Path:                "test-path",
			Context:             "test-context",
			IsManagementCluster: true,
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "updated-test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	err := SetContext(ctx, false)
	assert.NoError(t, err)

	c, err := GetContext(ctx.Name)
	assert.NoError(t, err)
	assert.Equal(t, c.Name, ctx.Name)
	assert.Equal(t, c.ClusterOpts.Endpoint, "old-test-endpoint")
	assert.Equal(t, c.ClusterOpts.Path, ctx.ClusterOpts.Path)
	assert.Equal(t, c.ClusterOpts.Context, ctx.ClusterOpts.Context)
}

func TestSetContextWithDiscoverySourceWithNewFields(t *testing.T) {
	// Setup config data
	_, cleanUp := setupTestConfig(t, &CfgTestData{})

	defer func() {
		cleanUp()
	}()

	tests := []struct {
		name    string
		src     *configtypes.ClientConfig
		ctx     *configtypes.Context
		current bool
		errStr  string
	}{
		{
			name: "should add new context with new discovery sources to empty client config",
			src:  &configtypes.ClientConfig{},
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
				DiscoverySources: []configtypes.PluginDiscovery{
					{
						GCP: &configtypes.GCPDiscovery{
							Name:         "test",
							Bucket:       "updated-test-bucket",
							ManifestPath: "test-manifest-path",
						},
					},
				},
			},
			current: true,
		},
		{
			name: "should update existing context",
			src: &configtypes.ClientConfig{
				KnownContexts: []*configtypes.Context{
					{
						Name:   "test-mc",
						Target: configtypes.TargetK8s,
						ClusterOpts: &configtypes.ClusterServer{
							Endpoint:            "test-endpoint",
							Path:                "test-path",
							Context:             "test-context",
							IsManagementCluster: true,
						},
						DiscoverySources: []configtypes.PluginDiscovery{
							{
								GCP: &configtypes.GCPDiscovery{
									Name:         "test",
									Bucket:       "test-bucket",
									ManifestPath: "test-manifest-path",
								},
							},
						},
					},
				},
				KnownServers: []*configtypes.Server{
					{
						Name: "test-mc",
						Type: configtypes.ManagementClusterServerType,
						ManagementClusterOpts: &configtypes.ManagementClusterServer{
							Endpoint: "test-endpoint",
							Path:     "test-path",
							Context:  "test-context",
						},
						DiscoverySources: []configtypes.PluginDiscovery{
							{
								GCP: &configtypes.GCPDiscovery{
									Name:         "test",
									Bucket:       "test-bucket",
									ManifestPath: "test-manifest-path",
								},
							},
						},
					},
				},
				CurrentServer: "test-mc",
				CurrentContext: map[configtypes.Target]string{
					configtypes.TargetK8s: "test-mc",
				},
			},
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "updated-test-endpoint",
					Path:                "updated-test-path",
					Context:             "updated-test-context",
					IsManagementCluster: true,
				},
				DiscoverySources: []configtypes.PluginDiscovery{
					{
						GCP: &configtypes.GCPDiscovery{
							Name:         "test",
							Bucket:       "updated-test-bucket",
							ManifestPath: "updated-test-manifest-path",
						},
					},
				},
			},
			current: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SetContext(tc.ctx, tc.current)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}

			ok, err := ContextExists(tc.ctx.Name)
			assert.True(t, ok)
			assert.NoError(t, err)
		})
	}
}

func TestSetContextWithDiscoverySource(t *testing.T) {
	// Setup config data
	_, cleanUp := setupTestConfig(t, &CfgTestData{})

	defer func() {
		cleanUp()
	}()

	tests := []struct {
		name    string
		src     *configtypes.ClientConfig
		ctx     *configtypes.Context
		current bool
		errStr  string
	}{
		{
			name: "should add new context with new discovery sources to empty client config",
			src:  &configtypes.ClientConfig{},
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
				DiscoverySources: []configtypes.PluginDiscovery{
					{
						GCP: &configtypes.GCPDiscovery{
							Name:         "test",
							Bucket:       "updated-test-bucket",
							ManifestPath: "test-manifest-path",
						},
					},
				},
			},
			current: true,
		},
		{
			name: "should update existing context",
			src: &configtypes.ClientConfig{
				KnownContexts: []*configtypes.Context{
					{
						Name:   "test-mc",
						Target: configtypes.TargetK8s,
						ClusterOpts: &configtypes.ClusterServer{
							Endpoint:            "test-endpoint",
							Path:                "test-path",
							Context:             "test-context",
							IsManagementCluster: true,
						},
						DiscoverySources: []configtypes.PluginDiscovery{
							{
								GCP: &configtypes.GCPDiscovery{
									Name:         "test",
									Bucket:       "test-bucket",
									ManifestPath: "test-manifest-path",
								},
							},
						},
					},
				},
				KnownServers: []*configtypes.Server{
					{
						Name: "test-mc",
						Type: configtypes.ManagementClusterServerType,
						ManagementClusterOpts: &configtypes.ManagementClusterServer{
							Endpoint: "test-endpoint",
							Path:     "test-path",
							Context:  "test-context",
						},
						DiscoverySources: []configtypes.PluginDiscovery{
							{
								GCP: &configtypes.GCPDiscovery{
									Name:         "test",
									Bucket:       "test-bucket",
									ManifestPath: "test-manifest-path",
								},
							},
						},
					},
				},
				CurrentServer: "test-mc",
				CurrentContext: map[configtypes.Target]string{
					configtypes.TargetK8s: "test-mc",
				},
			},
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "updated-test-endpoint",
					Path:                "updated-test-path",
					Context:             "updated-test-context",
					IsManagementCluster: true,
				},
				DiscoverySources: []configtypes.PluginDiscovery{
					{
						GCP: &configtypes.GCPDiscovery{
							Name:         "test",
							Bucket:       "updated-test-bucket",
							ManifestPath: "updated-test-manifest-path",
						},
					},
				},
			},
			current: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SetContext(tc.ctx, tc.current)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}

			ok, err := ContextExists(tc.ctx.Name)
			assert.True(t, ok)
			assert.NoError(t, err)
		})
	}
}

func setupForGetContext(t *testing.T) {
	// setup
	cfg := &configtypes.ClientConfig{
		KnownContexts: []*configtypes.Context{
			{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
			},
			{
				Name:   "test-mc-2",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint-2",
					Path:                "test-path-2",
					Context:             "test-context-2",
					IsManagementCluster: true,
				},
			},
			{
				Name:   "test-tmc",
				Target: configtypes.TargetTMC,
				GlobalOpts: &configtypes.GlobalServer{
					Endpoint: "test-endpoint",
				},
			},
		},
		CurrentContext: map[configtypes.Target]string{
			configtypes.TargetK8s: "test-mc-2",
			configtypes.TargetTMC: "test-tmc",
		},
	}
	func() {
		LocalDirName = TestLocalDirName
		err := StoreClientConfig(cfg)
		assert.NoError(t, err)
	}()
}

func TestGetContext(t *testing.T) {
	setupForGetContext(t)

	defer func() {
		cleanupDir(LocalDirName)
	}()

	tcs := []struct {
		name    string
		ctxName string
		errStr  string
	}{
		{
			name:    "success k8s",
			ctxName: "test-mc",
		},
		{
			name:    "success tmc",
			ctxName: "test-tmc",
		},
		{
			name:    "failure",
			ctxName: "test",
			errStr:  "context test not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c, err := GetContext(tc.ctxName)
			if tc.errStr == "" {
				assert.Equal(t, tc.ctxName, c.Name)
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
		})
	}
}

func TestContextExists(t *testing.T) {
	setupForGetContext(t)

	defer func() {
		cleanupDir(LocalDirName)
	}()

	tcs := []struct {
		name    string
		ctxName string
		ok      bool
	}{
		{
			name:    "success k8s",
			ctxName: "test-mc",
			ok:      true,
		},
		{
			name:    "success tmc",
			ctxName: "test-tmc",
			ok:      true,
		},
		{
			name:    "failure",
			ctxName: "test",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := ContextExists(tc.ctxName)
			assert.Equal(t, tc.ok, ok)
			assert.NoError(t, err)
		})
	}
}

func TestSetContext(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
		// setup data
		node := &yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{
					Kind:    yaml.MappingNode,
					Content: []*yaml.Node{},
				},
			},
		}
		err := persistConfig(node)
		assert.NoError(t, err)
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tcs := []struct {
		name    string
		ctx     *configtypes.Context
		current bool
		errStr  string
	}{
		{
			name: "should add new context and set current on empty config",
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
			},
			current: true,
		},

		{
			name: "should add new context but not current",
			ctx: &configtypes.Context{
				Name:   "test-mc2",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
			},
		},
		{
			name: "success tmc current",
			ctx: &configtypes.Context{
				Name:   "test-tmc1",
				Target: configtypes.TargetTMC,
				GlobalOpts: &configtypes.GlobalServer{
					Endpoint: "test-endpoint",
				},
			},
			current: true,
		},
		{
			name: "success tmc not_current",
			ctx: &configtypes.Context{
				Name:   "test-tmc2",
				Target: configtypes.TargetTMC,
				GlobalOpts: &configtypes.GlobalServer{
					Endpoint: "test-endpoint",
				},
			},
		},
		{
			name: "success update test-mc",
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "good-test-endpoint",
					Path:                "updated-test-path",
					Context:             "updated-test-context",
					IsManagementCluster: true,
				},
			},
		},
		{
			name: "success update tmc",
			ctx: &configtypes.Context{
				Name:   "test-tmc",
				Target: configtypes.TargetTMC,
				GlobalOpts: &configtypes.GlobalServer{
					Endpoint: "updated-test-endpoint",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// perform test
			err := SetContext(tc.ctx, tc.current)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
			ctx, err := GetContext(tc.ctx.Name)
			assert.NoError(t, err)
			assert.Equal(t, tc.ctx.Name, ctx.Name)
			s, err := GetServer(tc.ctx.Name)
			assert.NoError(t, err)
			assert.Equal(t, tc.ctx.Name, s.Name)
		})
	}
}

func TestRemoveContext(t *testing.T) {
	// setup
	setupForGetContext(t)
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tcs := []struct {
		name    string
		ctxName string
		target  configtypes.Target
		errStr  string
	}{
		{
			name:    "success k8s",
			ctxName: "test-mc",
			target:  configtypes.TargetK8s,
		},
		{
			name:    "success tmc",
			ctxName: "test-tmc",
			target:  configtypes.TargetTMC,
		},
		{
			name:    "failure",
			ctxName: "test",
			errStr:  "context test not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errStr == "" {
				ok, err := ContextExists(tc.ctxName)
				require.True(t, ok)
				require.NoError(t, err)
			}
			err := RemoveContext(tc.ctxName)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
			ok, err := ContextExists(tc.ctxName)
			assert.False(t, ok)
			assert.NoError(t, err)
			ok, err = ServerExists(tc.ctxName)
			assert.Nil(t, err)
			assert.False(t, ok)
		})
	}
}

func TestSetCurrentContext(t *testing.T) {
	// setup
	setupForGetContext(t)
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tcs := []struct {
		name       string
		target     configtypes.Target
		ctxName    string
		currServer string
		errStr     string
	}{
		{
			name:    "success tmc",
			ctxName: "test-tmc",
			target:  configtypes.TargetTMC,
		},
		{
			name:       "success k8s",
			ctxName:    "test-mc",
			target:     configtypes.TargetK8s,
			currServer: "test-mc",
		},
		{
			name:    "success tmc after setting k8s",
			ctxName: "test-tmc",
			target:  configtypes.TargetTMC,
		},
		{
			name:    "failure",
			ctxName: "test",
			errStr:  "context test not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			prevSrv, _ := GetCurrentServer()

			err := SetCurrentContext(tc.ctxName)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
			currCtx, err := GetCurrentContext(tc.target)
			if tc.errStr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.ctxName, currCtx.Name)
			} else {
				assert.Error(t, err)
			}
			currSrv, err := GetCurrentServer()
			assert.NoError(t, err)
			if tc.errStr == "" {
				if tc.currServer == "" {
					assert.Equal(t, prevSrv.Name, currSrv.Name)
				} else {
					assert.Equal(t, tc.currServer, currSrv.Name)
				}
			}
		})
	}

	currentContextMap, err := GetAllCurrentContextsMap()
	assert.NoError(t, err)
	assert.Equal(t, "test-mc", currentContextMap[configtypes.TargetK8s].Name)
	assert.Equal(t, "test-tmc", currentContextMap[configtypes.TargetTMC].Name)

	currentContextsList, err := GetAllCurrentContextsList()
	assert.NoError(t, err)
	assert.Contains(t, currentContextsList, "test-mc")
	assert.Contains(t, currentContextsList, "test-tmc")
}

func TestRemoveCurrentContext(t *testing.T) {
	// setup
	setupForGetContext(t)
	defer func() {
		cleanupDir(LocalDirName)
	}()

	err := RemoveCurrentContext(configtypes.TargetK8s)
	assert.NoError(t, err)

	currCtx, err := GetCurrentContext(configtypes.TargetK8s)
	assert.Equal(t, "no current context set for target \"kubernetes\"", err.Error())
	assert.Nil(t, currCtx)

	currSrv, err := GetCurrentServer()
	assert.Equal(t, "current server \"\" not found in tanzu config", err.Error())
	assert.Nil(t, currSrv)

	currCtx, err = GetCurrentContext(configtypes.TargetTMC)
	assert.NoError(t, err)
	assert.Equal(t, currCtx.Name, "test-tmc")
}

func TestSetSingleContext(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tcs := []struct {
		name    string
		ctx     *configtypes.Context
		current bool
		errStr  string
	}{
		{
			name: "success k8s current",
			ctx: &configtypes.Context{
				Name:   "test-mc",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := SetContext(tc.ctx, tc.current)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
			ok, err := ContextExists(tc.ctx.Name)
			assert.True(t, ok)
			assert.NoError(t, err)
			ok, err = ServerExists(tc.ctx.Name)
			assert.True(t, ok)
			assert.NoError(t, err)
		})
	}
}

func TestSetContextMultiFile(t *testing.T) {
	configBytes, configNextGenBytes := setupMultiCfgData()
	// Setup config data
	_, cleanUp := setupTestConfig(t, &CfgTestData{cfg: configBytes, cfgNextGen: configNextGenBytes})

	defer func() {
		cleanUp()
	}()

	ctx := &configtypes.Context{
		Name:   "test-mc",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			IsManagementCluster: true,
			Endpoint:            "test-endpoint",
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test2",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	ctx2 := &configtypes.Context{
		Name:   "test-mc2",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			Endpoint: "updated-test-endpoint",
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name: "test",
				},
			},
			{
				GCP: &configtypes.GCPDiscovery{
					Name: "test2",
				},
			},
		},
	}

	expectedCtx2 := &configtypes.Context{
		Name:   "test-mc2",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			IsManagementCluster: true,
			Endpoint:            "updated-test-endpoint",
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test2",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	c, err := GetCurrentContext(configtypes.TargetK8s)
	assert.NoError(t, err)
	assert.Equal(t, ctx, c)

	c, err = GetContext("test-mc")
	assert.NoError(t, err)
	assert.Equal(t, ctx, c)

	err = SetContext(ctx2, true)
	assert.NoError(t, err)

	c, err = GetContext(ctx2.Name)
	assert.NoError(t, err)
	assert.Equal(t, expectedCtx2, c)
}

func TestSetContextMultiFileAndMigrateToNewConfig(t *testing.T) {
	configBytes, configNextGenBytes := setupMultiCfgData()

	// Setup config data
	_, cleanUp := setupTestConfig(t, &CfgTestData{cfg: configBytes, cfgNextGen: configNextGenBytes, cfgMetadata: setupConfigMetadataWithMigrateToNewConfig()})

	defer func() {
		cleanUp()
	}()

	ctx := &configtypes.Context{
		Name:   "test-mc",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			IsManagementCluster: true,
			Endpoint:            "test-endpoint",
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test2",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	ctx2 := &configtypes.Context{
		Name:   "test-mc2",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			Endpoint: "updated-test-endpoint",
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name: "test",
				},
			},
			{
				GCP: &configtypes.GCPDiscovery{
					Name: "test2",
				},
			},
		},
	}

	expectedCtx2 := &configtypes.Context{
		Name:   "test-mc2",
		Target: configtypes.TargetK8s,
		ClusterOpts: &configtypes.ClusterServer{
			IsManagementCluster: true,
			Endpoint:            "updated-test-endpoint",
		},
		DiscoverySources: []configtypes.PluginDiscovery{
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
			{
				GCP: &configtypes.GCPDiscovery{
					Name:         "test2",
					Bucket:       "test-bucket",
					ManifestPath: "test-manifest-path",
				},
			},
		},
	}

	c, err := GetCurrentContext(configtypes.TargetK8s)
	assert.NoError(t, err)
	assert.Equal(t, ctx, c)

	c, err = GetContext("test-mc")
	assert.NoError(t, err)
	assert.Equal(t, ctx, c)

	err = SetContext(ctx2, true)
	assert.NoError(t, err)

	c, err = GetContext(ctx2.Name)
	assert.NoError(t, err)
	assert.Equal(t, expectedCtx2, c)
}

func TestSetContextWithUniquePermissions(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()

	defer func() {
		cleanupDir(LocalDirName)
	}()

	ctx := &configtypes.Context{
		Name:   "test-mc",
		Target: configtypes.TargetTMC,
		GlobalOpts: &configtypes.GlobalServer{
			Endpoint: "test-endpoint",
			Auth: configtypes.GlobalServerAuth{
				IDToken: "",
				Issuer:  "https://console-stg.cloud.vmware.com/csp/gateway/am/api",
				Permissions: []string{
					"external/25834195-19aa-4ffd-8933-f5f20094ab24/service:member",
					"csp:org_owner",
					"external/f52d39b0-c298-4adf-9c6f-0a4a07351cd7/service:admin",
					"csp:org_member",
					"external/f52d39b0-c298-4adf-9c6f-0a4a07351cd7/service:member",
				},
				RefreshToken: "XXX",
				Type:         "api-token",
				UserName:     "tanzu-core",
			},
		},
	}

	ctx2 := &configtypes.Context{
		Name:   "test-mc",
		Target: configtypes.TargetTMC,
		GlobalOpts: &configtypes.GlobalServer{
			Endpoint: "test-endpoint-updated",
			Auth: configtypes.GlobalServerAuth{
				IDToken: "",
				Issuer:  "https://console-stg.cloud.vmware.com/csp/gateway/am/api",
				Permissions: []string{
					"csp:org_member2",
					"external/f52d39b0-c298-4adf-9c6f-0a4a07351cd7/service:member",
				},
				RefreshToken: "XXX",
				Type:         "api-token",
				UserName:     "tanzu-core",
			},
		},
	}

	ctx3 := &configtypes.Context{
		Name:   "test-mc",
		Target: configtypes.TargetTMC,
		GlobalOpts: &configtypes.GlobalServer{
			Endpoint: "test-endpoint-updated3",
			Auth: configtypes.GlobalServerAuth{
				IDToken: "",
				Issuer:  "https://console-stg.cloud.vmware.com/csp/gateway/am/api",
				Permissions: []string{
					"external/25834195-19aa-4ffd-8933-f5f20094ab24/service:member",
					"csp:org_owner3",
				},
				RefreshToken: "XXX",
				Type:         "api-token",
				UserName:     "tanzu-core",
			},
		},
	}

	for i := 1; i <= 100; i++ {
		err := SetContext(ctx, true)
		assert.NoError(t, err)
		err = SetContext(ctx2, true)
		assert.NoError(t, err)
		err = SetContext(ctx3, true)
		assert.NoError(t, err)
		err = SetContext(ctx, true)
		assert.NoError(t, err)
	}

	c, err := GetContext("test-mc")
	assert.NoError(t, err)
	assert.Equal(t, 7, len(c.GlobalOpts.Auth.Permissions))

	s, err := GetServer("test-mc")
	assert.NoError(t, err)
	assert.Equal(t, 7, len(s.GlobalOpts.Auth.Permissions))
}

func TestSetContextWithEmptyName(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tcs := []struct {
		name    string
		ctx     *configtypes.Context
		current bool
		errStr  string
	}{
		{
			name: "success  current",
			ctx: &configtypes.Context{
				Name:   "",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
			},
			errStr: "context name cannot be empty",
		},
		{
			name: "success re empty current",
			ctx: &configtypes.Context{
				Name:   "",
				Target: configtypes.TargetK8s,
				ClusterOpts: &configtypes.ClusterServer{
					Endpoint:            "test-endpoint",
					Path:                "test-path",
					Context:             "test-context",
					IsManagementCluster: true,
				},
			},
			errStr: "context name cannot be empty",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := SetContext(tc.ctx, tc.current)
			if tc.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.errStr)
			}
		})
	}
}
