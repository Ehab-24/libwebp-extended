package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const MAX_WORKERS = 20

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("expected %d arguments, found %d", 3, len(os.Args))
	}

	srcDir := os.Args[1]
	outDir := os.Args[2]

	filenames := listFiles(srcDir)
	var wg sync.WaitGroup
	guard := make(chan int, MAX_WORKERS)
	defer close(guard)

	for _, fn := range filenames {
		if !isImageName(fn) {
			continue
		}

		srcFile := srcDir + "/" + fn
		outFile := outDir + "/" + fn[0:strings.LastIndex(fn, ".")] + ".webp"

		guard <- 1
		wg.Add(1)
		go func() {
			execCWebp("60", srcFile, outFile)
			<-guard
			wg.Done()
		}()
	}

	wg.Wait()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isImageName(str string) bool {
	return strings.HasSuffix(str, ".png") || strings.HasSuffix(str, ".jpg") || strings.HasSuffix(str, ".jpeg")
}

func execCWebp(q string, srcFile string, outFile string) {
	cmd := exec.Command("cwebp", "-q", q, srcFile, "-o", outFile)

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	err := cmd.Run()
	check(err)
}

func listFiles(dir string) []string {
	cmd := exec.Command("ls", dir)
	cmd.Stderr = log.Writer()

	out, err := cmd.Output()
	check(err)

	filenames := string(out)
	return strings.Split(filenames, "\n")
}
