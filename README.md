reverse tcp proxy

Tunneling raw tcp traffic from remote pc to your behind NAT local service and back.

USAGE: etunnel server|client path_to_config.json

REMOTE-PC <--> |publicPort <--> clientPort SERVER| <- (NAT) -> |CLIENT| <--> |servicePort SERVICE|

1. SERVER listens on localhost(publicHost):clientPort
2. CLIENT connects to publicHost:clientPort and log-in using digest method and waits. Repeat.
3. WhenREMOTE-PC connects to publicHost:publicPort SERVER activates CLIENT and CLIENT connects to servicePort
4. SERVICE and CLIENT transfer traffic between REMOTE-PC and SERVICE
