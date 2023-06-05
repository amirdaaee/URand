package rng

import (
	config "URand/Config"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/splode/fname"
)

func GenerateRandom(n uint, rng *fname.Generator) (*[]string, error) {
	var phraseList []string
	for i := uint(0); i < n; i++ {
		phrase, err := rng.Generate()
		if err != nil {
			logrus.WithField("generated", phraseList).WithError(err).Error("error while generating random phrase")
			return (&phraseList), nil
		}
		phraseList = append(phraseList, phrase)
	}
	return (&phraseList), nil
}
func GetGenerator(l uint) *fname.Generator {
	seed := int64(1)
	if !config.Config().FixSeed {
		seed = int64(time.Now().Second())
	}
	return fname.NewGenerator(fname.WithSize(l), fname.WithSeed(seed))
}
