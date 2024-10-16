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
	./k6 run examples/bug_application_properties.js

local/compose/up:
	docker compose up rabbit -d
	docker compose up sbemulator -d
	@export CONNECTION_STRING="Endpoint=sb://sbemulator/;SharedAccessKeyName=all;SharedAccessKey=CLwo3FQ3S39Z4pFOQDefaiUd1dSsli4XOAj3Y9Uh1E=;EnableAmqpLinkRedirect=false"

local/compose/down:
	docker compose down

test/run/local: install compile local/compose/up test/run local/compose/down
