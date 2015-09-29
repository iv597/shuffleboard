# Shuffleboard

Meet Shuffleboard. Shuffleboard likes to break your async code that forgot to
account for AJAX requests not necessarily returning in the order they were sent
out, or at the very least provide some realism to your dev environment, where
net requests are probably handled instantly and defensive asynchronous
programming can be easy to slack off on (or forget entirely).

## Usage
```
usage: shuffleboard [<flags>] [<command>...]

Flags:
  --help              Show context-sensitive help (also try --help-long and
                      --help-man).
  -c, --count=1       number of parallel executions - if your application is
                      asynchronous, the default of 1 is safe
  -b, --bind="localhost"
                      bind address (IP/hostname)
  -p, --port=8005     port to listen on
  -P, --innerPorts=INNERPORTS
                      comma-separated list (length of `count`) of ports to use
                      for spawned processes
  -a, --taskAddress="localhost"
                      address the spawned tasks are listening on
  -s, --taskSwitchLogic=1
                      logic to use for selecting which spawned process should
                      receive the request: 0 for sequential (NOT IMPLEMENTED), 1
                      for random
  -w, --minWait=0     the shortest (in ms) a request should be delayed
  -W, --maxWait=2500  the longest (in ms) a request should be delayed
  --version           Show application version.

Args:
  [<command>]  task to shuffle
```

`innerPorts` will default to a list of length `--count`, incrementing one port
number for each entry, starting at one port number above the Shuffleboard
port.

For example, if you have an asynchronous application which is hard-coded to
listen on something.local:4000, you might use the following:

`shuffleboard -a something.local -P 4000 node main.js`

Alternatively, let's say your application is a synchronous Flask application:

```
# Run four instances that know how to handle the $PORT env var natively
shuffleboard -a something.local -c 4 python main.py
```

## Legal

Shuffleboard is released by SpotOn, Inc. under the MIT License. See `LICENSE`
for details.
