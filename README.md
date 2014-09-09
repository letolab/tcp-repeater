tcp-repeater
============

Binary which listens on source host and copies all data to multiple destinations

## Build

Run:

```
$ git clone https://github.com/letolab/tcp-repeater.git
$ cd tcp-repeater
$ go build
```

## Usage

The following listens on source `localhost:9000` and repeats data to all destinations: `locahost:9001` and `localhost:9002`

`./tcp-repeater -s localhost:9000 -d localhost:9001 -d localhost:9002`

## Caveats

If any one of the destinations close the connection prematurely, it will fail to copy data to all the other destinations.
