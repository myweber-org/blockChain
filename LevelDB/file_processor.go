package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers   int
	batchSize int
	mu        sync.Mutex
	stats     map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(i, fileChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	fp.printStats()
	return nil
}

func (fp *FileProcessor) worker(id int, files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		start := time.Now()
		err := fp.processSingleFile(file)
		duration := time.Since(start)

		fp.mu.Lock()
		if err != nil {
			fp.stats["errors"]++
			fmt.Printf("Worker %d: Failed to process %s: %v\n", id, file, err)
		} else {
			fp.stats["processed"]++
			fmt.Printf("Worker %d: Processed %s in %v\n", id, file, duration)
		}
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) processSingleFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	switch ext {
	case ".txt":
		return fp.processTextFile(file)
	case ".log":
		return fp.processLogFile(file)
	default:
		return fp.processGenericFile(file)
	}
}

func (fp *FileProcessor) processTextFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		_ = scanner.Text()
	}

	fp.mu.Lock()
	fp.stats["lines"] += lineCount
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) processLogFile(file *os.File) error {
	_, err := io.Copy(io.Discard, file)
	return err
}

func (fp *FileProcessor) processGenericFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	fp.mu.Lock()
	fp.stats["bytes"] += int(info.Size())
	fp.mu.Unlock()

	return nil
}

func (fp *FileProcessor) printStats() {
	fmt.Println("\nProcessing Statistics:")
	fmt.Println("=====================")
	for key, value := range fp.stats {
		fmt.Printf("%s: %d\n", key, value)
	}
}

func main() {
	processor := NewFileProcessor(4, 10)

	files := []string{
		"data/file1.txt",
		"data/file2.log",
		"data/file3.csv",
		"data/file4.txt",
	}

	if err := processor.ProcessFiles(files); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}
}package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileStats struct {
	Path        string
	Size        int64
	LineCount   int
	WordCount   int
	ProcessTime time.Duration
	Error       error
}

func ProcessFile(path string) (FileStats, error) {
	start := time.Now()
	stats := FileStats{Path: path}

	file, err := os.Open(path)
	if err != nil {
		stats.Error = err
		return stats, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		stats.Error = err
		return stats, err
	}
	stats.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		lineCount++
		words := splitWords(scanner.Text())
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		stats.Error = err
		return stats, err
	}

	stats.LineCount = lineCount
	stats.WordCount = wordCount
	stats.ProcessTime = time.Since(start)

	return stats, nil
}

func splitWords(line string) []string {
	var words []string
	var currentWord []rune

	for _, r := range line {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			currentWord = append(currentWord, r)
		} else {
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = nil
			}
		}
	}

	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}

	return words
}

func ProcessFilesConcurrently(paths []string, maxWorkers int) ([]FileStats, error) {
	if len(paths) == 0 {
		return nil, errors.New("no files to process")
	}

	if maxWorkers <= 0 {
		maxWorkers = 4
	}

	jobs := make(chan string, len(paths))
	results := make(chan FileStats, len(paths))
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				stats, _ := ProcessFile(path)
				results <- stats
			}
		}()
	}

	for _, path := range paths {
		jobs <- path
	}
	close(jobs)

	wg.Wait()
	close(results)

	var allStats []FileStats
	for stats := range results {
		allStats = append(allStats, stats)
	}

	return allStats, nil
}

func FindFilesByPattern(rootDir, pattern string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

func WriteStatsToFile(stats []FileStats, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	header := "Path,Size (bytes),Lines,Words,Process Time (ms),Error\n"
	if _, err := writer.WriteString(header); err != nil {
		return err
	}

	for _, s := range stats {
		errorMsg := ""
		if s.Error != nil {
			errorMsg = s.Error.Error()
		}

		line := fmt.Sprintf("%s,%d,%d,%d,%v,%s\n",
			s.Path,
			s.Size,
			s.LineCount,
			s.WordCount,
			s.ProcessTime.Milliseconds(),
			errorMsg)

		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	files, err := FindFilesByPattern(".", "*.go")
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No Go files found")
		return
	}

	fmt.Printf("Processing %d files...\n", len(files))

	stats, err := ProcessFilesConcurrently(files, 4)
	if err != nil {
		fmt.Printf("Error processing files: %v\n", err)
		return
	}

	totalLines := 0
	totalWords := 0
	totalSize := int64(0)
	errors := 0

	for _, s := range stats {
		if s.Error != nil {
			errors++
			fmt.Printf("Error processing %s: %v\n", s.Path, s.Error)
		} else {
			totalLines += s.LineCount
			totalWords += s.WordCount
			totalSize += s.Size
			fmt.Printf("%s: %d lines, %d words, %d bytes, processed in %v\n",
				filepath.Base(s.Path),
				s.LineCount,
				s.WordCount,
				s.Size,
				s.ProcessTime)
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total files: %d\n", len(files))
	fmt.Printf("Successfully processed: %d\n", len(files)-errors)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Total words: %d\n", totalWords)
	fmt.Printf("Total size: %d bytes\n", totalSize)

	if err := WriteStatsToFile(stats, "file_stats.csv"); err != nil {
		fmt.Printf("Error writing stats: %v\n", err)
	} else {
		fmt.Println("Statistics written to file_stats.csv")
	}
}