package git

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}
	head, err := getHeadName(repo)
	if err != nil {
		return nil, err
	}
	endpoint, err := getRemoteEndpoint(repo)
	if err != nil {
		return nil, err
	}
	spaceKey, domain := extractSpaceKeyAndDomain(endpoint.Host)
	projectKey, repoName := extractProjectKeyAndRepoName(endpoint.Path)

	return &Repository{
		repo:       repo,
		head:       head,
		domain:     domain,
		spaceKey:   spaceKey,
		projectKey: projectKey,
		repoName:   repoName,
	}, nil
}

func getHeadName(repo *git.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().String(), nil
}

func getRemoteEndpoint(repo *git.Repository) (*transport.Endpoint, error) {
	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}
	cfg := remote.Config()
	if len(cfg.URLs) == 0 {
		return nil, errors.New("could not find remote URL")
	}
	urlStr := cfg.URLs[0]
	endpoint, err := transport.NewEndpoint(urlStr)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func extractSpaceKeyAndDomain(host string) (spaceKey, domain string) {
	delimitedHost := strings.Split(host, ".")
	spaceKey = delimitedHost[0]
	domain = strings.Join(delimitedHost[len(delimitedHost)-2:], ".")
	return
}

func extractProjectKeyAndRepoName(path string) (projectKey, repoName string) {
	epPath := strings.TrimPrefix(path, "/git")
	delimitedPath := strings.Split(epPath, "/")
	projectKey = delimitedPath[1]
	repoName = strings.TrimSuffix(delimitedPath[2], ".git")
	return
}

const (
	refPrefix     = "refs/"
	refPullPrefix = refPrefix + "pull/"
	refPullSuffix = "/head"
)

type Repository struct {
	repo       *git.Repository
	head       string // current branch
	domain     string
	spaceKey   string
	projectKey string
	repoName   string
}

func (r *Repository) BaseURL() string {
	return fmt.Sprintf("https://%s.%s/", r.spaceKey, r.domain)
}

func (r *Repository) Space() string {
	return r.spaceKey
}

func (r *Repository) Project() string {
	return r.projectKey
}

func (r *Repository) Name() string {
	return r.repoName
}

func (r *Repository) PullRequestNumberOfCurrentBranch() (int, error) {

	ref2rev, err := r.lsRemote()
	if err != nil {
		return 0, err
	}

	headRev, ok := ref2rev[r.head]
	if !ok {
		return 0, errors.New("not found a current branch in remote")
	}

	var nums []int
	for ref, rev := range ref2rev {
		if !isPRRef(ref) {
			continue
		}
		if rev != headRev {
			continue
		}
		num, err := extractPRNum(ref)
		if err != nil {
			continue
		}
		nums = append(nums, num)
	}

	if len(nums) == 0 {
		return 0, errors.New("not found a pull request related to current branch")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(nums)))

	return nums[0], nil
}

func (r *Repository) lsRemote() (map[string]string, error) {
	cmd := exec.Command("git", "ls-remote", "-q")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return ref2rev(out), nil
}

func (r *Repository) CurrentDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-prefix")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func isPRRef(ref string) bool {
	return strings.HasPrefix(ref, refPullPrefix) && strings.HasSuffix(ref, refPullSuffix)
}

func extractPRNum(ref string) (int, error) {
	num := strings.TrimPrefix(ref, refPullPrefix)
	num = strings.TrimSuffix(num, refPullSuffix)
	return strconv.Atoi(num)
}

func ref2rev(b []byte) map[string]string {
	refToHash := make(map[string]string)
	remotes := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	for _, v := range remotes {
		delimited := strings.Split(v, "\t")
		hash := delimited[0]
		ref := delimited[1]
		refToHash[ref] = hash
	}
	return refToHash
}
