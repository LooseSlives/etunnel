reverse tcp proxy
build command: go build .
Tunneling raw tcp traffic from remote pc to your behind NAT local service and back.
One server|client - one servicePort

USAGE: etunnel path_to_config.json

REMOTE-PC <--> |ServerPublicPort SERVER ServerClientPort| <- (NAT) -> |CLIENT| <--> |servicePort SERVICE|

1. SERVER listens on it's localhost (ServerPublicPort and ServerClientPort)
2. CLIENT connects to ServerHost:ServerClientPort and log-in using digest method and waits. Repeat to always have ClientReserve number of ready Client connections.
3. When REMOTE-PC connects to ServerHost:ServerPublicPort SERVER activates CLIENT and CLIENT connects to servicePort
4. SERVICE and CLIENT tunnelling all traffic between REMOTE-PC and SERVICE
