// The MIT License (MIT)
//
// Copyright (c) 2015 Jamie Alquiza
//
// http://knowyourmeme.com/memes/deal-with-it.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	endpoints = [][]string{
		[]string{"cluster-health", "_cluster/health"},
		[]string{"cluster-stats", "_cluster/stats"},
	}

	nodeIp         string
	nodePort       string
	updateInterval int
	requireMaster  bool
	graphiteIp     string
	graphitePort   string
	metricsPrefix  string

	stats   = make(map[string][]byte)
	metrics = make(map[string]int64)

	metricsChan = make(chan map[string]int64, 30)
)

func init() {
	hostname, _ := os.Hostname()
	flag.StringVar(&nodeIp, "ip", "127.0.0.1", "ElasticSearch IP address")
	flag.StringVar(&nodePort, "port", "9200", "ElasticSearch port")
	flag.IntVar(&updateInterval, "interval", 30, "update interval")
	flag.BoolVar(&requireMaster, "require-master", false, "Only poll if node is an elected master")
	flag.StringVar(&graphiteIp, "graphite-ip", "", "Destination Graphite IP address")
	flag.StringVar(&graphitePort, "graphite-port", "", "Destination Graphite plaintext port")
	flag.StringVar(&metricsPrefix, "metrics-prefix", hostname, "Top-level Graphite namespace prefix (defaults to hostname)")
	flag.Parse()
}

func pollEs(nodeName string) {
	pollInt := time.Tick(time.Duration(updateInterval) * time.Second)
	for _ = range pollInt {
		switch requireMaster {
		case false:
			m, err := fetchMetrics()
			if err != nil {
				log.Println(err)
			} else {
				metricsChan <- m
			}
		case true:
			masterName, err := getMasterName()
			if err != nil {
				log.Println(err)
			}
			if nodeName != masterName {
				log.Println("Node is not an elected master")
			} else {
				m, err := fetchMetrics()
				if err != nil {
					log.Println(err)
				} else {
					metricsChan <- m
				}
			}
		}
	}
}

func handleMetrics() {
	for {
		// Connect to Graphite.
		graphite, err := net.Dial("tcp", graphiteIp+":"+graphitePort)
		if err != nil {
			log.Printf("Graphite unreachable: %s", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Ship metrics.
		metrics := <-metricsChan
		log.Println("Metrics received")

		ts := metrics["timestamp"]
		delete(metrics, "timestamp")

		for k, v := range metrics {
			_, err := fmt.Fprintf(graphite, "%s.%s %d %d\n", metricsPrefix, k, v, ts)
			if err != nil {
				log.Printf("Error flushing to Graphite: %s", err)
			}
		}

		log.Println("Metrics flushed to Graphite")
		graphite.Close()
	}
}

func fetchMetrics() (map[string]int64, error) {
	for i := range endpoints {
		key, endpoint := endpoints[i][0], endpoints[i][1]

		resp, err := queryEndpoint(endpoint)
		if err != nil {
			return nil, err
		}

		stats[key] = resp
	}

	json.Unmarshal(stats["cluster-stats"], &clusterStats)
	json.Unmarshal(stats["cluster-health"], &clusterHealth)

	now := time.Now()
	ts := int64(now.Unix())
	metrics["timestamp"] = ts
	metrics["es-stats.state.red"] = 0
	metrics["es-stats.state.yellow"] = 0
	metrics["es-stats.state.green"] = 0
	// Flip value according to read state.
	metrics["es-stats.state."+clusterHealth.Status] = 1

	metrics["es-stats.shards.active_primary_shards"] = clusterHealth.ActivePrimaryShards
	metrics["es-stats.shards.active_shards"] = clusterHealth.ActiveShards
	metrics["es-stats.shards.relocating_shards"] = clusterHealth.RelocatingShards
	metrics["es-stats.shards.initializing_shards"] = clusterHealth.InitializingShards
	metrics["es-stats.shards.unassigned_shards"] = clusterHealth.UnassignedShards

	metrics["es-stats.indices"] = clusterStats.Indices.Count
	metrics["es-stats.docs"] = clusterStats.Indices.Docs.Count
	metrics["es-stats.cluster_cpu_cores"] = clusterStats.Nodes.Os.AvailableProcessors
	metrics["es-stats.cluster_memory"] = clusterStats.Nodes.Os.Mem.TotalInBytes

	metrics["es-stats.nodes.master"] = clusterStats.Nodes.Count.MasterOnly
	metrics["es-stats.nodes.data"] = clusterStats.Nodes.Count.DataOnly
	metrics["es-stats.nodes.master_data"] = clusterStats.Nodes.Count.MasterData
	metrics["es-stats.nodes.client"] = clusterStats.Nodes.Count.Client

	metrics["es-stats.fs.total"] = clusterStats.Nodes.Fs.TotalInBytes
	metrics["es-stats.fs.available"] = clusterStats.Nodes.Fs.AvailableInBytes
	storageUsed := metrics["es-stats.fs.total"] - metrics["es-stats.fs.available"]
	metrics["es-stats.fs.used"] = storageUsed

	metrics["es-stats.mem.jvm.heap_used_in_bytes"] = clusterStats.Nodes.Jvm.Mem.HeapUsedInBytes
	metrics["es-stats.mem.jvm.heap_max_in_bytes"] = clusterStats.Nodes.Jvm.Mem.HeapMaxInBytes
	metrics["es-stats.mem.store.size_in_bytes"] = clusterStats.Indices.Store.SizeInBytes
	metrics["es-stats.mem.store.throttle_time_in_millis"] = clusterStats.Indices.Store.ThrottleTimeInMillis
	metrics["es-stats.mem.fielddata.memory_size_in_bytes"] = clusterStats.Indices.Fielddata.MemorySizeInBytes
	metrics["es-stats.mem.fielddata.evictions"] = clusterStats.Indices.Fielddata.Evictions
	metrics["es-stats.mem.filter_cache.memory_size_in_bytes"] = clusterStats.Indices.FilterCache.MemorySizeInBytes
	metrics["es-stats.mem.filter_cache.evictions"] = clusterStats.Indices.FilterCache.Evictions
	metrics["es-stats.mem.id_cache.memory_size_in_bytes"] = clusterStats.Indices.IdCache.MemorySizeInBytes
	metrics["es-stats.mem.completion.size_in_bytes"] = clusterStats.Indices.Completion.SizeInBytes
	metrics["es-stats.mem.segments.count"] = clusterStats.Indices.Segments.Count
	metrics["es-stats.mem.segments.memory_in_bytes"] = clusterStats.Indices.Segments.MememoryInBytes
	metrics["es-stats.mem.segments.index_writer_memory_in_bytes"] = clusterStats.Indices.Segments.IndexWriterMemoryInBytes
	metrics["es-stats.mem.segments.index_writer_max_memory_in_bytes"] = clusterStats.Indices.Segments.IndexWriterMaxMemoryInBytes
	metrics["es-stats.mem.segments.version_map_memory_in_bytes"] = clusterStats.Indices.Segments.VersionMapMemoryInBytes
	metrics["es-stats.mem.segments.fixed_bit_set_memory_in_bytes"] = clusterStats.Indices.Segments.FixedBitSetMemoryInBytes

	return metrics, nil
}

func queryEndpoint(endpoint string) ([]byte, error) {
	resp, err := http.Get("http://" + nodeIp + ":" + nodePort + "/" + endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
func getNodeName() (string, error) {
	resp, err := http.Get("http://" + nodeIp + ":" + nodePort + "/_nodes/_local/name")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(contents, &nodesLocal)

	var name string
	for k, _ := range nodesLocal.Nodes.(map[string]interface{}) {
		name = k
	}
	return name, nil
}

func getMasterName() (string, error) {
	resp, err := http.Get("http://" + nodeIp + ":" + nodePort + "/_cluster/state/master_node")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(contents, &clusterState)

	return clusterState.MasterNode, nil
}

func main() {
	// Grab node name.
	var nodeName *string
	retry := time.Tick(time.Duration(updateInterval) * time.Second)

	for _ = range retry {
		name, err := getNodeName()
		if err != nil {
			log.Printf("ElasticSearch unreachable: %s", err)

		} else {
			nodeName = &name
			log.Printf("Connected to ElasticSearch: http://%s:%s\n", nodeIp, nodePort)
			break
		}
	}

	// Run.
	go handleMetrics()
	pollEs(*nodeName)
}
