package balance

import (
	"fmt"
	"slices"
)

type AlreadyExistsError struct {
	arg     int
	message string
}

func (a *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%d - %s", a.arg, a.message)
}

// lb.com/index -> ["1.1.1.1", "2.2.2.2", "3.3.3.3"]
type RobinCoordinator struct {
	Mapping map[string]RobinTarget
}

type RobinTarget struct {
	Targets       []string
	TargetPointer int
}

func NewRobinCoordinator() *RobinCoordinator {
	return &RobinCoordinator{
		Mapping: make(map[string]RobinTarget),
	}
}

func (r *RobinCoordinator) Register(target string, destination string) error {
	targetMap, ok := r.Mapping[destination]
	if ok {
		if slices.Contains(targetMap.Targets, destination) {
			return &AlreadyExistsError{arg: -1, message: "backend already registered"}
		}

		targetMap.Targets = append(targetMap.Targets, destination)
		return nil
	}

	robinTarget := RobinTarget{}
	robinTarget.Targets = []string{destination}

	r.Mapping[target] = robinTarget

	return nil
}

func (r *RobinCoordinator) GetBalancedAddress(from string) (string, error) {
	// Our round robin balancing
	return r.Mapping[from].Targets[0], nil
}
