package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/Ensono/stacks-cli/internal/constants"
	"gopkg.in/yaml.v2"
)

// using this as an inspiration https://github.com/moby/moby/blob/master/daemon/graphdriver/copy/copy.go
// CopyDirectory
func CopyDirectory(srcDir, dest string) error {
	entries, err := os.ReadDir(srcDir)
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
		if GetPlatformOS() != "windows" {
			// ensure file ownership for current process owner
			if err := os.Lchown(destPath, os.Geteuid(), os.Getegid()); err != nil {
				return err
			}
		}

		fsInfo, _ := entry.Info()
		isSymlink := fsInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fsInfo.Mode()); err != nil {
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

	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) || err != nil {
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
func Unzip(src, dest string) (string, error) {

	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// iterate around the files in the zip
	for _, file := range r.File {

		// Determine the path for the current file
		filePath := filepath.Join(dest, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return "", fmt.Errorf("%s: illegal file path", filePath)
		}

		// if the file is a directory, create it
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)

			// move onto the next iteration in the loop
			continue
		}

		// Make the file
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return "", err
		}

		// read the content of the current file so it can be set in the destination file
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}

		// open the current file
		rc, err := file.Open()
		if err != nil {
			return "", err
		}

		// copy the content of the current to the new one
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return "", err
		}

		// Ensure that the file is read/writable
		// This is so that the stackscli.yml file can be read, and so that the files can be deleted
		if GetPlatformOS() != "windows" {
			err = os.Chmod(filePath, 0775)
			if err != nil {
				return "", err
			}
		}

		// close the files
		outFile.Close()
		rc.Close()

	}

	// get the directory in the dest which is the cloned dir
	files, _ := os.ReadDir(dest)
	tmpRepoDir := filepath.Join(dest, files[0].Name())

	return tmpRepoDir, nil
}

// GetDefaultTempDir determines the path to be used for the temporary directory
// It does not create it, but it will be set in the config for the CLI to create
// as and when it is required
func GetDefaultTempDir() string {

	var tmpPath string

	// determine the directory within the temp dir to use
	dir := fmt.Sprintf("stackscli%s", RandomString(10))

	// set the tempDir path based
	// On Windows this has to be defined slightly differently because the os.TempDir() returns
	// the env var $TMP, or $TEMP which uses the 8.3 naming convention. In normal circumstances
	// this is OK, however dotnet does not like using the short name
	// For example this dir `C:\Users\RussellSeymour\AppData\Local\Temp` will come out as
	// `C:\Users\RUSSEL~1\AppData\Local\Temp`. So detect Windows and set the tempdir using the $UserProfile
	// environment var and then append AppData\Local\Temp to it
	// Other OSes will use thje os.TempDir
	if GetPlatformOS() == "windows" {
		tmpPath = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Temp", dir)
	} else {
		tmpPath = filepath.Join(os.TempDir(), dir)
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

// GetUserHomeDir returns the currenmt users home directory
func GetUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine user home directory")
	}

	return homeDir
}

// GetStacksCLIDir returns the directory that should be used for storing all configuration
func GetStacksCLIDir() string {
	return path.Join(GetUserHomeDir(), constants.ConfigFileDir)
}

// GetDefaultCacheDir returns the directory that should be used for caching all downloads
// that the CLI makes
func GetDefaultCacheDir() string {

	var path string

	// get details about the user
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to determine current user")
	}

	// get the cacheDir path
	path = filepath.Join(usr.HomeDir, ".stackscli", "cache")

	return path
}

// IsEmpty states if the specified directory is empty or not
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func WriteYAMLToFile(object interface{}, path string, perm uint32) error {

	var err error

	// Ensure that the directory for the fie already exists
	basepath := filepath.Dir(path)
	err = CreateIfNotExists(basepath, fs.FileMode(perm))
	if err != nil {
		return err
	}

	// write out the input for the configuration return any errors
	data, err := yaml.Marshal(object)
	if err != nil {
		return err
	}

	// write out the data to the specified path
	err = os.WriteFile(path, data, 0666)

	return err
}

func GetFileContent(path string) ([]byte, error) {

	var err error
	var data []byte

	// read in the file, if it exists and return the results
	if Exists(path) {
		data, err = os.ReadFile(path)

		if err != nil {
			return data, err
		}

		return data, err
	} else {
		return data, fmt.Errorf("file does not exist: %s", path)
	}
}

// GetFileList returns a list of files that match the specified pattern
func GetFileList(pattern string, parent string) ([]string, error) {
	var err error
	var filelist []string

	// determine if the pattern is an absolute path or not
	if !filepath.IsAbs(pattern) {
		pattern = filepath.Join(parent, pattern)
	}

	// perform a glob pattern match on the pattern to get the files
	filelist, err = filepath.Glob(pattern)

	return filelist, err
}
