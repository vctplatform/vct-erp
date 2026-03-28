package blog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Post struct tương ứng với cấu trúc JSON hiện tại của VCT Platform
type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Summary   string   `json:"summary"`
	Thumbnail string   `json:"thumbnail"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	Link      string   `json:"link"`
}

type BlogService struct {
	RepoPath  string
	PostsPath string
	ImagesDir string
	mu        sync.Mutex // Serialize tất cả thao tác ghi để tránh race condition
}

func NewBlogService(repoPath string) *BlogService {
	return &BlogService{
		RepoPath:  repoPath,
		PostsPath: filepath.Join(repoPath, "data", "posts.json"),
		ImagesDir: filepath.Join(repoPath, "assets", "images"),
	}
}

// GetPosts đọc file posts.json và parse thành mảng Post
func (s *BlogService) GetPosts() ([]Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.getPostsUnsafe()
}

// getPostsUnsafe đọc mà KHÔNG lock — chỉ gọi khi đã hold lock
func (s *BlogService) getPostsUnsafe() ([]Post, error) {
	file, err := os.Open(s.PostsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Post{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var posts []Post
	bytes, _ := io.ReadAll(file)
	if err := json.Unmarshal(bytes, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// SavePosts lưu mảng Post vào posts.json (thread-safe)
func (s *BlogService) SavePosts(posts []Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.savePostsUnsafe(posts)
}

// savePostsUnsafe ghi mà KHÔNG lock — chỉ gọi khi đã hold lock
func (s *BlogService) savePostsUnsafe(posts []Post) error {
	jsonData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.PostsPath, jsonData, 0644)
}

// CreatePost thêm bài viết mới (thread-safe: lock xuyên suốt read-modify-write)
func (s *BlogService) CreatePost(newPost Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	posts, err := s.getPostsUnsafe()
	if err != nil {
		return err
	}
	posts = append([]Post{newPost}, posts...)
	return s.savePostsUnsafe(posts)
}

// UpdatePost cập nhật bài viết theo ID (thread-safe)
func (s *BlogService) UpdatePost(id string, updated Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	posts, err := s.getPostsUnsafe()
	if err != nil {
		return err
	}

	found := false
	for i, p := range posts {
		if p.ID == id {
			updated.ID = id
			posts[i] = updated
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("not_found")
	}

	return s.savePostsUnsafe(posts)
}

// SyncToGit thực hiện Pull -> Add -> Commit -> Push
func (s *BlogService) SyncToGit(commitMsg string) error {
	repo, err := git.PlainOpen(s.RepoPath)
	if err != nil {
		return fmt.Errorf("không thể mở repo: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("không thể lấy worktree: %w", err)
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		fmt.Printf("Cảnh báo khi pull: %v\n", err)
	}

	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("lỗi git add: %w", err)
	}

	_, err = worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "VCT CMS Bot",
			Email: "cms-bot@vct-platform.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("lỗi git commit: %w", err)
	}

	err = repo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("lỗi git push: %w", err)
	}

	return nil
}
