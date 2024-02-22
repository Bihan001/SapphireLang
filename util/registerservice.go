package util

import (
	"fmt"
	"strconv"
)

type RegisterService struct {
	registerMap          map[string]bool
	lowerByteRegisterMap map[string]string
}

var allocator *RegisterService

func GetNewRegisterService() *RegisterService {
	if allocator == nil {
		allocator = &RegisterService{
			registerMap:          map[string]bool{"r8": false, "r9": false, "r10": false, "r11": false, "r12": false, "r13": false, "r14": false, "r15": false},
			lowerByteRegisterMap: map[string]string{"r8": "r8b", "r9": "r9b", "r10": "r10b", "r11": "r11b", "r12": "r12b", "r13": "r13b", "r14": "r14b", "r15": "r15b"},
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

func (a *RegisterService) GetLowerByte(register string) string {
	val, ok := a.lowerByteRegisterMap[register]
	if !ok {
		panic(fmt.Sprintf("cannot find lower byte of %s", register))
	}
	return val
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
