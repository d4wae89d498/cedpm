package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"io"
	"errors"
	"net/http"
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
_	"cedpm.org/internal"
)

/*\
 *		E x p o s e d
\*/

type Command struct {
	Before  	[]string
	After   	[]string
	On      	[]string
	Actions 	[]string
}

type CommandList struct {
	Paths     	[]string
	Commands 	map[string]Command
}

var GCommandList CommandList

/*\
 *		./Project 	F i l e
\*/

type projectFileAddons struct {
	Paths		[]string			`json:"paths"`
	Before  	map[string][]string	`json:"before"`
	After   	map[string][]string	`json:"after"`
	On      	map[string][]string	`json:"on"`
	Commands 	map[string][]string	`json:"commands"`
}

/*\
 *		./.dependencies.json 	F i l e
\*/

type dependencyManifest struct {
    SchemaVersion       int          	`json:"schemaVersion"`
    ResolvedDependencies []dependency	`json:"resolvedDependencies"`
}

type dependency struct {
    Name      			string     		`json:"-"`
    Type      	string   			`json:"type"`
    Uri       	string   			`json:"uri"`
    Checksums 	checksums 			`json:"checksums,omitempty"`
    Path      	string  			`json:"path,omitempty"`
}

type checksums struct {
    Sha256 		string 				`json:"sha256"`
}

/*\
 *		PKL - Generated  remote project JSON
\*/

type PackageInfo struct {
	Name 				 string 	`json:"name"`
	PackageZipUrl        string		`json:"packageZipUrl"`
	PackageZipChecksums  struct {
		Sha256 string 				 `json:"sha256"`
	} 								 `json:"packageZipChecksums"`
}

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	return err
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func extractZip(filePath, destDir string) error {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		fPath := fmt.Sprintf("%s/%s", destDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}
		dir := filepath.Dir(fPath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func getRemoteString(url string) (string, error) {
	// Send an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err // Handle errors connecting to the server
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status) // Handle non-OK HTTP responses
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err // Handle errors reading the response body
	}

	// Convert the body bytes to a string
	bodyString := string(bodyBytes)
	return bodyString, nil
}

func replacePrefix(s, prefix, newPrefix string) string {
	if strings.HasPrefix(s, prefix) {
		// TrimPrefix removes the prefix and then we add the newPrefix.
		return newPrefix + strings.TrimPrefix(s, prefix)
	}
	return s
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func copyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if !os.IsNotExist(err) {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//////////////////////////////////////////////////////////////


func ParseProjectDeps(projectDir string) error {
	fmt.Printf("Reading read deps json file !\n")
	var err error

	filePath := projectDir + "/.dependencies.json" // Constructing the file path

	fmt.Println(filePath)

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully read deps json file !\n")


	 // Assuming jsonStr is your JSON string
	 var rawManifest map[string]interface{}

	 err = json.Unmarshal([]byte(jsonData), &rawManifest)
	 if err != nil {
		 fmt.Println("Error parsing JSON:", err)
		 return err
	 }

	 manifest := dependencyManifest{
		 SchemaVersion: int(rawManifest["schemaVersion"].(float64)),
	 }

	 resolvedDeps := rawManifest["resolvedDependencies"].(map[string]interface{})
	 for name, pkg := range resolvedDeps {
		 pkgBytes, err := json.Marshal(pkg)
		 if err != nil {
			 fmt.Println("Error re-marshalling package:", err)
			 return err //continue
		 }

		 var packageItem dependency
		 err = json.Unmarshal(pkgBytes, &packageItem)
		 if err != nil {
			 fmt.Println("Error unmarshalling package:", err)
			 return err //continue
		 }

		 packageItem.Name = name
		 manifest.ResolvedDependencies = append(manifest.ResolvedDependencies, packageItem)
	 }

	// fmt.Printf("%+v\n", manifest)

	fmt.Println("Project deps:", manifest.ResolvedDependencies)

	for i, v := range manifest.ResolvedDependencies {
		fmt.Println("[i]:", i, " v=", v.Type)
		if v.Type == "local" {
			// move ? dont move ? should move.
			fmt.Println("Path=", v.Path)
		} else {
			// Download package zip url ?
			// When releasing a pakcage, have both Project and Pklproject files ? So that project one would stay
			fmt.Println("Uri=", v.Uri)
			fmt.Println("Downlonading ["+v.Name+"] ..... ")



			remoteString, err := getRemoteString(replacePrefix(v.Uri, "projectpackage://", "https://"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching remote string: %v\n", err)
				return err
			}
			fmt.Println(remoteString)
			var packageInfo PackageInfo
			if err := json.Unmarshal([]byte(remoteString), &packageInfo); err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return err
			}

			fileName := "package.zip"
			if err := downloadFile(packageInfo.PackageZipUrl, fileName); err != nil {
				fmt.Println("Error downloading file:", err)
				return err
			}

			calculatedHash, err := calculateSHA256(fileName)
			if err != nil {
				fmt.Println("Error calculating SHA256:", err)
				return err
			}

			if calculatedHash != packageInfo.PackageZipChecksums.Sha256 {
				fmt.Println("Checksum verification failed.")
				return err
			}

			packageOutputPath := "./.cedpm/packages/" + packageInfo.Name + "/"

			err = os.RemoveAll(packageOutputPath)
			if err != nil {
				// Handle the error
				fmt.Println("Error removing the directory:", err)
				return err
			}

			if err := extractZip(fileName, packageOutputPath); err != nil {
				fmt.Println("Error extracting ZIP file:", err)
				return err
			}

			err = os.Remove(fileName)
			if err != nil {
				// Handle the error, such as logging or printing it
				fmt.Println("Error deleting the file:", err)
				return err
			}

			fmt.Println("Package extracted successfully.")

		}
	}
	return nil
}

func ParseProjectManifest(jsonData string) error {
	var addons projectFileAddons
	if err := json.Unmarshal([]byte(jsonData), &addons); err != nil {
		//("Error parsing JSON:", err)
		return err
	}

	fmt.Println("Project File addons:", addons)


	GCommandList.Paths = append(
		GCommandList.Paths,
		addons.Paths...
	)

	for key, value := range addons.Commands {
		_, exists := GCommandList.Commands[key]
		if exists {
			return errors.New("Err")
		}
		_ = value
		_ = key
    }

	return nil
}
