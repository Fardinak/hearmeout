# Hear Me Out

A very simple chat client/server that uses OS X's `say` command to read the message to you on arrival.

## Usage
On the server
```
$ hearmeout-server [ip:port]
```
* `[ip:port]` defaults to 0.0.0.0:35000

On the clients
```
$ hearmeout-client hostname:port [username]
```
* `hostname:port` is mandatory
* `username` defaults to current system user's username

## DISCLAIMER
While being simple and pretty cool, it's a little rude. So you may wanna modify error messages before using it!

## License
None
