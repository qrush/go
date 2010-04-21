/*

This program does operations on a list of unrealistic "wagons" on a prairie. It is initialized with two wagons, a head and a tail. The size of the prairie defaults to 30x30 and can be set via command line arguments as well.

Once running, keystrokes from the user move and modify the wagon train. The keys u, d, l, and r move the head, while U, D, L, and R move the tail. When the head moves, the next wagon in the train takes the head's old position and all the other wagons move to the previous position of their leader. Moving the tail works similarly.

Pressing the "a" key creates a new wagon in the upper left corner of the prairie, which becomes the new head. Pressing "A" makes a new wagon in the lower right corner to act as the new tail.

The "q" key quits the program.

*/
package documentation
