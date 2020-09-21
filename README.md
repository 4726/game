# game

Implementation of a backend game server in Go and gRPC using microservice architecture.

# Architecture

![img](https://github.com/4726/game/blob/master/images/design.png?raw=true)

# Services

### Constants

Stores Key-Value constants.

### Custom Match

Allows players to create custom lobbies with password support.

### Match History

Stores finished match results.

### Live Matches

Stores matches currently going on.

### Matchmaking Queue

Allows players to join a queue to look for matches.

### Ranking

Stores player matchmaking ranking.

### Poll

Stores player matchmaking ranking.

### Chat

Stores direct messages between two players. Sends notification to NSQ when a new direct message is sent.

### Friends

Stores friend status between players.