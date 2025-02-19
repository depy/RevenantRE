package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/depy/RevenantRE/graphics"
)

func main() {
	path := flag.String("path", "", "Path to the directory to search for unique headers")
	flag.Parse()

	if *path == "" {
		fmt.Println("Path is required")
		return
	}

	uniqueHeaders := make(map[uint32]string)

	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() ||
			filepath.Ext(path) == ".bmp" ||
			filepath.Ext(path) == ".TN" ||
			filepath.Ext(path) == ".tn" ||
			filepath.Ext(path) == ".def" ||
			filepath.Ext(path) == ".wav" ||
			filepath.Ext(path) == ".mp3" ||
			filepath.Ext(path) == ".s" ||
			filepath.Base(path) == "quickload.dat" {
			return nil
		}

		//fmt.Println("Processing file: ", path)
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer file.Close()

		fmt.Println("Creating file resource for: ", path)
		fr, err := graphics.NewFileResource(file, true)
		if err != nil {
			fmt.Println(err)
			return err
		}

		//fmt.Println("Found bitmaps: ", len(fr.Bitmaps))

		for _, bm := range fr.Bitmaps {
			fval := bm.Header.Flags
			//fmt.Println("Extracted headers flags", fval, path)
			if _, ok := uniqueHeaders[fval]; !ok {
				//fmt.Println("Unique header: ", fval, " found in ", path)
				uniqueHeaders[fval] = path
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	for k, v := range uniqueHeaders {
		f := graphics.NewBitmapFlags(k)
		fmt.Println("----- ", v, " -----")
		graphics.PrintBitmapFlags(&f)
		fmt.Println("----------------------------------------")
	}
}
