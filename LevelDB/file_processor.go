
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	fileList []string
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		fileList: make([]string, 0),
	}
}

func (fp *FileProcessor) ScanDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fp.mu.Lock()
			fp.fileList = append(fp.fileList, path)
			fp.mu.Unlock()
		}
		return nil
	})
}

func (fp *FileProcessor) ProcessFiles() {
	var wg sync.WaitGroup
	results := make(chan string, len(fp.fileList))

	for _, file := range fp.fileList {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			result := fp.processSingleFile(f)
			results <- result
		}(file)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println(result)
	}
}

func (fp *FileProcessor) processSingleFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("Error reading file %s: %v", filePath, err)
	}

	return fmt.Sprintf("File: %s, Lines: %d", filepath.Base(filePath), lineCount)
}

func (fp *FileProcessor) GetFileCount() int {
	return len(fp.fileList)
}

func main() {
	processor := NewFileProcessor()
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		return
	}

	dirPath := os.Args[1]
	err := processor.ScanDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	fmt.Printf("Found %d files\n", processor.GetFileCount())
	processor.ProcessFiles()
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  map[string]int
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", path, err)
			return
		}
		defer file.Close()
		
		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", path, err)
			return
		}
		
		fp.mu.Lock()
		fp.results[path] = lineCount
		fp.mu.Unlock()
		
		fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	}()
	
	return nil
}

func (fp *FileProcessor) Wait() {
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]int {
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		os.Exit(1)
	}
	
	processor := NewFileProcessor()
	
	for _, filePath := range os.Args[1:] {
		if err := processor.ProcessFile(filePath); err != nil {
			fmt.Printf("Failed to process %s: %v\n", filePath, err)
		}
	}
	
	processor.Wait()
	
	fmt.Println("\nProcessing complete. Results:")
	for file, lines := range processor.GetResults() {
		fmt.Printf("%s: %d lines\n", filepath.Base(file), lines)
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

type FileProcessor struct {
	BatchSize   int
	Workers     int
	mu          sync.Mutex
	processed   int
	errors      []error
}

func NewFileProcessor(batchSize, workers int) *FileProcessor {
	if batchSize <= 0 {
		batchSize = 100
	}
	if workers <= 0 {
		workers = 4
	}
	return &FileProcessor{
		BatchSize: batchSize,
		Workers:   workers,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, len(paths))
	resultChan := make(chan error, len(paths))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan, resultChan, processor)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	close(resultChan)

	for err := range resultChan {
		if err != nil {
			fp.mu.Lock()
			fp.errors = append(fp.errors, err)
			fp.mu.Unlock()
		}
	}

	if len(fp.errors) > 0 {
		return fmt.Errorf("encountered %d errors during processing", len(fp.errors))
	}
	return nil
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan string, results chan<- error, processor func(string) error) {
	defer wg.Done()
	for file := range files {
		err := processor(file)
		results <- err
		fp.mu.Lock()
		fp.processed++
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

func (fp *FileProcessor) CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (fp *FileProcessor) FindFilesByPattern(dir, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func (fp *FileProcessor) Stats() (processed int, errorCount int) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, len(fp.errors)
}

func main() {
	processor := NewFileProcessor(50, 8)
	
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	
	err := processor.ProcessFiles(files, func(path string) error {
		lines, err := processor.CountLines(path)
		if err != nil {
			return err
		}
		fmt.Printf("File %s has %d lines\n", path, lines)
		return nil
	})
	
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
	}
	
	processed, errors := processor.Stats()
	fmt.Printf("Processed %d files with %d errors\n", processed, errors)
	
	time.Sleep(100 * time.Millisecond)
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
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		stats.Error = err
		return stats, err
	}
	stats.LineCount = lineCount
	stats.ProcessTime = time.Since(start)

	return stats, nil
}

func ProcessFilesConcurrently(paths []string, maxWorkers int) ([]FileStats, error) {
	if maxWorkers <= 0 {
		return nil, errors.New("maxWorkers must be positive")
	}

	var wg sync.WaitGroup
	workChan := make(chan string, len(paths))
	resultsChan := make(chan FileStats, len(paths))
	stats := make([]FileStats, 0, len(paths))

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range workChan {
				fileStats, _ := ProcessFile(path)
				resultsChan <- fileStats
			}
		}()
	}

	for _, path := range paths {
		workChan <- path
	}
	close(workChan)

	wg.Wait()
	close(resultsChan)

	for stat := range resultsChan {
		stats = append(stats, stat)
	}

	return stats, nil
}

func FindFilesByExtension(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func PrintStats(stats []FileStats) {
	totalSize := int64(0)
	totalLines := 0
	var totalTime time.Duration

	fmt.Println("File Processing Results:")
	fmt.Println("=======================")
	for _, s := range stats {
		fmt.Printf("Path: %s\n", s.Path)
		fmt.Printf("  Size: %d bytes\n", s.Size)
		fmt.Printf("  Lines: %d\n", s.LineCount)
		fmt.Printf("  Process Time: %v\n", s.ProcessTime)
		if s.Error != nil {
			fmt.Printf("  Error: %v\n", s.Error)
		}
		fmt.Println()

		totalSize += s.Size
		totalLines += s.LineCount
		totalTime += s.ProcessTime
	}

	fmt.Printf("Total Files: %d\n", len(stats))
	fmt.Printf("Total Size: %d bytes\n", totalSize)
	fmt.Printf("Total Lines: %d\n", totalLines)
	fmt.Printf("Total Process Time: %v\n", totalTime)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	files, err := FindFilesByExtension(dir, ".txt")
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No .txt files found in directory")
		return
	}

	fmt.Printf("Found %d .txt files\n", len(files))
	stats, err := ProcessFilesConcurrently(files, 4)
	if err != nil {
		fmt.Printf("Error processing files: %v\n", err)
		os.Exit(1)
	}

	PrintStats(stats)
}