<?php

fwrite(STDOUT, "This is a normal message\n");
fwrite(STDERR, "This is an error message\n");

// TODO This output is always emptp. We would love for PHP to see the ENV data, because it is often used in docker containers.
print_r($_ENV);