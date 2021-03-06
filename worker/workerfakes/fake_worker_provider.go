// This file was generated by counterfeiter
package workerfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/worker"
)

type FakeWorkerProvider struct {
	RunningWorkersStub        func() ([]worker.Worker, error)
	runningWorkersMutex       sync.RWMutex
	runningWorkersArgsForCall []struct{}
	runningWorkersReturns     struct {
		result1 []worker.Worker
		result2 error
	}
	runningWorkersReturnsOnCall map[int]struct {
		result1 []worker.Worker
		result2 error
	}
	GetWorkerStub        func(string) (worker.Worker, bool, error)
	getWorkerMutex       sync.RWMutex
	getWorkerArgsForCall []struct {
		arg1 string
	}
	getWorkerReturns struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}
	getWorkerReturnsOnCall map[int]struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}
	FindWorkerForResourceCheckContainerStub        func(logger lager.Logger, teamID int, resourceUser dbng.ResourceUser, resourceType string, resourceSource atc.Source, types atc.VersionedResourceTypes) (worker.Worker, bool, error)
	findWorkerForResourceCheckContainerMutex       sync.RWMutex
	findWorkerForResourceCheckContainerArgsForCall []struct {
		logger         lager.Logger
		teamID         int
		resourceUser   dbng.ResourceUser
		resourceType   string
		resourceSource atc.Source
		types          atc.VersionedResourceTypes
	}
	findWorkerForResourceCheckContainerReturns struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}
	findWorkerForResourceCheckContainerReturnsOnCall map[int]struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}
	FindContainerForIdentifierStub        func(worker.Identifier) (db.SavedContainer, bool, error)
	findContainerForIdentifierMutex       sync.RWMutex
	findContainerForIdentifierArgsForCall []struct {
		arg1 worker.Identifier
	}
	findContainerForIdentifierReturns struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}
	findContainerForIdentifierReturnsOnCall map[int]struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}
	GetContainerStub        func(string) (db.SavedContainer, bool, error)
	getContainerMutex       sync.RWMutex
	getContainerArgsForCall []struct {
		arg1 string
	}
	getContainerReturns struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}
	getContainerReturnsOnCall map[int]struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWorkerProvider) RunningWorkers() ([]worker.Worker, error) {
	fake.runningWorkersMutex.Lock()
	ret, specificReturn := fake.runningWorkersReturnsOnCall[len(fake.runningWorkersArgsForCall)]
	fake.runningWorkersArgsForCall = append(fake.runningWorkersArgsForCall, struct{}{})
	fake.recordInvocation("RunningWorkers", []interface{}{})
	fake.runningWorkersMutex.Unlock()
	if fake.RunningWorkersStub != nil {
		return fake.RunningWorkersStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.runningWorkersReturns.result1, fake.runningWorkersReturns.result2
}

func (fake *FakeWorkerProvider) RunningWorkersCallCount() int {
	fake.runningWorkersMutex.RLock()
	defer fake.runningWorkersMutex.RUnlock()
	return len(fake.runningWorkersArgsForCall)
}

func (fake *FakeWorkerProvider) RunningWorkersReturns(result1 []worker.Worker, result2 error) {
	fake.RunningWorkersStub = nil
	fake.runningWorkersReturns = struct {
		result1 []worker.Worker
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkerProvider) RunningWorkersReturnsOnCall(i int, result1 []worker.Worker, result2 error) {
	fake.RunningWorkersStub = nil
	if fake.runningWorkersReturnsOnCall == nil {
		fake.runningWorkersReturnsOnCall = make(map[int]struct {
			result1 []worker.Worker
			result2 error
		})
	}
	fake.runningWorkersReturnsOnCall[i] = struct {
		result1 []worker.Worker
		result2 error
	}{result1, result2}
}

func (fake *FakeWorkerProvider) GetWorker(arg1 string) (worker.Worker, bool, error) {
	fake.getWorkerMutex.Lock()
	ret, specificReturn := fake.getWorkerReturnsOnCall[len(fake.getWorkerArgsForCall)]
	fake.getWorkerArgsForCall = append(fake.getWorkerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetWorker", []interface{}{arg1})
	fake.getWorkerMutex.Unlock()
	if fake.GetWorkerStub != nil {
		return fake.GetWorkerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.getWorkerReturns.result1, fake.getWorkerReturns.result2, fake.getWorkerReturns.result3
}

func (fake *FakeWorkerProvider) GetWorkerCallCount() int {
	fake.getWorkerMutex.RLock()
	defer fake.getWorkerMutex.RUnlock()
	return len(fake.getWorkerArgsForCall)
}

func (fake *FakeWorkerProvider) GetWorkerArgsForCall(i int) string {
	fake.getWorkerMutex.RLock()
	defer fake.getWorkerMutex.RUnlock()
	return fake.getWorkerArgsForCall[i].arg1
}

func (fake *FakeWorkerProvider) GetWorkerReturns(result1 worker.Worker, result2 bool, result3 error) {
	fake.GetWorkerStub = nil
	fake.getWorkerReturns = struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) GetWorkerReturnsOnCall(i int, result1 worker.Worker, result2 bool, result3 error) {
	fake.GetWorkerStub = nil
	if fake.getWorkerReturnsOnCall == nil {
		fake.getWorkerReturnsOnCall = make(map[int]struct {
			result1 worker.Worker
			result2 bool
			result3 error
		})
	}
	fake.getWorkerReturnsOnCall[i] = struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) FindWorkerForResourceCheckContainer(logger lager.Logger, teamID int, resourceUser dbng.ResourceUser, resourceType string, resourceSource atc.Source, types atc.VersionedResourceTypes) (worker.Worker, bool, error) {
	fake.findWorkerForResourceCheckContainerMutex.Lock()
	ret, specificReturn := fake.findWorkerForResourceCheckContainerReturnsOnCall[len(fake.findWorkerForResourceCheckContainerArgsForCall)]
	fake.findWorkerForResourceCheckContainerArgsForCall = append(fake.findWorkerForResourceCheckContainerArgsForCall, struct {
		logger         lager.Logger
		teamID         int
		resourceUser   dbng.ResourceUser
		resourceType   string
		resourceSource atc.Source
		types          atc.VersionedResourceTypes
	}{logger, teamID, resourceUser, resourceType, resourceSource, types})
	fake.recordInvocation("FindWorkerForResourceCheckContainer", []interface{}{logger, teamID, resourceUser, resourceType, resourceSource, types})
	fake.findWorkerForResourceCheckContainerMutex.Unlock()
	if fake.FindWorkerForResourceCheckContainerStub != nil {
		return fake.FindWorkerForResourceCheckContainerStub(logger, teamID, resourceUser, resourceType, resourceSource, types)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.findWorkerForResourceCheckContainerReturns.result1, fake.findWorkerForResourceCheckContainerReturns.result2, fake.findWorkerForResourceCheckContainerReturns.result3
}

func (fake *FakeWorkerProvider) FindWorkerForResourceCheckContainerCallCount() int {
	fake.findWorkerForResourceCheckContainerMutex.RLock()
	defer fake.findWorkerForResourceCheckContainerMutex.RUnlock()
	return len(fake.findWorkerForResourceCheckContainerArgsForCall)
}

func (fake *FakeWorkerProvider) FindWorkerForResourceCheckContainerArgsForCall(i int) (lager.Logger, int, dbng.ResourceUser, string, atc.Source, atc.VersionedResourceTypes) {
	fake.findWorkerForResourceCheckContainerMutex.RLock()
	defer fake.findWorkerForResourceCheckContainerMutex.RUnlock()
	return fake.findWorkerForResourceCheckContainerArgsForCall[i].logger, fake.findWorkerForResourceCheckContainerArgsForCall[i].teamID, fake.findWorkerForResourceCheckContainerArgsForCall[i].resourceUser, fake.findWorkerForResourceCheckContainerArgsForCall[i].resourceType, fake.findWorkerForResourceCheckContainerArgsForCall[i].resourceSource, fake.findWorkerForResourceCheckContainerArgsForCall[i].types
}

func (fake *FakeWorkerProvider) FindWorkerForResourceCheckContainerReturns(result1 worker.Worker, result2 bool, result3 error) {
	fake.FindWorkerForResourceCheckContainerStub = nil
	fake.findWorkerForResourceCheckContainerReturns = struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) FindWorkerForResourceCheckContainerReturnsOnCall(i int, result1 worker.Worker, result2 bool, result3 error) {
	fake.FindWorkerForResourceCheckContainerStub = nil
	if fake.findWorkerForResourceCheckContainerReturnsOnCall == nil {
		fake.findWorkerForResourceCheckContainerReturnsOnCall = make(map[int]struct {
			result1 worker.Worker
			result2 bool
			result3 error
		})
	}
	fake.findWorkerForResourceCheckContainerReturnsOnCall[i] = struct {
		result1 worker.Worker
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) FindContainerForIdentifier(arg1 worker.Identifier) (db.SavedContainer, bool, error) {
	fake.findContainerForIdentifierMutex.Lock()
	ret, specificReturn := fake.findContainerForIdentifierReturnsOnCall[len(fake.findContainerForIdentifierArgsForCall)]
	fake.findContainerForIdentifierArgsForCall = append(fake.findContainerForIdentifierArgsForCall, struct {
		arg1 worker.Identifier
	}{arg1})
	fake.recordInvocation("FindContainerForIdentifier", []interface{}{arg1})
	fake.findContainerForIdentifierMutex.Unlock()
	if fake.FindContainerForIdentifierStub != nil {
		return fake.FindContainerForIdentifierStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.findContainerForIdentifierReturns.result1, fake.findContainerForIdentifierReturns.result2, fake.findContainerForIdentifierReturns.result3
}

func (fake *FakeWorkerProvider) FindContainerForIdentifierCallCount() int {
	fake.findContainerForIdentifierMutex.RLock()
	defer fake.findContainerForIdentifierMutex.RUnlock()
	return len(fake.findContainerForIdentifierArgsForCall)
}

func (fake *FakeWorkerProvider) FindContainerForIdentifierArgsForCall(i int) worker.Identifier {
	fake.findContainerForIdentifierMutex.RLock()
	defer fake.findContainerForIdentifierMutex.RUnlock()
	return fake.findContainerForIdentifierArgsForCall[i].arg1
}

func (fake *FakeWorkerProvider) FindContainerForIdentifierReturns(result1 db.SavedContainer, result2 bool, result3 error) {
	fake.FindContainerForIdentifierStub = nil
	fake.findContainerForIdentifierReturns = struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) FindContainerForIdentifierReturnsOnCall(i int, result1 db.SavedContainer, result2 bool, result3 error) {
	fake.FindContainerForIdentifierStub = nil
	if fake.findContainerForIdentifierReturnsOnCall == nil {
		fake.findContainerForIdentifierReturnsOnCall = make(map[int]struct {
			result1 db.SavedContainer
			result2 bool
			result3 error
		})
	}
	fake.findContainerForIdentifierReturnsOnCall[i] = struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) GetContainer(arg1 string) (db.SavedContainer, bool, error) {
	fake.getContainerMutex.Lock()
	ret, specificReturn := fake.getContainerReturnsOnCall[len(fake.getContainerArgsForCall)]
	fake.getContainerArgsForCall = append(fake.getContainerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetContainer", []interface{}{arg1})
	fake.getContainerMutex.Unlock()
	if fake.GetContainerStub != nil {
		return fake.GetContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.getContainerReturns.result1, fake.getContainerReturns.result2, fake.getContainerReturns.result3
}

func (fake *FakeWorkerProvider) GetContainerCallCount() int {
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	return len(fake.getContainerArgsForCall)
}

func (fake *FakeWorkerProvider) GetContainerArgsForCall(i int) string {
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	return fake.getContainerArgsForCall[i].arg1
}

func (fake *FakeWorkerProvider) GetContainerReturns(result1 db.SavedContainer, result2 bool, result3 error) {
	fake.GetContainerStub = nil
	fake.getContainerReturns = struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) GetContainerReturnsOnCall(i int, result1 db.SavedContainer, result2 bool, result3 error) {
	fake.GetContainerStub = nil
	if fake.getContainerReturnsOnCall == nil {
		fake.getContainerReturnsOnCall = make(map[int]struct {
			result1 db.SavedContainer
			result2 bool
			result3 error
		})
	}
	fake.getContainerReturnsOnCall[i] = struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorkerProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.runningWorkersMutex.RLock()
	defer fake.runningWorkersMutex.RUnlock()
	fake.getWorkerMutex.RLock()
	defer fake.getWorkerMutex.RUnlock()
	fake.findWorkerForResourceCheckContainerMutex.RLock()
	defer fake.findWorkerForResourceCheckContainerMutex.RUnlock()
	fake.findContainerForIdentifierMutex.RLock()
	defer fake.findContainerForIdentifierMutex.RUnlock()
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeWorkerProvider) recordInvocation(key string, args []interface{}) {
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

var _ worker.WorkerProvider = new(FakeWorkerProvider)
