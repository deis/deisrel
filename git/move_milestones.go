package git

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"
)

func MoveMilestone(
	ghClient *github.Client,
	repo string,
	oldMilestoneName string,
	newMilestoneName string,
	includeClosed bool,
) error {
	is := ghClient.Issues
	milestones, _, err := is.ListMilestones("deis", repo, &github.MilestoneListOptions{})
	if err != nil {
		return err
	}
	oldMilestone, err := getMilestoneFromMilestoneList(milestones, oldMilestoneName)
	if err != nil {
		return err
	}
	newMilestone, err := getMilestoneFromMilestoneList(milestones, newMilestoneName)
	if err != nil {
		return err
	}
	issueState := "open"
	if includeClosed {
		issueState = "all"
	}
	// This list will ALSO include PRs:
	oldMilestoneIssues, _, err := is.ListByRepo("deis", repo, &github.IssueListByRepoOptions{
		Milestone: fmt.Sprintf("%d", *oldMilestone.Number),
		State:     issueState,
		ListOptions: github.ListOptions{
			PerPage: 10000,
		},
	})
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	done := make(chan bool)
	errCh := make(chan error)
	defer close(errCh)
	for _, issue := range oldMilestoneIssues {
		wg.Add(1)
		go func(issue github.Issue) {
			defer wg.Done()
			ir := &github.IssueRequest{
				Milestone: newMilestone.Number,
			}
			if _, _, err := is.Edit("deis", repo, *issue.Number, ir); err != nil {
				errCh <- err
			}
		}(issue)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	errs := []error{}
	for {
		select {
		case <-done:
			if len(errs) > 0 {
				var errStr string
				for _, err := range errs {
					errStr = fmt.Sprintf("%s%s\n", errStr, err)
				}
				return fmt.Errorf(errStr)
			}
			return nil
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
}

func getMilestoneFromMilestoneList(milestones []github.Milestone, milestoneName string) (*github.Milestone, error) {
	for _, milestone := range milestones {
		if *milestone.Title == milestoneName {
			return &milestone, nil
		}
	}
	return nil, newErrMilestoneNotFound(milestoneName)
}

type errMilestoneNotFound struct {
	milestone string
}

func newErrMilestoneNotFound(milestone string) errMilestoneNotFound {
	return errMilestoneNotFound{milestone: milestone}
}

func (e errMilestoneNotFound) Error() string {
	return fmt.Sprintf("milestone %s not found", e.milestone)
}
