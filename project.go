package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Project struct {
	RepoURL       string
	Directory     string
	CloneDir      string
	Layer         int
	ProjectOnDisk bool
	DefaultBranch string
	State         string
	Name          string
}

func addProject(directory string, url string, name string) error {
	c, configErr := loadUserConfig()
	if configErr != nil {
		return configErr
	}
	c.Projects = append(c.Projects, Project{RepoURL: url, Directory: directory, Name: name, DefaultBranch: "main", State: "added"})
	saveUserConfig(*c)

	return nil
}

func cloneProjects() error {

	triples, configErr := loadUserTurtle()

	if configErr != nil {
		return configErr
	}
	ns := os.Getenv("O8_META_NAMESPACE")
	root := os.Getenv("O8ROOT")
	var projectList []*Project
	projects := findProjects(triples, ns)

	for _, spo := range projects {
		pro := &Project{}
		pro.DefaultBranch = "main"
		for _, project := range spo {
			_, p, o := trimNs(ns, project)
			switch p {
			case "projectLayer":
				i, _ := strconv.Atoi(o)
				pro.Layer = i

			case "projectOnDisk":
				pro.ProjectOnDisk = o == "true"

			case "projectRepository":
				pro.RepoURL = o

			case "projectFilePath":
				pro.Directory = o

			case "projectDefaultBranch":
				pro.DefaultBranch = o

			case "rootFilePath":
				pro.CloneDir = o

			}

		}
		projectList = append(projectList, pro)
	}

	for _, pr := range projectList {
		if !pr.ProjectOnDisk {
			dir := filepath.Join(root, pr.Directory)

			print(dir, "WILL CLONE")
			_, err := git.PlainClone(dir, false, &git.CloneOptions{
				URL: pr.RepoURL,

				Progress: os.Stdout,
			})

			if err != nil {
				fmt.Printf("Error: %s", err)
				return err
			}
		}

	}

	// saveUserConfig(*c)

	return nil
}

func findProjects(list RDFTripleList, ns string) BySubjectUnprefixed {
	m := groupBySubject(list, ns)
	toReturn := make(BySubjectUnprefixed)

	for subject, triples := range m {

		if subjectIsProject(triples, ns, m) {
			print(subject, "PROJECT")
			toReturn[subject] = triples
			for _, project := range triples {
				_, p, o := trimNs(ns, project)
				print(p, o)

			}
		}
	}
	return toReturn

}

type BySubjectUnprefixed map[string]RDFTripleList

func groupBySubject(list RDFTripleList, ns string) BySubjectUnprefixed {
	m := make(BySubjectUnprefixed)
	for _, t := range list {
		s, _, _ := trimNs(ns, t)
		m[s] = append(m[s], t)
	}

	return m
}

func subjectIsProject(subjectTriples RDFTripleList, ns string, m BySubjectUnprefixed) bool {
	response := false
	for _, t := range subjectTriples {
		o, b := isProject(ns, t)

		if b {
			response = true
		}
		for _, g := range m[o] {
			if isSubclassOfProject(ns, g) {
				response = true
			}
		}

	}

	return response
}

func isProject(ns string, t RDFTriple) (string, bool) {
	_, p, o := trimNs(ns, t)
	return o, (p == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" && trimNsValue(ns, o) == "Project")
}

func isSubclassOfProject(ns string, t RDFTriple) bool {
	_, p, o := trimNs(ns, t)
	return (p == "http://www.w3.org/2000/01/rdf-schema#subClassOf" && trimNsValue(ns, o) == "Project")
}

func trimNs(ns string, triple RDFTriple) (string, string, string) {
	var s, p, o string
	toRemove := ns + "/"
	s = strings.Replace(triple.Subject, toRemove, "", 1)
	p = strings.Replace(triple.Predicate, toRemove, "", 1)
	o = strings.Replace(triple.Object, toRemove, "", 1)

	return s, p, o
}

func trimNsValue(ns string, term string) string {
	toRemove := ns + "/"
	return strings.Replace(term, toRemove, "", 1)
}
