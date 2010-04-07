/*

Dag reads a dependency graph from a file and performs
depth-first traversals.

Usage:
  dag [-f file] target...
  
The default file is mkfile, the default target is the first 
target in the file.

The file contains blocks for targets, separated by one or more
blank lines.

The first line of a block contains a target and optionally
prerequisite targets, separated by white space.

Dag is intended as a building block for a make-like command.

*/
package documentation
