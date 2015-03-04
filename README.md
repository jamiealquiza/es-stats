# es-wtf

Reads key cluster metrics from ElasticSearch and writes to Graphite. 

Just built, needs thorough testing.

### Installation

es-wtf has no external dependencies (you're welcome). Assuming Go is installed (built/tested with 1.4.x):

- `go get github.com/jamiealquiza/go-wtf`
- `go build github.com/jamiealquiza/go-wtf`

Binary will be found at: `$GOPATH/bin/go-wtf`

### Usage

Flags:
<pre>
./es-wtf -h
Usage of ./es-wtf:
  -graphite-ip="": Destination Graphite IP address
  -graphite-port="2003": Destination Graphite plaintext port
  -interval=30: Metrics polling interval
  -ip="127.0.0.1": ElasticSearch IP address
  -metrics-prefix="elasticsearch": Top-level Graphite namespace prefix (defaults to hostname)
  -port="9200": ElasticSearch port
  -require-master=false: Only poll if node is an elected master
</pre>

Running:
<pre>
% ./es-wtf -ip="192.168.100.204" -interval=5 -graphite-ip="192.168.100.175" -graphite-port="2013"
2015/03/04 15:27:56 Connect to ElasticSearch: http://192.168.100.204:9200
2015/03/04 15:27:56 Connected to Graphite: 192.168.100.175 port 2013
2015/03/04 15:28:01 Metrics received
2015/03/04 15:28:01 Metrics flushed to Graphite
2015/03/04 15:28:06 Metrics received
2015/03/04 15:28:06 Metrics flushed to Graphite
2015/03/04 15:28:11 Metrics received
2015/03/04 15:28:11 Metrics flushed to Graphite
2015/03/04 15:28:16 Metrics received
2015/03/04 15:28:16 Metrics flushed to Graphite
</pre>

Get metrics:
![ScreenShot](http://us-east.manta.joyent.com/jalquiza/public/github/es-clusterstats-graphite.png)

Make pretty graphs like:
![ScreenShot](http://us-east.manta.joyent.com/jalquiza/public/github/es-stats-grafana.png)

and:
![ScreenShot](http://us-east.manta.joyent.com/jalquiza/public/github/es-clusterstats.png)
