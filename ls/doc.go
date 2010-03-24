/*

ls recreates basic functionality of the Unix "ls" command, but in a lisp-like programmable form. The command takes a file or directory and a script to apply to that file.

A simple script that will print the filename for a regular file and do nothing for a directory:
	( (file (name) (nl) ) )
Example usage:
	% ./ls somefile simplescript.ls
	somefile
	% ./ls somedir simplescript.ls
	%

A more complex script that displays the filename if given a file and recursively displays all the contents of a directory.
	( (file (name) (nl) )
		(dir
			(name) : (nl)
			(sub (dir (name) (nl) ) )
			(sub (file (name) (nl) ) )
			(sub (dir (recurse) ) ) 
		)
	)
Example usage:
	% ./ls somefile complexscript.ls
	somefile
	% ./ls test/complextree/ complexscript.ls
	test/complextree:
	test/complextree/eins
	test/complextree/gamma
	test/complextree/beta
	test/complextree/alpha
	test/complextree/eins:
	test/complextree/eins/one
	test/complextree/eins/foobar
	test/complextree/eins/rawr
	test/complextree/eins/one:
	test/complextree/eins/one/foodir
	test/complextree/eins/one/foodir:
	test/complextree/eins/one/foodir/blah
	%

The script must be enclosed at the top level with a pair of parentheses. The available functions are:
	(file ... ): If the current item is a file, execute the items in this expression.
	(dir ... ): If the current item is a directory, execute the items in this expression.
	(sub ... ): Loop over the direct descendants of a directory
	(recurse): Re-execute the entire script for the current item. 
	(tab): Insert a tab character
	(name): Print the current item's name
	(user): Print the numeric user ID of the current item
	(group): Print the numeric group ID of the current item
	(size): Print the size of the current item in bytes
	(human_size): Print the size of the current items in bytes/kilobytes/megabytes/gigabytes as necessary
	(nl): Print a newline

Any other token, such as ":" in the example above, is simply printed.

*/
package documentation
