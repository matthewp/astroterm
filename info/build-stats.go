package info

import (
	"os"
	"path/filepath"
	"strings"
)

type BuildStats struct {
	HasBuilt      bool
	NumberOfPages int
	PageStats     []*BuildPageStats
}

type BuildPageStats struct {
	Path     string
	HtmlSize int64
	// More stuff to follow, such as the amount of JS
}

func (stats *BuildStats) CollectStatsForStaticOutDir(outDir string) {
	stats.HasBuilt = stats.crawlStaticOutputDir(outDir, outDir)
}

func (stats *BuildStats) crawlStaticOutputDir(outDir string, subDir string) bool {
	files, err := os.ReadDir(subDir)
	if err != nil {
		return false
	}
	for _, file := range files {
		fullPath := filepath.Join(subDir, file.Name())
		if file.IsDir() {
			stats.crawlStaticOutputDir(outDir, fullPath)
			continue
		}
		switch true {
		case strings.HasSuffix(file.Name(), ".html"):
			stats.NumberOfPages += 1
			path := ""
			if relPath, err := filepath.Rel(outDir, subDir); err == nil {
				path = filepath.Join("/", relPath)
			}
			var htmlSize int64 = 0
			if fi, err := os.Stat(fullPath); err == nil {
				htmlSize = fi.Size()
			}
			pageStats := &BuildPageStats{
				Path:     path,
				HtmlSize: htmlSize,
			}
			stats.PageStats = append(stats.PageStats, pageStats)
			break
		// TODO some JS stuff
		default:
			break
		}
	}
	return true
}
