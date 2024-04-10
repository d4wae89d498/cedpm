package project_manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*\
 *		./.dependencies.json 	F i l e
\*/

type dependencyManifest struct {
    SchemaVersion       int          	`json:"schemaVersion"`
    ResolvedDependencies []Dependency	`json:"resolvedDependencies"`
}

type Dependency struct {
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

///////////////////////////////////////

func GetResolvedDependencies(projectDir string) ([]Dependency, error) {
	fmt.Printf("Reading read deps json file !\n")
	var err error

	filePath := projectDir + "/.dependencies.json" // Constructing the file path

	fmt.Println(filePath)

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Successfully read deps json file !\n")


	 // Assuming jsonStr is your JSON string
	 var rawManifest map[string]interface{}

	 err = json.Unmarshal([]byte(jsonData), &rawManifest)
	 if err != nil {
		 fmt.Println("Error parsing JSON:", err)
		 return nil, err
	 }

	 manifest := dependencyManifest{
		 SchemaVersion: int(rawManifest["schemaVersion"].(float64)),
	 }

	 resolvedDeps := rawManifest["resolvedDependencies"].(map[string]interface{})
	 for name, pkg := range resolvedDeps {
		 pkgBytes, err := json.Marshal(pkg)
		 if err != nil {
			 fmt.Println("Error re-marshalling package:", err)
			 return nil, err //continue
		 }

		 var packageItem Dependency
		 err = json.Unmarshal(pkgBytes, &packageItem)
		 if err != nil {
			 fmt.Println("Error unmarshalling package:", err)
			 return nil, err //continue
		 }

		 packageItem.Name = name
		 manifest.ResolvedDependencies = append(manifest.ResolvedDependencies, packageItem)
	 }
	return manifest.ResolvedDependencies, nil
}
