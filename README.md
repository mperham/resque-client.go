resque-client.go
=======================

Resque is a nice, simple queueing system built on top of Redis's
native list structure but its workers are horribly inefficient since
they are single threaded.  This client can read jobs pushed from a
Rails app and process them in parallel using goroutines.

Installation
------------------

None just yet.  Still in development.


Notes
--------------

Contains a forked copy of the Go-Redis driver in the redis/ directory
since I couldn't get it to install normally.
