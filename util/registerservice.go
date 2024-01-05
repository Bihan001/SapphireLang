package util

import (
	"strconv"
)

type RegisterService struct {
	registerMap map[string]bool
}

var allocator *RegisterService

func GetNewRegisterService() *RegisterService {
	if allocator == nil {
		allocator = &RegisterService{
			registerMap: map[string]bool{"r8": false, "r9": false, "r10": false, "r11": false, "r12": false, "r13": false, "r14": false, "r15": false},
		}
	}
	return allocator
}

func (a *RegisterService) GetNewRegister() string {
	for key, val := range a.registerMap {
		if val == false {
			a.registerMap[key] = true
			return key
		}
	}
	panic("failed to get new register")
}

func (a *RegisterService) FreeRegister(allc string) {
	res, ok := a.registerMap[allc]
	if !ok || res == false {
		panic("failed to free register")
	}
	a.registerMap[allc] = false
}

func (a *RegisterService) FreeAllRegisters() {
	for key, _ := range a.registerMap {
		a.registerMap[key] = false
	}
}

func (a *RegisterService) getVariableFromCounter(count int) string {
	return "%" + strconv.Itoa(count)
}
