package eocs

import (
	"github.com/exlskills/eocsutil/wsenv"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var validEnvKeys = map[string]struct{}{
	"java_default_free": {},
	"python_2_7_free":   {},
	"python_3_4_free":   {},
}

type BlockREPL struct {
	APIVersion     int                             `yaml:"api_version"`
	EnvironmentKey string                          `yaml:"environment"`
	SourcePath     string                          `yaml:"src_path"`
	Explanation    string                          `yaml:"explanation,omitempty"`
	TmplPath       string                          `yaml:"tmpl_path,omitempty"`
	TestPath       string                          `yaml:"test_path,omitempty"`
	Display        *BlockREPLDisplay               `yaml:"display,omitempty"`
	SrcFiles       map[string]*wsenv.WorkspaceFile `yaml:"-"`
	TmplFiles      map[string]*wsenv.WorkspaceFile `yaml:"-"`
	TestFiles      map[string]*wsenv.WorkspaceFile `yaml:"-"`
}

type BlockREPLDisplay struct {
	Height string `yaml:"height"`
}

func (repl *BlockREPL) IsAPIVersionValid() bool {
	return repl.APIVersion == 1
}

func (repl *BlockREPL) IsEnvironmentKeyValid() bool {
	if _, exists := validEnvKeys[repl.EnvironmentKey]; exists {
		return true
	}
	return false
}

func (repl *BlockREPL) LoadFilesFromFS(rootDir string) error {
	if repl.SourcePath != "" {
		files, err := loadFilesFromFSForEnv(repl.EnvironmentKey, filepath.Join(rootDir, repl.SourcePath))
		if err != nil {
			return err
		}
		repl.SrcFiles = files
	}
	if repl.TmplPath != "" {
		files, err := loadFilesFromFSForEnv(repl.EnvironmentKey, filepath.Join(rootDir, repl.TmplPath))
		if err != nil {
			return err
		}
		repl.TmplFiles = files
	}
	if repl.TestPath != "" {
		files, err := loadFilesFromFSForEnv(repl.EnvironmentKey, filepath.Join(rootDir, repl.TestPath))
		if err != nil {
			// TODO make this fatal later...
			Log.Warn("Unable to load test directory, despite a path being supplied...")
			//return err
		}
		repl.TestFiles = files
	}
	return nil
}

func loadFilesFromFSForEnv(envKey, dir string) (files map[string]*wsenv.WorkspaceFile, err error) {
	fillinFiles, err := loadFilesFromDirRecursive(dir)
	if err != nil {
		return nil, err
	}
	switch envKey {
	case "java_default_free":
		files = map[string]*wsenv.WorkspaceFile{
			"src": {
				Name:       "src",
				IsDir:      true,
				IsTmplFile: true,
				Children: map[string]*wsenv.WorkspaceFile{
					"main": {
						Name:       "main",
						IsDir:      true,
						IsTmplFile: true,
						Children: map[string]*wsenv.WorkspaceFile{
							"java": {
								Name:       "java",
								IsDir:      true,
								IsTmplFile: true,
								Children: map[string]*wsenv.WorkspaceFile{
									"exlcode": {
										Name:       "exlcode",
										IsDir:      true,
										IsTmplFile: true,
										Children:   fillinFiles,
									},
								},
							},
						},
					},
				},
			},
		}
		return
	default:
		// This is essentially just what we get directly from the FS
		return fillinFiles, nil
	}
}

func loadFilesFromDirRecursive(dir string) (files map[string]*wsenv.WorkspaceFile, err error) {
	dirListing, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files = make(map[string]*wsenv.WorkspaceFile)
	for _, fi := range dirListing {
		if fi.IsDir() {
			f := &wsenv.WorkspaceFile{
				Name:  fi.Name(),
				IsDir: true,
			}
			f.Children, err = loadFilesFromDirRecursive(filepath.Join(dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			files[fi.Name()] = f
		} else {
			bContents, err := ioutil.ReadFile(filepath.Join(dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			files[fi.Name()] = &wsenv.WorkspaceFile{
				Name:     fi.Name(),
				Contents: string(bContents),
			}
		}
	}
	return files, nil
}

func isValidProblemREPLShebang(sb string) bool {
	if strings.HasPrefix(sb, "#!exl::repl('") && strings.HasSuffix(sb, "')") {
		return true
	}
	return false
}

func getProblemREPLPath(shebang string) (path string, err error) {
	if !isValidProblemREPLShebang(shebang) {
		return "", errors.New("invalid problem REPL shebang")
	}
	return filepath.Clean(strings.Replace(strings.Replace(shebang, "#!exl::repl('", "", 1), "')", "", 1)), nil
}
