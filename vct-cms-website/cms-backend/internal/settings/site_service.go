package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ─── SITE CONFIG STRUCTS (map 1:1 với site.json) ───

type SiteGeneral struct {
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Edition     string `json:"edition"`
	Year        int    `json:"year"`
	Copyright   string `json:"copyright"`
	URL         string `json:"url"`
	Locale      string `json:"locale"`
	Icon        string `json:"icon"`
}

type SiteContact struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	WorkingHours string `json:"workingHours"`
}

type SiteSocial struct {
	Facebook string `json:"facebook"`
	YouTube  string `json:"youtube"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
	Zalo     string `json:"zalo"`
}

type SiteSEO struct {
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	Keywords        string `json:"keywords"`
	OGImage         string `json:"ogImage"`
	CanonicalURL    string `json:"canonicalUrl"`
}

// SiteConfig chứa toàn bộ cấu trúc site.json
type SiteConfig struct {
	Site       SiteGeneral            `json:"site"`
	Stats      map[string]interface{} `json:"stats"`
	Navigation []interface{}          `json:"navigation"`
	FooterLink []interface{}          `json:"footer_links"`
	Contact    SiteContact            `json:"contact"`
	Social     SiteSocial             `json:"social"`
	SEO        SiteSEO                `json:"seo"`
	TechStack  []interface{}          `json:"tech_stack"`
	Blog       []interface{}          `json:"blog"`
}

type SiteService struct {
	RepoPath string
	FilePath string
}

func NewSiteService(repoPath string) *SiteService {
	return &SiteService{
		RepoPath: repoPath,
		FilePath: filepath.Join(repoPath, "data", "site.json"),
	}
}

// GetSiteConfig đọc toàn bộ site.json
func (s *SiteService) GetSiteConfig() (*SiteConfig, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc site.json: %w", err)
	}
	var config SiteConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("lỗi parse site.json: %w", err)
	}
	return &config, nil
}

// SaveSiteConfig ghi toàn bộ SiteConfig vào site.json
func (s *SiteService) SaveSiteConfig(config *SiteConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0644)
}

// SyncToGit commit & push thay đổi settings
func (s *SiteService) SyncToGit(commitMsg string) error {
	repo, err := git.PlainOpen(s.RepoPath)
	if err != nil {
		return fmt.Errorf("không thể mở repo: %w", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("không thể lấy worktree: %w", err)
	}
	_ = wt.Pull(&git.PullOptions{RemoteName: "origin"})

	if _, err = wt.Add("."); err != nil {
		return fmt.Errorf("lỗi git add: %w", err)
	}
	if _, err = wt.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "VCT CMS Bot",
			Email: "cms-bot@vct-platform.com",
			When:  time.Now(),
		},
	}); err != nil {
		return fmt.Errorf("lỗi git commit: %w", err)
	}
	if err = repo.Push(&git.PushOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("lỗi git push: %w", err)
	}
	return nil
}
