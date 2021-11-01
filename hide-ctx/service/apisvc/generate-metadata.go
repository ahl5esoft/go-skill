package apisvc

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"
)

const (
	metadataTpl = `package api

import (
    "github.com/ahl5esoft/go-skill/hide-ctx/contract"
    {{- range .packages }}
    {{ .Name }} "{{ $.workspace }}/{{ join .RelativePathParts "/" }}"
    {{- end }}
)

func Register(apiFactory contract.IApiFactory) {
	{{- range $i, $r := .packages }}{{ range $ci, $cr := $r.Apis }}
    apiFactory.Register("{{ $r.Endpoint }}", "{{ $cr.Route }}", {{ $r.Name }}.{{ $cr.Struct }}Api{}){{ end }}{{ end }}
}`
	metadataFilename = "metadata.go"
)

var (
	regApi   = regexp.MustCompile(`type\s(\w+)Api`)
	tplFuncs = template.FuncMap{
		"join": func(elems []string, sep string) string {
			return strings.Join(elems, sep)
		},
	}
)

type apiData struct {
	Struct string
	Route  string
}

type packageData struct {
	Apis              []apiData
	Endpoint          string
	Name              string
	RelativePathParts []string
}

func GenerateMetadata(ioFactory contract.IIOFactory, ioPath contract.IIOPath, workspace string) (err error) {
	packages := make([]packageData, 0)
	apiDir := ioFactory.BuildDirectory(
		ioPath.GetRoot(),
		"api",
	)
	err = readGoFiles(apiDir, &packages, workspace)
	if err != nil {
		return
	}

	var tpl *template.Template
	if tpl, err = template.New("").Funcs(tplFuncs).Parse(metadataTpl); err != nil {
		return
	}

	var bf bytes.Buffer
	err = tpl.Execute(&bf, map[string]interface{}{
		"packages":  packages,
		"workspace": workspace,
	})
	if err != nil {
		return
	}

	err = ioFactory.BuildFile(
		apiDir.GetPath(),
		metadataFilename,
	).Write(bf)
	return
}

func readGoFiles(dir contract.IIODirectory, packages *[]packageData, workspace string) (err error) {
	files := dir.FindFiles()

	apis := make([]apiData, 0)
	for _, r := range files {
		if r.GetExt() != ".go" || r.GetName() == metadataFilename {
			continue
		}

		isTest := strings.Contains(
			r.GetName(),
			"_test",
		)
		if isTest {
			continue
		}

		api := apiData{
			Route: strings.Replace(
				r.GetName(),
				r.GetExt(),
				"",
				1,
			),
		}

		var text string
		if err = r.Read(&text); err != nil {
			return
		}

		matches := regApi.FindStringSubmatch(text)
		if len(matches) == 0 {
			continue
		}

		api.Struct = matches[1]
		apis = append(apis, api)
	}

	if len(apis) > 0 {
		pkg := packageData{
			Apis:              apis,
			RelativePathParts: make([]string, 0),
		}
		var temp contract.IIODirectory
		for {
			if len(pkg.RelativePathParts) == 0 {
				temp = dir
			} else {
				temp = temp.GetParent()
			}

			if temp.GetName() == workspace {
				break
			}

			pkg.RelativePathParts = append([]string{
				temp.GetName(),
			}, pkg.RelativePathParts...)
		}

		if pkg.RelativePathParts[len(pkg.RelativePathParts)-2] == "api" {
			pkg.Endpoint = pkg.RelativePathParts[len(pkg.RelativePathParts)-1]
			pkg.Name = pkg.RelativePathParts[len(pkg.RelativePathParts)-1]
		} else {
			pkg.Endpoint = strings.Join(
				pkg.RelativePathParts[len(pkg.RelativePathParts)-2:],
				"/",
			)
			pkg.Name = strings.Join(
				pkg.RelativePathParts[len(pkg.RelativePathParts)-2:],
				"",
			)
		}
		pkg.Name = strings.Replace(pkg.Name, "-", "", -1)

		*packages = append(*packages, pkg)
	}

	childDirs := dir.FindDirectories()
	if len(childDirs) == 0 {
		return
	}

	for _, r := range childDirs {
		readGoFiles(r, packages, workspace)
	}
	return
}
