/*

Expr recreates basic functionality of the Unix "expr" command.
It is implemented as a recursive descent parser, as explained on http://math.hws.edu/javanotes/c9/s5.html. The following BNF, from the same source, was used:
	<expression>  ::=  [ "-" ] <term> [ ( "+" | "-" ) <term> ]...
	<term>  ::=  <factor> [ ( "*" | "/" ) <factor> ]...
	<factor>  ::=  <number>  |  "(" <expression> ")"

This expr understands +,-,*,/, and parenthesized expressions. Examples:
	% ./expr 1 + 2
	3
	% ./expr \( 5 + 1 \) \* 2
	12
	%./expr 4 \* 2 \* 2 + 1
	17

*/
package documentation
