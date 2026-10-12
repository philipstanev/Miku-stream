package nameparser

import (
	"errors"
	"regexp"
	"strings"
)

type Episode struct {
	series           string
	name             string
	seasonAndEpisode string
}

func ParseToStruct(fName string) (Episode, error) {
	output := Episode{}
	splitNoMetadata := removeMetadata(fName)
	episodeArr, err := fullEpisodeArr(splitNoMetadata)
	if err != nil {
		return Episode{}, err
	}

	output.series = episodeArr[0]
	output.seasonAndEpisode = episodeArr[1]
	output.name = episodeArr[2]

	return output, nil
}

func removeMetadata(fName string) []string {
	splitStr := strings.Split(fName, " ")
	//pattern := regexp.MustCompile(`^\[(.*)\]$`)
	pattern := regexp.MustCompile(`\[|\]`)
	splitStr[len(splitStr)-1] = strings.TrimSuffix(splitStr[len(splitStr)-1], ".mkv")
	output := []string{}
	for _, str := range splitStr {
		if !pattern.MatchString(str) && str != "-" && str != "" {
			output = append(output, str)
		}
	}
	return output
}

func searchForSeasonEpisodeString(split []string) int {
	pattern := regexp.MustCompile("^(S[0-9]+E[0-9]+|[0-9]+x[0-9]+)$")
	for i, str := range split {
		if pattern.MatchString(str) {
			return i
		}
	}
	return -1
}

func fullEpisodeArr(split []string) ([]string, error) {
	output := [3]string{}
	//Index of season-episode string
	mid := searchForSeasonEpisodeString(split)
	if mid == -1 {
		return []string{}, errors.New("Format unsuported")
	}

	output[1] = split[mid]

	for i := 0; i < mid; i++ {
		output[0] += split[i] + " "
	}

	for i := mid + 1; i < len(split); i++ {
		output[2] += split[i] + " "
	}

	output[0] = strings.TrimSpace(output[0])
	output[2] = strings.TrimSpace(output[2])

	return output[:], nil
}
