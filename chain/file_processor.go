package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

type FileProcessor struct {
    inputDir  string
    outputDir string
    workers   int
}

func NewFileProcessor(input, output string, workers int) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
        workers:   workers,
    }
}

func (fp *FileProcessor) ProcessFiles() error {
    files, err := ioutil.ReadDir(fp.inputDir)
    if err != nil {
        return fmt.Errorf("failed to read input directory: %w", err)
    }

    var wg sync.WaitGroup
    fileChan := make(chan os.FileInfo, len(files))
    errChan := make(chan error, fp.workers)

    for i := 0; i < fp.workers; i++ {
        wg.Add(1)
        go fp.worker(&wg, fileChan, errChan)
    }

    for _, file := range files {
        if !file.IsDir() {
            fileChan <- file
        }
    }
    close(fileChan)

    wg.Wait()
    close(errChan)

    for err := range errChan {
        if err != nil {
            return err
        }
    }

    return nil
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, files <-chan os.FileInfo, errChan chan<- error) {
    defer wg.Done()

    for file := range files {
        if err := fp.processSingleFile(file.Name()); err != nil {
            errChan <- err
        }
    }
}

func (fp *FileProcessor) processSingleFile(filename string) error {
    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, filename)

    data, err := ioutil.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", filename, err)
    }

    processedData := transformData(data)

    if err := ioutil.WriteFile(outputPath, processedData, 0644); err != nil {
        return fmt.Errorf("failed to write file %s: %w", filename, err)
    }

    fmt.Printf("Processed: %s -> %s\n", inputPath, outputPath)
    return nil
}

func transformData(data []byte) []byte {
    result := make([]byte, len(data))
    for i, b := range data {
        result[i] = b ^ 0xFF
    }
    return result
}

func main() {
    processor := NewFileProcessor("./input", "./output", 4)
    if err := processor.ProcessFiles(); err != nil {
        fmt.Printf("Error processing files: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("File processing completed successfully")
}