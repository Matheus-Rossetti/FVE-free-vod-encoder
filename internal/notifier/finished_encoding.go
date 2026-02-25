package notifier

import "fmt"

func (n *Notifier) FinishedEncoding() {
	fmt.Println("Notifying everyone that an encoding job has finished...")

	// delete the downloaded video
}
