package main

func WatchDog(pendingLinks chan string, pendingJobCount chan int) {
	count := 0

	for val := range pendingJobCount {
		count += val

		//log.Printf("Pending Job Tally: %d", count)

		if count == 0 {
			close(pendingLinks)
			close(pendingJobCount)
		}
	}
}
