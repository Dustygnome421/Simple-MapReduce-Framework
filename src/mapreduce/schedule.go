package mapreduce

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
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO:

	registerChan := make(chan string)

	// Create a channel for each worker
	go func ()  {
		for {
			worker := <-mr.registerChannel
			registerChan <- worker
		}
	}()

	taskChannel := make(chan int, ntasks)
	for i := 0; i < ntasks; i++ {
		taskChannel <- i
	}

	for len(taskChannel) > 0 {
		task := <-taskChannel
		go func (taskNumber int, phase jobPhase, nios int)  {
			worker := <- registerChan
			args := &DoTaskArgs{
				JobName:       mr.jobName,
				File:          mr.files[taskNumber],
				Phase:         phase,
				TaskNumber:    taskNumber,
				NumOtherPhase: nios,
			}
			ok := call(worker, "Worker.DoTask", args, new(struct{}))
			if ok {
				registerChan <- worker
				<- taskChannel
			} else {
				taskChannel <- taskNumber
			}
		}(task, phase, nios)
	}

	<- taskChannel


	// ---------------------------

	debug("Schedule: %v phase done\n", phase)
}
