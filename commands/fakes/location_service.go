// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	sync "sync"

	commands "github.com/pivotal-cf/om/commands"
)

type Location struct {
	ContainerStub        func(string) (commands.Container, error)
	containerMutex       sync.RWMutex
	containerArgsForCall []struct {
		arg1 string
	}
	containerReturns struct {
		result1 commands.Container
		result2 error
	}
	containerReturnsOnCall map[int]struct {
		result1 commands.Container
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Location) Container(arg1 string) (commands.Container, error) {
	fake.containerMutex.Lock()
	ret, specificReturn := fake.containerReturnsOnCall[len(fake.containerArgsForCall)]
	fake.containerArgsForCall = append(fake.containerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Container", []interface{}{arg1})
	fake.containerMutex.Unlock()
	if fake.ContainerStub != nil {
		return fake.ContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.containerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Location) ContainerCallCount() int {
	fake.containerMutex.RLock()
	defer fake.containerMutex.RUnlock()
	return len(fake.containerArgsForCall)
}

func (fake *Location) ContainerCalls(stub func(string) (commands.Container, error)) {
	fake.containerMutex.Lock()
	defer fake.containerMutex.Unlock()
	fake.ContainerStub = stub
}

func (fake *Location) ContainerArgsForCall(i int) string {
	fake.containerMutex.RLock()
	defer fake.containerMutex.RUnlock()
	argsForCall := fake.containerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Location) ContainerReturns(result1 commands.Container, result2 error) {
	fake.containerMutex.Lock()
	defer fake.containerMutex.Unlock()
	fake.ContainerStub = nil
	fake.containerReturns = struct {
		result1 commands.Container
		result2 error
	}{result1, result2}
}

func (fake *Location) ContainerReturnsOnCall(i int, result1 commands.Container, result2 error) {
	fake.containerMutex.Lock()
	defer fake.containerMutex.Unlock()
	fake.ContainerStub = nil
	if fake.containerReturnsOnCall == nil {
		fake.containerReturnsOnCall = make(map[int]struct {
			result1 commands.Container
			result2 error
		})
	}
	fake.containerReturnsOnCall[i] = struct {
		result1 commands.Container
		result2 error
	}{result1, result2}
}

func (fake *Location) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.containerMutex.RLock()
	defer fake.containerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Location) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ commands.Location = new(Location)
