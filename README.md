# Telnet client

This is a simple telnet client written on pure [golang](https://golang.org/)

The implementation generally follows specifications from [OTUS basic telnet client task](https://github.com/OtusGolang/home_work/tree/master/hw11_telnet_client) 


#### Call examples
```bash
$ go-telnet --timeout=10s host port
$ go-telnet mysite.ru 8080
$ go-telnet --timeout=3s 1.1.1.1 123
```


#### Connection examples

1) server closes connection
```bash
$ nc -l localhost 4242
Hello from NC
I'm telnet client
Bye, client!          
^C
```

```bash
$ go-telnet --timeout=5s localhost:4242
Hello from NC
I'm telnet client
Bye, client!
Bye-bye 
```

2) client closes connection
```bash
$ go-telnet localhost:4242
I
will be
back!
^D
```

```bash
$ nc -l localhost 4242
I
will be
back!
```

