install:
	go install go.k6.io/xk6/cmd/xk6@latest

compile:
	xk6 build --with xk6-azservicebus=.

test/run:
	./k6 run tests/test.js
