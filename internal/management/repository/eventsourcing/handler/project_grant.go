package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type ProjectGrant struct {
	handler
	eventstore    eventstore.Eventstore
	projectEvents *proj_event.ProjectEventstore
	orgEvents     *org_event.OrgEventstore
}

const (
	grantedProjectTable = "management.project_grants"
)

func (p *ProjectGrant) ViewModel() string {
	return grantedProjectTable
}

func (p *ProjectGrant) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestProjectGrantSequence()
	if err != nil {
		return nil, err
	}
	return proj_event.ProjectQuery(sequence.CurrentSequence), nil
}

func (p *ProjectGrant) Reduce(event *models.Event) (err error) {
	grantedProject := new(view_model.ProjectGrantView)
	switch event.Type {
	case es_model.ProjectChanged:
		project, err := p.view.ProjectByID(event.AggregateID)
		if err != nil {
			return err
		}
		return p.updateExistingProjects(project)
	case es_model.ProjectGrantAdded:
		err = grantedProject.AppendEvent(event)
		if err != nil {
			return err
		}
		project, err := p.getProject(grantedProject.ProjectID)
		if err != nil {
			return err
		}
		grantedProject.Name = project.Name

		org, err := p.orgEvents.OrgByID(context.TODO(), org_model.NewOrg(grantedProject.OrgID))
		if err != nil {
			return err
		}
		resourceOwner, err := p.orgEvents.OrgByID(context.TODO(), org_model.NewOrg(grantedProject.ResourceOwner))
		if err != nil {
			return err
		}
		p.fillOrgData(grantedProject, org, resourceOwner)
	case es_model.ProjectGrantChanged, es_model.ProjectGrantCascadeChanged:
		grant := new(view_model.ProjectGrant)
		err := grant.SetData(event)
		if err != nil {
			return err
		}
		grantedProject, err = p.view.ProjectGrantByID(grant.GrantID)
		if err != nil {
			return err
		}
		err = grantedProject.AppendEvent(event)
	case es_model.ProjectGrantRemoved:
		grant := new(view_model.ProjectGrant)
		err := grant.SetData(event)
		if err != nil {
			return err
		}
		return p.view.DeleteProjectGrant(grant.GrantID, event.Sequence)
	case es_model.ProjectRemoved:
		projectGrants, err := p.view.ProjectGrantsByProjectID(event.AggregateID)
		if err != nil {
			logging.LogWithFields("HANDL-h8jVH", "id", event.AggregateID).WithError(err).Warn("could not update existing projects")
		}
		for _, existing := range projectGrants {
			err := p.view.DeleteProjectGrant(existing.GrantID, event.Sequence)
			if err != nil {
				return err
			}
		}
	default:
		return p.view.ProcessedProjectGrantSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutProjectGrant(grantedProject)
}

func (p *ProjectGrant) fillOrgData(grantedProject *view_model.ProjectGrantView, org, resourceOwner *org_model.Org) {
	grantedProject.OrgName = org.Name
	grantedProject.ResourceOwnerName = resourceOwner.Name
}

func (p *ProjectGrant) getProject(projectID string) (*proj_model.Project, error) {
	return p.projectEvents.ProjectByID(context.Background(), projectID)
}

func (p *ProjectGrant) updateExistingProjects(project *view_model.ProjectView) error {
	projectGrants, err := p.view.ProjectGrantsByProjectID(project.ProjectID)
	if err != nil {
		logging.LogWithFields("SPOOL-los03", "id", project.ProjectID).WithError(err).Warn("could not update existing projects")
	}
	for _, existing := range projectGrants {
		existing.Name = project.Name
		err := p.view.PutProjectGrant(existing)
		if err != nil {
			logging.LogWithFields("SPOOL-sjwi3", "id", existing.ProjectID).WithError(err).Warn("could not update existing project")
			return err
		}
	}

	return p.view.ProcessedProjectGrantSequence(project.Sequence)
}

func (p *ProjectGrant) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-is8wa", "id", event.AggregateID).WithError(err).Warn("something went wrong in granted projecthandler")
	return spooler.HandleError(event, err, p.view.GetLatestProjectGrantFailedEvent, p.view.ProcessedProjectGrantFailedEvent, p.view.ProcessedProjectGrantSequence, p.errorCountUntilSkip)
}
