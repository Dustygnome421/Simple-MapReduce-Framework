package mapreduce

import "sync"

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
		// for i := 0; i < len(mr.workers); i++ {
		// 	mr.doneChannel <- true
		// }
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO:

	// ---------------------------

	// Buffered Channel to store the task numbers
	taskChannel := make(chan int, ntasks)
	for i := 0; i < ntasks; i++ {
		taskChannel <- i
	}

	// WaitGroup to wait for all tasks to be completed
	var wg sync.WaitGroup
	wg.Add(ntasks)

	// Loop through all tasks
	for i := 0 ; i < ntasks; i++ {
		// Get the task number from the channel
		task := <-taskChannel
		// Start a goroutine to do the task
		go func (taskNumber int, phase jobPhase, nios int)  {
			// Get a worker from the register channel
			worker := <- mr.registerChannel
			// Create the arguments for the DoTask RPC
			args := &DoTaskArgs{
				JobName:       mr.jobName,
				File:          mr.files[taskNumber],
				Phase:         phase,
				TaskNumber:    taskNumber,
				NumOtherPhase: nios,
			}
			// Call the DoTask RPC
			ok := call(worker, "Worker.DoTask", args, new(ShutdownReply))
			
			// If the RPC was successful, put the worker back in the register channel
			// Else, put the task number back in the task channel
			if ok {
				go func() { 
					// Decrement the wait group
					wg.Done()
					mr.registerChannel <- worker 
				}()
			} else {
				taskChannel <- taskNumber
			}
		}(task, phase, nios)
	}

	// Wait for all tasks to be completed
	wg.Wait()

	// ---------------------------

	debug("Schedule: %v phase done\n", phase)
}


func getNextTask(jobName string, file []string, taskNumber int, phase jobPhase, nios int) (task DoTaskArgs){
	
	task = DoTaskArgs{
		JobName:       jobName,
		File:          file[taskNumber],
		Phase:         phase,
		TaskNumber:    taskNumber,
		NumOtherPhase: nios,
	}
	return task
}