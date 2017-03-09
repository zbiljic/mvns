// Copyright 2016 Nemanja Zbiljic
//

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	description = "CLI tool to search artifacts in maven central repository"
	version     = "0.2.0"
)

var (
	fGroupId     = flag.String("g", "", "specify groupId")
	fArtifactId  = flag.String("a", "", "specify artifactId")
	fVersion     = flag.String("v", "", "specify version")
	fAllVersions = flag.Bool("A", false, "show all versions")
	fMax         = flag.Int("m", 20, "limit number of result")
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "NAME:\n   %s - %s\n", os.Args[0], description)
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "USAGE:\n   %s [options] [query]\n", os.Args[0])
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "VERSION:\n   %s\n", version)
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "OPTIONS:\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

// JSONResponse solrsearch response.
type JSONResponse struct {
	Response struct {
		Docs []struct {
			ID            string
			LatestVersion string
			Timestamp     int64
		}
	}
}

var maxVLen int

func collect(data JSONResponse) {
	// calculate row size
	for _, d := range data.Response.Docs {
		vlen := 4 // + 4 spaces
		vlen += len(d.ID)
		vlen += len(d.LatestVersion)
		if vlen > maxVLen {
			maxVLen = vlen
		}
	}
	// print results
	for _, d := range data.Response.Docs {
		var line string
		if len(d.LatestVersion) == 0 {
			line = fmt.Sprintf("%s", d.ID)
		} else {
			line = fmt.Sprintf("%s:%s", d.ID, d.LatestVersion)
		}
		fmt.Printf("compile '%s'", color(line))
		fillLine(line)
		fmt.Printf("%6s", msToTime(d.Timestamp).Format("2006-01-02"))
		fmt.Println()
	}
}

func color(s string) string {
	id := strings.Split(s, ":")
	return fmt.Sprintf("%s:%s:%s", colorGroupId(id[0]), colorArtifactId(id[1]), colorVersion(id[2]))
}

func colorGroupId(s string) string {
	return fmt.Sprintf("%s%s%s", "\x1b[32m", s, "\x1b[0m")
}

func colorArtifactId(s string) string {
	return fmt.Sprintf("%s%s%s", "\x1b[35m", s, "\x1b[0m")
}

func colorVersion(s string) string {
	return fmt.Sprintf("%s%s%s%s", "\x1b[34m", "\x1b[1m", s, "\x1b[0m")
}

func msToTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

func fillLine(line string) {
	count := maxVLen - len(line)
	for i := 0; i < count; i++ {
		fmt.Print(" ")
	}
}

func request(params url.Values) (JSONResponse, error) {
	endpoint := "https://search.maven.org/solrsearch/select?" + params.Encode()
	res := JSONResponse{}
	resp, err := http.Get(endpoint)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return res, errors.New(strconv.Itoa(resp.StatusCode))
	}
	// parse JSON with anonymous struct.
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&res); err != nil {
		return res, err
	}
	return res, nil
}

func formatParams(q string) url.Values {

	query := make([]string, 0)

	if len(q) > 0 {
		query = append(query, q)
	}

	groupId := *fGroupId
	if len(groupId) > 0 {
		query = append(query, appendQuery(query, "g", groupId))
	}

	artifactId := *fArtifactId
	if len(artifactId) > 0 {
		query = append(query, appendQuery(query, "a", artifactId))
	}

	version := *fVersion
	if len(version) > 0 {
		query = append(query, appendQuery(query, "v", version))
	}

	params := url.Values{
		"wt":   []string{"json"},
		"rows": []string{strconv.Itoa(*fMax)},
		"q":    []string{strings.Join(query, "")},
	}

	if *fAllVersions {
		params["core"] = []string{"gav"}
	}

	return params
}

func appendQuery(query []string, key, value string) string {
	if len(query) == 0 {
		return fmt.Sprintf("%s:\"%s\"", key, value)
	}
	return fmt.Sprintf(" AND %s:\"%s\"", key, value)
}

func containsNoQueryOptions() bool {
	if len(*fGroupId) > 0 {
		return true
	}
	if len(*fArtifactId) > 0 {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////
// main logic
////////////////////////////////////////////////////////////////////////

func run(args []string) (err error) {

	// Handles cases when too many arguments passed, and no arguments passed
	// but also no (then) required options
	if (len(args) > 2) || (len(args) < 1 && !containsNoQueryOptions()) {
		err = fmt.Errorf("Usage: %s [options] [query]", os.Args[0])
		return
	}

	var query string

	if len(args) < 1 {
		query = ""
	} else {
		query = args[0]
	}

	params := formatParams(query)

	data, err1 := request(params)
	if err1 != nil {
		fmt.Printf("Returns status code %s. Please try it again.\n", err1)
	}

	collect(data)

	fmt.Println()

	return
}

func main() {
	flag.Usage = Usage

	flag.Parse()

	err := run(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
