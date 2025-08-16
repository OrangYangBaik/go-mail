package controllers

import (
	"fmt"
	"go-mail/services"
	"os"
	"time"
)

type PriorimailController interface {
	ProcessEmails(filename string) error
}

type priorimailController struct {
	Svc services.PriorimailService
}

func NewPriorimailController(svc services.PriorimailService) PriorimailController {
	return &priorimailController{Svc: svc}
}

func (c *priorimailController) ProcessEmails(filename string) error {
	start := time.Now()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("starting batch\n"); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	errs := c.Svc.ProcessEmails()
	end := time.Now()

	if len(errs) > 0 {
		for _, e := range errs {
			if _, wErr := file.WriteString(e + "\n"); wErr != nil {
				return fmt.Errorf("failed to write error to log file: %w", wErr)
			}
		}
	}

	if _, err := file.WriteString(fmt.Sprintf("batch completed in %s\n", end.Sub(start))); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}
