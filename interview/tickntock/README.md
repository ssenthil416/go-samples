tick tock and bong

Requirment
=========

Create a clock application written in go that will print the following values at the following intervals to stdout:

- "tick" every second

- "tock" every minute

- "bong" every hour
Only one value should be printed in a given second, i.e. when printing "bong" on the hour, the "tick" and "tock" values should not be printed.

It should run for three hours and then exit.

A mechanism should exist for the user to alter any of the printed values while the program is running, i.e. after the clock has run for 10 minutes I should, without stopping the program, be able to change it so that it stops printing "tick" every second and starts printing "quack" instead.

Please provide appropriate test coverage

Design
======
To change "tick" to "quack", custom.txt file is used for it. 

build & run
==========
go run tickntock.go

test
====
go run -v *.go
