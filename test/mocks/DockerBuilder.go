// Copyright 2022 Giuseppe De Palma, Matteo Trentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	docker "github.com/funlessdev/fl-cli/pkg/docker"
	mock "github.com/stretchr/testify/mock"
)

// DockerBuilder is an autogenerated mock type for the DockerBuilder type
type DockerBuilder struct {
	mock.Mock
}

// BuildSource provides a mock function with given fields: ctx, srcPath
func (_m *DockerBuilder) BuildSource(ctx context.Context, srcPath string) error {
	ret := _m.Called(ctx, srcPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, srcPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PullBuilderImage provides a mock function with given fields: ctx
func (_m *DockerBuilder) PullBuilderImage(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenameCodeWasm provides a mock function with given fields: name
func (_m *DockerBuilder) RenameCodeWasm(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Setup provides a mock function with given fields: client, language, dest
func (_m *DockerBuilder) Setup(client docker.DockerClient, language string, dest string) error {
	ret := _m.Called(client, language, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(docker.DockerClient, string, string) error); ok {
		r0 = rf(client, language, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDockerBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewDockerBuilder creates a new instance of DockerBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDockerBuilder(t mockConstructorTestingTNewDockerBuilder) *DockerBuilder {
	mock := &DockerBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
