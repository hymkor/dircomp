dircomp
=======

Compare two directory like `svn status`.

Executing `dircomp X Y`, the differences between X and Y are output as follows.

- `A  fname` ... X/fname does not exist, but Y/fname exists.
- `M  fname` ... Both X/fname and Y/fname exist, but they are different.
- `D  fname` ... X/fname exists, but Y/fname does not exist.

If X/fname and Y/fname equal to each other, they are not list up.

Option:

- `-i WILDCARD` ... filtering
