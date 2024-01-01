package util

import (
	"strconv"
)

type IRVariableService struct {
	allocated_map map[string]bool
	free_stack    []string
	counter       int
}

var allocator *IRVariableService

func GetNewIRVariableService() *IRVariableService {
	if allocator == nil {
		allocator = &IRVariableService{
			allocated_map: make(map[string]bool),
			free_stack:    make([]string, 0),
			counter:       0,
		}
	}

	return allocator
}

/*
Original feature: If previously allocated variables are freed, then its stored in the free_stack and reused in future instead of creating a new variable by incrementing the counter
Current Feature: We just increment the counter and don't reuse the previously used variables from free_stack
Reason: The reason we don't reuse them is because LLVM IR doesn't support reassignment of variables
How to revert to original: Just uncomment the commented code and also change the condition from true to the commented one in GetNewAllocation method
*/

func (a *IRVariableService) GetNewAllocation() string {
	//if len(a.free_stack) == 0 {
	if true {
		val := a.getVariableFromCounter(a.counter)
		//a.allocated_map[val] = true
		a.counter += 1
		return val
	}

	allc := a.free_stack[0]
	a.free_stack = a.free_stack[1:]
	a.allocated_map[allc] = true
	return allc
}

func (a *IRVariableService) FreeAllocation(allc string) error {
	//res, ok := a.allocated_map[allc]
	//if !ok || res == false {
	//	return errors.New(allc + " Not allocated!")
	//}
	//a.allocated_map[allc] = false
	//a.free_stack = append(a.free_stack, allc)
	return nil
}

func (a *IRVariableService) FreeAllAllocation() {

}

func (a *IRVariableService) getVariableFromCounter(count int) string {
	return "%" + strconv.Itoa(count)
}
