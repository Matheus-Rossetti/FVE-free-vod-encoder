package notifier

import "fmt"

func (n *Notifier) FinishedEncoding() {
	fmt.Println("Notifying everyone...")
}
