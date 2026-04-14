package dispatcher

import (
	"os"
)

func (d *dispatcher) Start() {
	for job := range d.dispatchQueue {
		d.slog.Info("Received a job!", "dispatching contents from", job.DispatchJob.FromDir)

		for _, provider := range d.storageProviders {
			d.slog.Info("Dispatching", "to", provider.Name())

			err := provider.Dispatch(d.ctx, job.DispatchJob)
			if err != nil {
				d.slog.Error("dispatch failed", "provider", provider.Name(), "err", err)
				d.slog.Warn("Starting cleanup...", "at", provider.Name())
				provider.HandleError(job.DispatchJob)
			}

			d.slog.Info("Finished dispatching!", "to", provider.Name())
		}

		if !d.options.Dispatch.Local.Enabled { // local renames the dir
			err := os.RemoveAll(job.DispatchJob.FromDir) // so no need to remove
			if err != nil {
				d.slog.Error("error deleting a dir after dispatch", "err", err)
			}
		}

		d.slog.Info("Finished!")
	}

	// After queue closes
	d.slog.Info("Shutting down...")
}
