package managed

import (
	"context"
	"time"

	cronlib "github.com/hjwalt/runway/cron"
	"github.com/hjwalt/runway/inverse"
)

const (
	QualifierCron = "cron"
)

func AddCron(ic inverse.Container) {
	AddService(ic, &cron{
		cron: cronlib.New(
			cronlib.WithLocation(time.UTC),
			cronlib.WithSeconds(),
		),
	})
}

func AddCronJob(ic inverse.Container, cronjob cronlib.Job) {
	inverse.GenericAddVal(ic, QualifierCron, cronjob)
}

func AddCronJobResolver(ic inverse.Container, cronjob inverse.Injector[cronlib.Job]) {
	inverse.GenericAdd(ic, QualifierCron, cronjob)
}

type cron struct {
	lifecycle Lifecycle
	cron      *cronlib.Cron
}

func (r *cron) Name() string {
	return "cron"
}

func (r *cron) Register(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *cron) Resolve(ctx context.Context, ic inverse.Container) error {
	lifecycle, lifecycleErr := GetLifecycle(ic, ctx)
	if lifecycleErr != nil {
		return lifecycleErr
	}
	r.lifecycle = lifecycle

	runnables, runnableErr := inverse.GenericGetAll[cronlib.Job](ic, ctx, QualifierCron)
	if runnableErr != nil {
		return runnableErr
	}

	for _, runnable := range runnables {
		_, entryerr := r.cron.AddJob(runnable)
		if entryerr != nil {
			return entryerr
		}
	}

	return nil
}

func (r *cron) Clean() error {
	return nil
}

func (r *cron) Start() error {
	r.cron.Start()

	return nil
}

func (r *cron) Stop() error {
	r.cron.Stop()
	return nil
}
