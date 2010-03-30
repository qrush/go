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
	dag.Main(dag.DagTargetFactory,
		func(t dag.Target) os.Error {
			fmt.Println("LOLOLOL: " + t.Name())
			return nil
		})
}
