package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// using this as an inspiration https://github.com/moby/moby/blob/master/daemon/graphdriver/copy/copy.go
// CopyDirectory
func CopyDirectory(srcDir, dest string) error {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		// TODO: test this without
		// var gid, uid int

		// if runtime.GOOS != "windows" {
		// 	//  REMOVING this as syscall is kind of depracated in favour of sys ==> platform specific low level abstraction
		// 	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		// 	if !ok {
		// 		return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		// 	}
		// 	uid = int(stat.Uid)
		// 	gid = int(stat.Gid)
		// }

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		// // On Windows, it always returns the syscall.EWINDOWS error, wrapped in *PathError.
		if runtime.GOOS != "windows" {
			// ensure file ownership for current process owner
			if err := os.Lchown(destPath, os.Geteuid(), os.Getegid()); err != nil {
				return err
			}
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
// Returns all unzipped files (abs path)
func Unzip(src, dest string) error {

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// iterate around the files in the zip
	for _, file := range r.File {

		// Determine the path for the current file
		filePath := filepath.Join(dest, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", filePath)
		}

		// if the file is a directory, create it
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)

			// move onto the next iteration in the loop
			continue
		}

		// Make the file
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// read the content of the current file so it can be set in the destination file
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		// open the current file
		rc, err := file.Open()
		if err != nil {
			return err
		}

		// copy the content of the current to the new one
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}

		// close the files
		outFile.Close()
		rc.Close()

	}

	return nil
}

// GetDefaultTempDir determines and creates a temporary directory for packages and source
// code to be downloaded to.
func GetDefaultTempDir() string {
	tmpPath, err := os.MkdirTemp("", "stackscli")
	if err != nil {
		log.Fatalf("Unable to create temporary directory")
	}
	return tmpPath
}

// GetDefaultWorkingDir returns the current directory as the default working directory
// for where projects will be created
func GetDefaultWorkingDir() string {

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to determine current directory")
	}

	return workingDir
}
