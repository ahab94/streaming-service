# Streaming Service

__Streaming Service__  connects to a __Channel Source Service__ and transmit appropriate live channel streams to the connected __Users__.


```
                                   +------------------+
                                   |                  |
                                   |                  |/register HTTP(8080) POST
                                   |                  +----------------------------+
                                   |                  |endpoint for registering
                                   |                  |users
                                   |                  |
                                   |                  |
+---------------+                  |                  |
|               | TCP(9110)        |                  | TCP (9111)
|Channel source +<---------------->+ Streaming service+-----------------------------+
|               | Unordered stream |                  | endpoint for transmitting
+---------------+ of streaming     |                  | ordered streaming packets
                  packets for all  |                  | payload for requested
                  channels         |                  | channal id to a connected user
                                   |                  |
                                   |                  |
                                   |                  |
                                   |                  | /usage/{user_id} HTTP(8080) GET
                                   |                  +-----------------------------+
                                   +------------------+ endpoint for exposing usage
                                                        data in seconds

```


Streaming packets for any given channel could arrive out of order from the __Channel Source Service__ and should be ordered correctly before sending payload part to registered __Users__ watching these channels at the moment of receiving the streaming packet from __Channel Source Service__.

__Streaming Service__ reads those packets from channel source, process them and publish on the appropirate users.

## API

- Usage: Get last 24 hour of usage for a particular user
    - path: /usage/{user-id}
    - method: GET
    - statuses: 200, 500
    - example: 
        ```bash
        ahab@ahab-VirtualBox:~/work/streaming-service$ curl -i localhost:8080/usage/user-1
        HTTP/1.1 200 OK
        Date: Mon, 30 Mar 2020 08:09:59 GMT
        Content-Length: 39
        Content-Type: text/plain; charset=utf-8
        
        {"sessions":{"zato01":31,"zato02":31}} 
        ```

- Register: to stream packets, a user must be registered
    - path: /register
    - method: POST
    - statuses: 200, 400, 500
    - Body: user-id
        - just the user-id as part of text body
    - example: 
        ```bash
        ahab@ahab-VirtualBox:~/work/streaming-service$ curl -i -d 'user-1' localhost:8080/register
        HTTP/1.1 200 OK
        Date: Mon, 30 Mar 2020 08:19:06 GMT
        Content-Length: 0
        ```

## Development & Requirements
Makefile is being used for the development of __Streaming Service__.

To start contributing to __Stream Service__ you will need to have golang and docker installed.

__Stream Service__ builds and uses dockerfile.dev container for various development activities like
formatting, coverage, linting and testing.

When you run  `$ make any_target`, it builds the development container, it can take some time depending on your network speed.



```bash
Usage:
  ## Build Commands
    build           Compile the project.
  ## Demo Command
    demo            automated demonstration
  ## Develop / Test Commands
    dep             Update go modules.
    dep-verify      Verify go modules.
    dep-why         Question the need of imported go modules.
    dep-cache       cache go modules.
    format          Run code formatter.
    check           Run static code analysis (lint).
    test            Run tests on project.
    cover           Run tests and capture code coverage metrics on project.
    todo            Generate a TODO list for project.

  ## Local Commands
    drma            Removes all stopped containers.
    drmia           Removes all unlabelled images.
    drmvu           Removes all unused container volumes.

```   
## Deployment

__Streaming Service__ can be deployed as a binary or as a container it can be deployed on standalone Docker, Kubernetes, Openshift and other container schedulers.

### Configuration

```bash
	# env var for server
	"SERVER_HOST"
	"SERVER_PORT"

	# env var for stream server
	"STREAM_HOST"
	"STREAM_PORT"

	# env var for channel source
	"CHANNEL_SOURCE_HOST"
	"CHANNEL_SOURCE_PORT"

	# env var for packets batch size per read
	"PACKET_BATCH_SIZE"
```

## Demo

To run a quick demonstration of __Streaming Service__ just dial in `$ make demo`

```bash
ahab@ahab-VirtualBox:~/work/streaming-service$ make demo
starting db...
setting up db...
starting channel source
bash ./scripts/demo.sh
Starting Streaming Server

{"level":"info","msg":"Server is ready to handle tcp requests at 127.0.0.1:9111","time":"2020-03-30T12:24:08+05:00"}
{"level":"info","msg":"Server is ready to handle requests at 127.0.0.1:8080","time":"2020-03-30T12:24:08+05:00"}
going to interact with stream service...
user-5: registration succeed. Will stream soon...
{"level":"debug","msg":"setting up connection 127.0.0.1:49188","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"payload user-5 zato02 from 127.0.0.1:49188","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"processed get query: SELECT * FROM users WHERE user_id='user-5' ","time":"2020-03-30T12:24:08+05:00"}
user-4: registration succeed. Will stream soon...
user-2: registration succeed. Will stream soon...
user-3: registration succeed. Will stream soon...
user-1: registration succeed. Will stream soon...
{"level":"debug","msg":"setting up connection 127.0.0.1:49190","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"payload user-4 zato01 from 127.0.0.1:49190","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"processed get query: SELECT * FROM users WHERE user_id='user-4' ","time":"2020-03-30T12:24:08+05:00"}
{"level":"info","msg":"adding user user-5 for channel zato02","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"setting up connection 127.0.0.1:49192","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"payload user-2 zato02 from 127.0.0.1:49192","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"processed get query: SELECT * FROM users WHERE user_id='user-2' ","time":"2020-03-30T12:24:08+05:00"}
{"level":"info","msg":"adding user user-4 for channel zato01","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"setting up connection 127.0.0.1:49194","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"payload user-3 unknown from 127.0.0.1:49194","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"processed get query: SELECT * FROM users WHERE user_id='user-3' ","time":"2020-03-30T12:24:08+05:00"}
{"level":"info","msg":"adding user user-2 for channel zato02","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"setting up connection 127.0.0.1:49196","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"payload user-1 zato01 from 127.0.0.1:49196","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"processed get query: SELECT * FROM users WHERE user_id='user-1' ","time":"2020-03-30T12:24:08+05:00"}
{"level":"info","msg":"adding user user-3 for channel unknown","time":"2020-03-30T12:24:08+05:00"}
{"level":"info","msg":"adding user user-1 for channel zato01","time":"2020-03-30T12:24:08+05:00"}
{"level":"debug","msg":"sending 1501 packets to user user-1","time":"2020-03-30T12:24:23+05:00"}
{"level":"debug","msg":"sending 1501 packets to user user-4","time":"2020-03-30T12:24:23+05:00"}
{"level":"debug","msg":"sending 1499 packets to user user-2","time":"2020-03-30T12:24:23+05:00"}
{"level":"debug","msg":"sending 1499 packets to user user-5","time":"2020-03-30T12:24:23+05:00"}
```

Makefile target & demo script spawn cassandra, channel source,
streaming server and user demo (which runs 5 users that seek streaming for different channels).

## Having Issues?

Please feel free to reach out at 'abdshah94@gmail.com'.
