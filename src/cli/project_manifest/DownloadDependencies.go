package project_manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"cedpm.org/utils"
)


func DownloadDependencies(projectDir string) error {

	resolvedDependencies, err := GetResolvedDependencies(projectDir)
	if err != nil {
		return err
	}

	for i, v := range resolvedDependencies {
		fmt.Println("[i]:", i, " v=", v.Type)
		if v.Type == "local" {
			// move ? dont move ? should move.
			fmt.Println("Path=", v.Path)
		} else {
			// Download package zip url ?
			// When releasing a pakcage, have both Project and Pklproject files ? So that project one would stay
			fmt.Println("Uri=", v.Uri)
			fmt.Println("Downlonading ["+v.Name+"] ..... ")



			remoteString, err := utils.GetRemoteString(utils.ReplacePrefix(v.Uri, "projectpackage://", "https://"))
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
			if err := utils.DownloadFile(packageInfo.PackageZipUrl, fileName); err != nil {
				fmt.Println("Error downloading file:", err)
				return err
			}

			calculatedHash, err := utils.CalculateSHA256(fileName)
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

			if err := utils.ExtractZip(fileName, packageOutputPath); err != nil {
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
