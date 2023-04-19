package mapreduce

import (
	"encoding/json"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Use checkError to handle errors.

	// To store ALL the maps aka keyValue pairs
	var intermediate []KeyValue

	/*
		Opening all the Encoded files created by doMap function
		and decoding them in-order to insert them into the
		intermediate array
	*/
	for i := 0; i < nMap; i++ {
		file, err := os.Open(reduceName(jobName, i, reduceTaskNumber))
		checkError(err)

		dec := json.NewDecoder(file)
		// Decoding the files until an error aka EOF is encountered
		for {
			var kv KeyValue
			err := dec.Decode(&kv)
			if err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
		file.Close()
	}

	// Sort intermediate key/value pairs by key
	// (Since key is a string, it sort them lexicographically)
	sort.Slice(intermediate, func(i, j int) bool {
		return intermediate[i].Key < intermediate[j].Key
	})

	// Creates a merge file by using mergeName for the file name
	outFile, err := os.Create( mergeName(jobName, reduceTaskNumber))
	checkError(err)

	enc := json.NewEncoder(outFile)

	/*
		Iterate over intermediate's keyValue pairs to merge all
		with the same keys and insert them into the new merge file
	*/
	i := 0
	for i < len(intermediate) {
		// Using j to count the number of values in Intermediate with the same
		// key value for a specific key
		j := i + 1
		for j < len(intermediate) && intermediate[i].Key == intermediate[j].Key {
			j++
		}

		// Using j, appending j values for the reduce function
		values := make([]string, 0)
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}

		// Calling the reduce function to reduce all the keys with same
		// keys in-order to merge them and encode them into the new file
		enc.Encode(KeyValue{intermediate[i].Key, reduceF(intermediate[i].Key, values)})
		i = j
	}

	outFile.Close()

	// returns nothing other than the one new merged file

}
