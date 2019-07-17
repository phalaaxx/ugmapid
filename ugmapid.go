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
		tgt_uid, tgt_gid := int(stat.Uid), int(stat.Gid)
		/* check ranges */
		if tgt_uid < offset || tgt_uid > offset+65536 {
			tgt_uid = offset + (tgt_uid - base)
		}
		if tgt_gid < offset || tgt_gid > offset+65536 {
			tgt_gid = offset + (tgt_gid - base)
		}
		/* change file/directory owner */
		if tgt_uid != int(stat.Uid) || tgt_gid != int(stat.Gid) {
			if err = os.Lchown(path, tgt_uid, tgt_gid); err != nil {
				log.Println("  ERROR: ", path)
				return err
			}
			/* show some progress */
			if verbose {
				log.Printf("%s : uid %d -> %d, gid %d -> %d\n",
					path, stat.Uid, tgt_uid, stat.Gid, tgt_gid)
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
