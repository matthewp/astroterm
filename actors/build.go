package actors

import (
	"astroterm/astro"
	"astroterm/db"
	"astroterm/info"
	"astroterm/project"

	"github.com/Licoy/stail"
)

type BuildActor struct {
	DB  *db.Database
	Pid int

	project *project.Project
	si      stail.STailItem
	config  *info.ConfigInfo
	stats   *info.BuildStats

	blogs  *broker[string]
	bstats *broker[*info.BuildStats]
}

func NewBuildActor(project *project.Project) *BuildActor {
	return &BuildActor{
		DB:      db.NewDatabase(),
		project: project,

		blogs:  newBroker[string](),
		bstats: newBroker[*info.BuildStats](),
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

func (b *BuildActor) SubscribeToStats() chan *info.BuildStats {
	ch := b.bstats.Subscribe()

	// If there are already stats, funnel the current into the channel
	if b.stats != nil {
		go func(stats *info.BuildStats) {
			ch <- stats
		}(b.stats)
	}

	return ch
}

func RunBuild() chan bool {
	ch := make(chan bool)

	go func() {

	}()

	return ch
}

// Private
func (b *BuildActor) startup() {
	go b.blogs.Start()
	//go b.figureOutDistSituation()
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

func (b *BuildActor) figureOutDistSituation() error {
	err := b.loadConfig()
	if err != nil {
		return err
	}
	b.stats = &info.BuildStats{
		NumberOfPages: 0,
	}
	if b.config.Output == "server" {

	} else {
		// static
		b.stats.CollectStatsForStaticOutDir(b.config.OutDir)
	}
	return nil
}

func (b *BuildActor) loadConfig() error {
	config, err := info.OpenConfig(b.project.Dir)
	b.config = config

	if err != nil {
		// TODO publish error maybe
		return err
	}

	return err
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
