package core

type Notification interface {

	// Notify everyone that an encoding process has finished, this includes the downloadSlots and the cleanupWorker (cleans an already encoded video)
	FinishedEncoding()
}
