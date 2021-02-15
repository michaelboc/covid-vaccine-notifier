//
// Program scrapes for covid availability in MA state run facilities.
//

package main

import (
    "os"
    "io/ioutil"
    "strings"
    "net/http"
    "log" 
    "github.com/PuerkitoBio/goquery" 
    "gopkg.in/yaml.v2"
)

func LookupSiteAvailability(vaccineURL string, vaccineSites map[string]bool) (sucess, foundAvailableSite bool) {
    foundAvailableSite = false

    response, err := http.Get(vaccineURL)
    if err != nil {
        log.Fatal(err)
        return false, false
    }
    defer response.Body.Close()

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal(err)
        return false, false
     }

    doc.Find(".table tr").Each(func(i int, rowNode *goquery.Selection){

        for key, _ := range vaccineSites {
            siteLocation := rowNode.Find("a").Text()

            if strings.Contains(siteLocation, key) {
                rowNode.Find(".text").Each(func(j int, columnNode *goquery.Selection){

                    if strings.Contains(columnNode.Text(), "Available") {
                        vaccineSites[key] = true
                        foundAvailableSite = true 
                    }
                })
            }
        }
    })
    return true, foundAvailableSite
}

type ProgConf struct {
    LogFile         string      `yaml:"logFile"`
    VaccineURL      string      `yaml:"vaccineURL"` 
    SiteTargets     []string    `yaml:"siteTargets,flow"` 
    EmailTargets     []string    `yaml:"emailTargets,flow"` 
    EmailHost       string      `yaml:"emailHost"` 
    EmailFrom       string      `yaml:"emailFrom"` 
    EmailPassword   string      `yaml:"emailPassword"` 
    EmailPort       string      `yaml:"emailPort"` 
}

func ReadConfigurationFile(filepath string, conf *ProgConf) (bool) {
    yamlFile, err := ioutil.ReadFile(filepath)
    if err != nil {
        log.Fatal(err)
        return false 
    }
    err = yaml.Unmarshal(yamlFile, conf)
    if err != nil {
        log.Fatal(err)
        return false 
    }     
    return true
}


func main() {
    // Get program vars
    var conf ProgConf 
    ReadConfigurationFile(os.Args[1], &conf)  

    // Set up logging to file
    f, err := os.OpenFile(conf.LogFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)

    siteData := make(map[string]bool)
    for i := 0; i < len(conf.SiteTargets); i++ {
        siteData[conf.SiteTargets[i]] = false
    }

    sucess, foundSite := LookupSiteAvailability(conf.VaccineURL, siteData) 
    if sucess && foundSite {
        log.Println("Found a valid site") 
        emailVars := EmailVars{conf.EmailHost, conf.EmailFrom, conf.EmailTargets, conf.EmailPassword, conf.EmailPort}
        SendEmailSMTP(GenerateEmailBody(siteData), emailVars )
    } else {
        log.Println("No valid sites found.")
    } 
}
