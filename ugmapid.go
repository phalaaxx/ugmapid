package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

/* change uid and gid for a single path element */
func walkFn(base, offset int, verbose bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		stat := info.Sys().(*syscall.Stat_t)
		/* get target uid and gid values */
		tgtUid, tgtGid := int(stat.Uid), int(stat.Gid)
		/* check ranges */
		if tgtUid < offset || tgtUid > offset+65536 {
			tgtUid = offset + (tgtUid - base)
		}
		if tgtGid < offset || tgtGid > offset+65536 {
			tgtGid = offset + (tgtGid - base)
		}
		/* change file/directory owner */
		if tgtUid != int(stat.Uid) || tgtGid != int(stat.Gid) {
			if err = os.Lchown(path, tgtUid, tgtGid); err != nil {
				log.Println("  ERROR: ", path)
				return err
			}
			/* show some progress */
			if verbose {
				log.Printf("%s : uid %d -> %d, gid %d -> %d\n",
					path, stat.Uid, tgtUid, stat.Gid, tgtGid)
			}
		}
		return nil
	}
}

/* main program */
func main() {
	/* initialize arguments */
	var (
		directory string
		offset    int
		verbose   bool
	)
	flag.StringVar(
		&directory,
		"directory",
		"",
		"Filesystem directory to traverse")
	flag.IntVar(
		&offset,
		"offset",
		100000,
		"uid/gid map offset")
	flag.BoolVar(
		&verbose,
		"verbose",
		false,
		"Display verbose output")
	flag.Parse()
	/* check proper directory value */
	if len(directory) == 0 {
		log.Fatal("error: directory argument is required")
	}
	st, err := os.Stat(directory)
	if err != nil {
		log.Fatal("error:", err)
	}
	if !st.Mode().IsDir() {
		log.Fatal("error: not a directory:", directory)
	}
	/* check for sane offset value */
	if offset < 0 {
		log.Fatal("error: offset must be a positive integer value")
	}
	/* calculate base id */
	base := int(st.Sys().(*syscall.Stat_t).Uid)
	log.Println("calculated base", base)
	/* walk files and update uid and gid */
	if err = filepath.Walk(directory, walkFn(base, offset, verbose)); err != nil {
		log.Fatal("Error: ", err)
	}
}
