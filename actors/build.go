package actors

import (
	"astroterm/astro"
	"astroterm/db"
	"astroterm/project"

	"github.com/Licoy/stail"
)

type BuildActor struct {
	DB  *db.Database
	Pid int

	project *project.Project
	si      stail.STailItem

	blogs *broker[string]
}

func NewBuildActor(project *project.Project) *BuildActor {
	return &BuildActor{
		DB:      db.NewDatabase(),
		project: project,

		blogs: newBroker[string](),
	}
}

func (b *BuildActor) Start() *BuildActor {
	go b.startup()
	return b
}

func (b *BuildActor) StartBuild() chan bool {
	done := make(chan bool)
	go func() {
		b.startBuild()
		done <- true
	}()
	return done
}

func (b *BuildActor) SubscribeToLogs() chan string {
	return b.blogs.Subscribe()
}

// Private
func (b *BuildActor) startup() {
	go b.blogs.Start()
}

func (b *BuildActor) startBuild() error {
	pid, logpth, err := astro.RunBuildAndPipeToLog(b.project.Dir)
	if err != nil {
		return err
	}
	b.Pid = pid
	b.tailLogFile(logpth)
	return nil
}

func (b *BuildActor) tailLogFile(logpth string) error {
	st, err := stail.New(stail.Options{})
	if err != nil {
		return err
	}
	si, err := st.Tail(logpth, 0, func(content string) {
		b.blogs.Publish(content)
	})
	b.si = si

	if err != nil {
		return err
	}
	b.si.Watch()
	return nil
}
