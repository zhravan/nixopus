package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/vmihailenco/taskq/v3"
)

const (
	queueCustomDomain        = "custom-domain"
	taskRegisterCustomDomain = "task_register_custom_domain"
	taskRemoveCustomDomain   = "task_remove_custom_domain"
)

type CustomDomainPayload struct {
	DomainID  string `json:"domain_id"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	GuestIP   string `json:"guest_ip"`
}

type RemoveCustomDomainPayload struct {
	DomainID string `json:"domain_id"`
	Domain   string `json:"domain"`
}

var (
	onceCustomDomainQueues   sync.Once
	CustomDomainQueue        taskq.Queue
	TaskRegisterCustomDomain *taskq.Task
	TaskRemoveCustomDomain   *taskq.Task
)

func SetupCustomDomainQueue() {
	onceCustomDomainQueues.Do(func() {
		CustomDomainQueue = registerProducerQueue(&taskq.QueueOptions{
			Name: queueCustomDomain,
		})

		TaskRegisterCustomDomain = taskq.RegisterTask(&taskq.TaskOptions{
			Name:       taskRegisterCustomDomain,
			RetryLimit: 1,
			Handler: func(ctx context.Context, payload CustomDomainPayload) error {
				return nil
			},
		})

		TaskRemoveCustomDomain = taskq.RegisterTask(&taskq.TaskOptions{
			Name:       taskRemoveCustomDomain,
			RetryLimit: 1,
			Handler: func(ctx context.Context, payload RemoveCustomDomainPayload) error {
				return nil
			},
		})
	})
}

func EnqueueRegisterCustomDomain(ctx context.Context, payload CustomDomainPayload) error {
	if CustomDomainQueue == nil {
		return fmt.Errorf("custom domain queue not initialized - call SetupCustomDomainQueue first")
	}
	return CustomDomainQueue.Add(TaskRegisterCustomDomain.WithArgs(ctx, payload))
}

func EnqueueRemoveCustomDomain(ctx context.Context, payload RemoveCustomDomainPayload) error {
	if CustomDomainQueue == nil {
		return fmt.Errorf("custom domain queue not initialized - call SetupCustomDomainQueue first")
	}
	return CustomDomainQueue.Add(TaskRemoveCustomDomain.WithArgs(ctx, payload))
}
