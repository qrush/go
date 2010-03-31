/*

The dag package implements a Directed Acyclic Graph according to the Set and Target interfaces. The command bin/dag uses the dag package to parse and traverse mkfiles.

Running it with no arguments loads a file called "mkfile" in the current directory and executes the first target. Giving it the "-f filename" option allows specification of the file name. Any other arguments are taken to be targets which should be executed.

When a mkfile is read and executed, the various targets are traversed in a depth-first fashion; thus, for the example below, targetD would be applied before targetC, since targetC depends on targetD.

Example mkfile:

	targetA	targetB	targetC

	targetB
		additional strings (currently ignored)

	targetC	targetD

	targetD

Example command line, default options:
	% bin/dag
	targetD
	targetC targetD
	targetB
	targetA targetC targetB
	%


*/
package documentation
