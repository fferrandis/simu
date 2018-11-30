# simu
simple simulator to manage data dispatch across several servers and disks

## Build and install
```
export GOPATH=<path_to_repository>
go install sim
```

## start simulator 
```
./bin/sim [--config=<path_to_config_file>]
```

## Send some data and get time needed to perform the IO
```
 curl    -XPUT   'localhost:8080/put?datalen=134100000'; 
{"scal-response-time" : 1341000000}

```
The answer contains a simple json giving us the time needed to perform the IO
This value is a nsec number, so in that example
{"scal-response-time" : 1341000000} is 1s and 341ms


## Configuration

A default configuration is declared in config.go file. Those values can be changed.
You can directly supplied a simple json, with that format :

```json
{
  "write_speed"     : 100000000,
  "read_speed"      : 200000000,
  "extent_size"     : 134200000,
  "data_scheme"     :5,
  "coding_scheme"   : 2,
  "network_bdwidth" : 125000000,
    "hdservers" : [
      {
        "nr_disk" : 20,
        "capacity": 5000000000
      },
      {
        "nr_disk" : 20,
        "capacity": 5000000000
      }
    ]
}

```
* write\_speed is the speed in b/s to perform writes. This is global to all disk
* read\_speed is the same but for read operation
* extent\_size  is the size of a container that receive ata
* {data,coding}\_scheme is the ECE schema used
* network\_bdwidth is the network upload speed in bytes/s
* hdservers section describe all the servers you will declare.  
  - nr\_disk is the number od disk for one server
  - capacity is size in bytes


## Time IO computation

The simulator assume to receive data in a continuous way.
That means that, with the configuration used above, receiving 125000000 bytes
will advance our internal timer to 1 second (since the network bandwith has a capacity
of 125000000 bytes by second)
So when we start the simulator, the internal timer is at 0.
After receiving 125000000b, the internal timer is at 1 000 000 000 ns (1s)
And so on ....

When a put is routed to a disk
