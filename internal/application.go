package internal

import (
	"encoding/json"
	"fmt"
	"github.com/npcnixel/genanki-go"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const DataDirectory = "data"

var Levels = []string{"A1", "A2"}

type Application struct {
	logger *zap.Logger

	llm    *llmClient
	result map[string][]*VerbRecord
}

type VerbRecord struct {
	Infinitive  string       `json:"infinitive"`
	Present     string       `json:"present"`
	Past        string       `json:"past"`
	Translation *Translation `json:"translation"`
	Examples    []*Example   `json:"examples"`
}

type Translation struct {
	English string `json:"en"`
	Russian string `json:"ru"`
}

type Example struct {
	Sentence    string       `json:"sentence"`
	Translation *Translation `json:"translation"`
}

func NewApplication(options ...Option) *Application {
	a := &Application{}
	a.logger, _ = zap.NewDevelopment()
	a.result = map[string][]*VerbRecord{}

	for _, option := range options {
		option(a)
	}

	return a
}

type Option func(application *Application)

func WithLLM() Option {
	return func(a *Application) {
		a.llm = newLLMClient(a.logger)
	}
}

func (a *Application) WriteResult() {
	for _, level := range Levels {
		res := a.result[level]

		sort.Slice(res,
			func(i, j int) bool {
				return res[i].Infinitive < res[j].Infinitive
			},
		)

		bytes, _ := json.MarshalIndent(res, "", "  ")

		f, _ := os.Create(filepath.Join(DataDirectory, fmt.Sprintf("%s.json", level)))
		_, _ = f.Write(bytes)
		_ = f.Close()
	}
}

func (a *Application) LoadResult() {
	var bytes []byte
	var err error

	for _, level := range Levels {
		bytes, err = os.ReadFile(filepath.Join(DataDirectory, fmt.Sprintf("%s.json", level)))
		if err != nil {
			panic(err)
		}

		tmp := []*VerbRecord{}
		if err = json.Unmarshal(bytes, &tmp); err != nil {
			panic(err)
		}
		a.result[level] = tmp
	}
}

func (a *Application) ExtractVerbsFromImages() {
	for _, level := range Levels {
		path := filepath.Join(DataDirectory, "images", level)
		a.result[level] = []*VerbRecord{}
		images, err := os.ReadDir(
			path,
		)
		if err != nil {
			a.logger.Warn("Failed to read images", zap.String("level", level), zap.Error(err))
		}

		for _, image := range images {
			data := a.llm.extractVerbsFromImage(
				filepath.Join(path, image.Name()),
			)
			a.logger.Info("processing image",
				zap.String("level", level),
				zap.String("image", image.Name()),
				zap.Int("count", len(data)),
			)
			a.result[level] = append(a.result[level], data...)
		}
	}
}

func (a *Application) AddVerbExamples() {
	var (
		err  error
		curr *VerbRecord
		cnt  int
		lvl  int
	)

	lvl = 1
	for _, level := range Levels {
		cnt = 1
		for idx, verb := range a.result[level] {
			curr, err = a.llm.addExampleSentences(verb)
			cnt++
			a.logger.Info(
				fmt.Sprintf("%d/%d; %d/%d", lvl, len(Levels), cnt, len(a.result[level])),
				zap.Error(err),
			)
			if err == nil {
				a.result[level][idx] = curr
			}
		}
		lvl++
	}
	a.logger.Info("examples added",
		zap.Int("total verbs", cnt),
	)
}

func (a *Application) MakeAnkicards() {
	model := genanki.StandardBasicModel("verbs")

	for _, level := range Levels {
		deck := genanki.StandardDeck(
			fmt.Sprintf("German verbs %s", level),
			"List of German verbs with details for given level",
		)
		for _, verb := range a.result[level] {
			note := genanki.NewNote(
				model.ID,
				[]string{
					verb.Infinitive,
					verb.Translation.English,
					verb.Translation.Russian,
					verb.Present,
					verb.Past,
					func(examples []*Example) string {
						res := []string{}
						for _, e := range examples {
							res = append(res,
								[]string{e.Sentence, e.Translation.English}...,
							)
						}
						return strings.Join(res, "\n")
					}(verb.Examples),
				},
				[]string{
					"infinitive",
					"english",
					"russian",
					"present",
					"past",
					"examples",
				},
			)
			deck.AddNote(note)

		}
		pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model.Model)
		if err := pkg.WriteToFile(
			filepath.Join(DataDirectory, fmt.Sprintf("%s.apkg", level)),
		); err != nil {
			a.logger.Info(
				"Failed to write package",
				zap.String("level", level),
				zap.Error(err),
			)
		}
	}
}
