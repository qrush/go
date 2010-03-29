///////////////////////////////////////////////////////////////////////////////
// dagmake
// John Floren, Nick Quaranto
///////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"dag"
	"os"
)

func main() {
	fmt.Println("in main")
	dag.Main(func(s dag.Set, strs []string, fac dag.TargetFactory) (dag.Target, os.Error) {
		return nil, nil
	},
		func(t dag.Target) os.Error { return nil })
}
