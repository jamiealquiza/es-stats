// The MIT License (MIT)
//
// Copyright (c) 2015 Jamie Alquiza
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

var nodesLocal struct {
	Nodes interface {
	} `json:"nodes"`
}

var clusterState struct {
	MasterNode string `json:"master_node"`
}

var clusterHealth struct {
	Status              string `json:"status"`
	ActivePrimaryShards int64  `json:"active_primary_shards"`
	ActiveShards        int64  `json:"active_shards"`
	RelocatingShards    int64  `json:"relocating_shards"`
	InitializingShards  int64  `json:"initializing_shards"`
	UnassignedShards    int64  `json:"unassigned_shards"`
}

var clusterStats struct {
	Indices struct {
		Count int64
		Docs  struct {
			Count int64 `json:"count"`
		} `json:"docs"`
		Store struct {
			SizeInBytes          int64 `json:"size_in_bytes"`
			ThrottleTimeInMillis int64 `json:"throttle_time_in_millis"`
		} `json:"store"`
		Fielddata struct {
			MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
			Evictions         int64 `json:"evictions"`
		} `json:"fielddata"`
		FilterCache struct {
			MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
			Evictions         int64 `json:"evictions"`
		} `json:"filter_cache"`
		IdCache struct {
			MemorySizeInBytes int64 `json:"memory_size_in_bytes"`
		} `json:"id_cache"`
		Completion struct {
			SizeInBytes int64 `json:"size_in_bytes"`
		} `json:"completion"`
		Segments struct {
			Count                       int64 `json:"count"`
			MememoryInBytes             int64 `json:"memory_in_bytes"`
			IndexWriterMemoryInBytes    int64 `json:"index_writer_memory_in_bytes"`
			IndexWriterMaxMemoryInBytes int64 `json:"index_writer_max_memory_in_bytes"`
			VersionMapMemoryInBytes     int64 `json:"version_map_memory_in_bytes"`
			FixedBitSetMemoryInBytes    int64 `json:"fixed_bit_set_memory_in_bytes"`
		} `json:"segments"`
	} `json:"indices"`
	Nodes struct {
		Count struct {
			MasterOnly int64 `json:"master_only"`
			DataOnly   int64 `json:"data_only"`
			MasterData int64 `json:"master_data"`
			Client     int64 `json:"client"`
		} `json:"count"`
		Os struct {
			AvailableProcessors int64 `json:"available_processors"`
			Mem                 struct {
				TotalInBytes int64 `json:"total_in_bytes"`
			} `json:"mem"`
		} `json:"os"`
		Jvm struct {
			Mem struct {
				HeapUsedInBytes int64 `json:"heap_used_in_bytes"`
				HeapMaxInBytes  int64 `json:"heap_max_in_bytes"`
			} `json:"mem"`
		} `json:"jvm"`
		Fs struct {
			TotalInBytes     int64 `json:"total_in_bytes"`
			AvailableInBytes int64 `json:"available_in_bytes"`
		} `json:"fs"`
	} `json:"nodes"`
}
