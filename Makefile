install:
	go install go.k6.io/xk6/cmd/xk6@latest

compile:
	xk6 build --with xk6-azservicebus=.

test/run:
	./k6 run examples/queue_batch_message.js
	./k6 run examples/queue_batch_string_message.js
	./k6 run examples/queue_single_message.js
	./k6 run examples/queue_single_string_message.js
	./k6 run examples/topic_batch_message.js
	./k6 run examples/topic_batch_string_message.js
	./k6 run examples/topic_single_message.js
	./k6 run examples/topic_single_string_message.js

