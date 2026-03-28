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

// TranslationEntry đại diện cho 1 key dịch, chứa cả 2 ngôn ngữ
type TranslationEntry struct {
	Key string `json:"key"`
	Vi  string `json:"vi"`
	En  string `json:"en"`
}

type I18nService struct {
	RepoPath string
	ViPath   string
	EnPath   string
}

func NewI18nService(repoPath string) *I18nService {
	return &I18nService{
		RepoPath: repoPath,
		ViPath:   filepath.Join(repoPath, "data", "lang", "vi.json"),
		EnPath:   filepath.Join(repoPath, "data", "lang", "en.json"),
	}
}

// readJSON đọc một file JSON flat (key-value string) từ đĩa
func readJSON(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// writeJSON ghi map[string]string ra file JSON có thụt lề
func writeJSON(path string, data map[string]string) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// GetMergedTranslations đọc cả 2 file, merge theo key trả về mảng
func (s *I18nService) GetMergedTranslations() ([]TranslationEntry, error) {
	viData, err := readJSON(s.ViPath)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc vi.json: %w", err)
	}
	enData, err := readJSON(s.EnPath)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc en.json: %w", err)
	}

	// Thu thập tất cả key duy nhất từ cả 2 file
	allKeys := make(map[string]bool)
	for k := range viData {
		allKeys[k] = true
	}
	for k := range enData {
		allKeys[k] = true
	}

	// Merge thành mảng TranslationEntry
	entries := make([]TranslationEntry, 0, len(allKeys))
	for key := range allKeys {
		entries = append(entries, TranslationEntry{
			Key: key,
			Vi:  viData[key],
			En:  enData[key],
		})
	}

	return entries, nil
}

// SaveTranslations tách mảng entry ngược lại thành 2 file vi.json & en.json
func (s *I18nService) SaveTranslations(entries []TranslationEntry) error {
	viData := make(map[string]string, len(entries))
	enData := make(map[string]string, len(entries))

	for _, e := range entries {
		if e.Vi != "" {
			viData[e.Key] = e.Vi
		}
		if e.En != "" {
			enData[e.Key] = e.En
		}
	}

	if err := writeJSON(s.ViPath, viData); err != nil {
		return fmt.Errorf("lỗi ghi vi.json: %w", err)
	}
	if err := writeJSON(s.EnPath, enData); err != nil {
		return fmt.Errorf("lỗi ghi en.json: %w", err)
	}

	return nil
}

// SyncToGit commit & push thay đổi i18n lên GitHub
func (s *I18nService) SyncToGit(commitMsg string) error {
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
