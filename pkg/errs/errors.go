package custom_errors

import (
	"fmt"
)

type errs struct {
}

func InitErrs() errs {
	return errs{}
}

type Errs interface {
	NoMatches(game string) error
	NotAllMatchesFetched(number int) error
}

func (e errs) NoMatches(game string) error {
	return fmt.Errorf("no matches played in %v", game)
}

// maybe no need in this function
func (e errs) NotAllMatchesFetched(number int) error {
	return fmt.Errorf("info about first %v matches failed to fetch due to external reasons", number)
}
