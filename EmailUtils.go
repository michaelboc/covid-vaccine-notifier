package main

import (
    "log" 
    "net/smtp" 
)

type EmailVars struct {
    Host            string
    SourceAccnt     string 
    EmailTargets    []string
    Password        string
    Port            string
} 

func SendEmailSMTP(emailBody string, emailVars EmailVars) (bool) {
    emailAuth := smtp.PlainAuth("", emailVars.SourceAccnt, emailVars.Password, emailVars.Host)

    const subject = "Subject: COIVD-19 Possible Appointment Notification Email\n"
    msg := []byte(subject + "\n" + emailBody)

    err := smtp.SendMail(emailVars.Host + ":" + emailVars.Port, emailAuth, emailVars.SourceAccnt, emailVars.EmailTargets, msg)
    if err != nil {
        log.Fatal(err)
        return false
    }
    return true
}

func GenerateEmailBody(vaccineSites map[string]bool) (emailBody string) {
    emailBody = "One or more of the requested vaccination site(s) have availablities:\n" +
    "---------\n"

    for key, value := range vaccineSites {
        var siteState string
        if value {
            siteState = "AVAILABLE"
        } else {
            siteState = "UNAVAILABLE"
        } 
        emailBody += key + ": " + siteState + "\n" 
    }
    return
}
