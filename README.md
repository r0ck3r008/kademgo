## Kademlia DHT Implementation

The project aims to adhere to vanilla Kademlia information protocol as described by _Maymounkov et. al._ in the paper __Kademlia: A Peer to Peer Information System based on XOR Metric__

#### Dependencies
* go
* godoc (optional: for documentation)

#### Usage
The kademgo library can be imported in any project like regular libraries.
In the project root run,
```bash
$> go get github.com/r0ck3r008/kademgo
```
or, to get the conainer image,
```bash
$> podman pull ghcr.io/r0ck3r008/kademgo:latest # using podman (recommended)
$> docker pull ghcr.io/r0ck3r008/kademgo:latest # Using docker
```
Then use,
```go
import "github.com/r0ck3r008/kademgo"
```
in the project.
The project has actively mantained [godoc](https://blog.golang.org/godoc) based comments. This makes easier for community to figure out how code is structured.

To access the documentation,
```bash
# Clone the repository
git clone https://github.com/r0ck3r008/kademgo
# Go to directory root
cd kademgo
# Use Godoc HTTP server
godoc
```
Then in your browser, visit <http://localhost:6060/pkg/github.com/r0ck3r008/kademgo>.

#### Contributing
This project aims to follow all the general guidelines mentioned in official go documentation in [Effective Go](https://golang.org/doc/effective_go.html).
This is project currently under development and is open for all contributions.

#### Future Trajectory
The future trajectory of this project is to make a _Command and Control_ framework for a botnet inspired by _Starnberger et. al_ in the paper __Overbot: a botnet protocol based on Kademlia__.
