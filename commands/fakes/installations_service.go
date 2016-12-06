// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/pivotal-cf/om/api"
)

type InstallationsService struct {
	TriggerStub        func() (api.InstallationsServiceOutput, error)
	triggerMutex       sync.RWMutex
	triggerArgsForCall []struct{}
	triggerReturns     struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}
	StatusStub        func(id int) (api.InstallationsServiceOutput, error)
	statusMutex       sync.RWMutex
	statusArgsForCall []struct {
		id int
	}
	statusReturns struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}
	LogsStub        func(id int) (api.InstallationsServiceOutput, error)
	logsMutex       sync.RWMutex
	logsArgsForCall []struct {
		id int
	}
	logsReturns struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}
	RunningInstallationStub        func() (api.InstallationsServiceOutput, error)
	runningInstallationMutex       sync.RWMutex
	runningInstallationArgsForCall []struct{}
	runningInstallationReturns     struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *InstallationsService) Trigger() (api.InstallationsServiceOutput, error) {
	fake.triggerMutex.Lock()
	fake.triggerArgsForCall = append(fake.triggerArgsForCall, struct{}{})
	fake.recordInvocation("Trigger", []interface{}{})
	fake.triggerMutex.Unlock()
	if fake.TriggerStub != nil {
		return fake.TriggerStub()
	} else {
		return fake.triggerReturns.result1, fake.triggerReturns.result2
	}
}

func (fake *InstallationsService) TriggerCallCount() int {
	fake.triggerMutex.RLock()
	defer fake.triggerMutex.RUnlock()
	return len(fake.triggerArgsForCall)
}

func (fake *InstallationsService) TriggerReturns(result1 api.InstallationsServiceOutput, result2 error) {
	fake.TriggerStub = nil
	fake.triggerReturns = struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}{result1, result2}
}

func (fake *InstallationsService) Status(id int) (api.InstallationsServiceOutput, error) {
	fake.statusMutex.Lock()
	fake.statusArgsForCall = append(fake.statusArgsForCall, struct {
		id int
	}{id})
	fake.recordInvocation("Status", []interface{}{id})
	fake.statusMutex.Unlock()
	if fake.StatusStub != nil {
		return fake.StatusStub(id)
	} else {
		return fake.statusReturns.result1, fake.statusReturns.result2
	}
}

func (fake *InstallationsService) StatusCallCount() int {
	fake.statusMutex.RLock()
	defer fake.statusMutex.RUnlock()
	return len(fake.statusArgsForCall)
}

func (fake *InstallationsService) StatusArgsForCall(i int) int {
	fake.statusMutex.RLock()
	defer fake.statusMutex.RUnlock()
	return fake.statusArgsForCall[i].id
}

func (fake *InstallationsService) StatusReturns(result1 api.InstallationsServiceOutput, result2 error) {
	fake.StatusStub = nil
	fake.statusReturns = struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}{result1, result2}
}

func (fake *InstallationsService) Logs(id int) (api.InstallationsServiceOutput, error) {
	fake.logsMutex.Lock()
	fake.logsArgsForCall = append(fake.logsArgsForCall, struct {
		id int
	}{id})
	fake.recordInvocation("Logs", []interface{}{id})
	fake.logsMutex.Unlock()
	if fake.LogsStub != nil {
		return fake.LogsStub(id)
	} else {
		return fake.logsReturns.result1, fake.logsReturns.result2
	}
}

func (fake *InstallationsService) LogsCallCount() int {
	fake.logsMutex.RLock()
	defer fake.logsMutex.RUnlock()
	return len(fake.logsArgsForCall)
}

func (fake *InstallationsService) LogsArgsForCall(i int) int {
	fake.logsMutex.RLock()
	defer fake.logsMutex.RUnlock()
	return fake.logsArgsForCall[i].id
}

func (fake *InstallationsService) LogsReturns(result1 api.InstallationsServiceOutput, result2 error) {
	fake.LogsStub = nil
	fake.logsReturns = struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}{result1, result2}
}

func (fake *InstallationsService) RunningInstallation() (api.InstallationsServiceOutput, error) {
	fake.runningInstallationMutex.Lock()
	fake.runningInstallationArgsForCall = append(fake.runningInstallationArgsForCall, struct{}{})
	fake.recordInvocation("RunningInstallation", []interface{}{})
	fake.runningInstallationMutex.Unlock()
	if fake.RunningInstallationStub != nil {
		return fake.RunningInstallationStub()
	} else {
		return fake.runningInstallationReturns.result1, fake.runningInstallationReturns.result2
	}
}

func (fake *InstallationsService) RunningInstallationCallCount() int {
	fake.runningInstallationMutex.RLock()
	defer fake.runningInstallationMutex.RUnlock()
	return len(fake.runningInstallationArgsForCall)
}

func (fake *InstallationsService) RunningInstallationReturns(result1 api.InstallationsServiceOutput, result2 error) {
	fake.RunningInstallationStub = nil
	fake.runningInstallationReturns = struct {
		result1 api.InstallationsServiceOutput
		result2 error
	}{result1, result2}
}

func (fake *InstallationsService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.triggerMutex.RLock()
	defer fake.triggerMutex.RUnlock()
	fake.statusMutex.RLock()
	defer fake.statusMutex.RUnlock()
	fake.logsMutex.RLock()
	defer fake.logsMutex.RUnlock()
	fake.runningInstallationMutex.RLock()
	defer fake.runningInstallationMutex.RUnlock()
	return fake.invocations
}

func (fake *InstallationsService) recordInvocation(key string, args []interface{}) {
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
