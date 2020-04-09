package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type ProjectCache struct {
	projectCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*ProjectCache, error) {
	projectCache, err := conf.Config.NewCache()
	logging.Log("EVENT-vDneN").OnError(err).Panic("unable to create project cache")

	return &ProjectCache{projectCache: projectCache}, nil
}

func (c *ProjectCache) getProject(ID string) (project *Project, sequence uint64) {
	project = &Project{ObjectRoot: models.ObjectRoot{ID: ID}}
	if err := c.projectCache.Get(ID, project); err == nil {
		sequence = project.Sequence
	} else {
		logging.Log("EVENT-4eTZh").WithError(err).Debug("error in getting cache")
	}
	return project, sequence
}

func (c *ProjectCache) cacheProject(project *Project) {
	err := c.projectCache.Set(project.ID, project)
	if err != nil {
		logging.Log("EVENT-ThnBb").WithError(err).Debug("error in setting project cache")
	}
}
