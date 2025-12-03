package git_registry

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type AuthOptions struct {
	Username        string
	Token           string
	AuthorName      string
	GitStoreRepoUrl string
	PathRepo        string
}

type GitRegistry struct {
	options    AuthOptions
	workTree   *git.Worktree
	repository *git.Repository
	sync.Mutex
}

func cloneRepo(options AuthOptions) (*git.Repository, error) {
	return git.PlainClone(fmt.Sprintf("./%s", options.PathRepo), false, &git.CloneOptions{
		URL: options.GitStoreRepoUrl,
		Auth: &http.BasicAuth{
			Username: options.Username,
			Password: options.Token,
		},
	})
}

func openRepo(options AuthOptions) (*git.Repository, error) {
	return git.PlainOpen(options.PathRepo)
}
func NewGitRegistry(options AuthOptions) (*GitRegistry, error) {
	var err error
	var repo *git.Repository
	var workTree *git.Worktree
	if _, err = os.Stat(options.PathRepo); os.IsNotExist(err) {
		repo, err = cloneRepo(options)
		if err != nil {
			return nil, err
		}
	} else {
		repo, err = openRepo(options)
		if err != nil {
			return nil, err
		}
	}
	if repo == nil {
		return nil, fmt.Errorf("repo is nil")
	}
	workTree, err = repo.Worktree()
	if err != nil {
		return nil, err
	}
	return &GitRegistry{
		options:    options,
		workTree:   workTree,
		repository: repo,
	}, nil
}

func (gr *GitRegistry) checkoutBranch(branch string, isCreate bool) error {
	gr.Lock()
	defer gr.Unlock()
	err := gr.workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: isCreate,
	})
	return err
}

func (gr *GitRegistry) CreateBranch(branch string) error {
	return gr.checkoutBranch(branch, true)
}

func (gr *GitRegistry) CheckoutBranch(branch string) error {
	return gr.checkoutBranch(branch, false)
}

func (gr *GitRegistry) Pull() error {
	gr.Lock()
	defer gr.Unlock()
	return gr.workTree.Pull(&git.PullOptions{})
}

func (gr *GitRegistry) AddChanges(folder ...string) error {
	gr.Lock()
	defer gr.Unlock()
	for _, file := range folder {
		err := gr.workTree.AddWithOptions(&git.AddOptions{
			Path:       file,
			SkipStatus: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (gr *GitRegistry) CommitChanges(message string) error {
	gr.Lock()
	defer gr.Unlock()
	_, err := gr.workTree.Commit(message, &git.CommitOptions{})
	return err
}

func (gr *GitRegistry) Push() error {
	gr.Lock()
	defer gr.Unlock()
	err := gr.repository.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: gr.options.Username,
			Password: gr.options.Token,
		},
	})
	return err
}
