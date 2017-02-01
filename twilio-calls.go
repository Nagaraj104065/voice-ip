package main

import (
  "encoding/xml"
  "net/http"
  "net/url"
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type TwiML struct {
  XMLName xml.Name `xml:"Response"`

  Say    string `xml:",omitempty"`
  Play   string `xml:",omitempty"`
}

func main() {
  http.HandleFunc("/twiml", twiml)
  http.HandleFunc("/call", call);
  http.ListenAndServe(":3000", nil)
}

func twiml(w http.ResponseWriter, r *http.Request) {
  twiml := TwiML{Play: "http://demo.rickyrobinett.com/huh.mp3"}

  x, err := xml.Marshal(twiml)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/xml")
  w.Write(x)
}

func call(w http.ResponseWriter, r *http.Request) {
  // Let's set some initial default variables
  accountSid := "A*********************************"
  authToken := "8*******************************"
  urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls.json"

  // Build out the data for our message
  v := url.Values{}
  v.Set("To","+91**********")
  v.Set("From","+91**********")
  v.Set("Url","http://127.0.0.1:4040/twiml")
  rb := *strings.NewReader(v.Encode())

// Create Client
  client := &http.Client{}

  req, _ := http.NewRequest("POST", urlStr, &rb)
  req.SetBasicAuth(accountSid, authToken)
  req.Header.Add("Accept", "application/json")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

  // make request
  resp, _ := client.Do(req)
  if( resp.StatusCode >= 200 && resp.StatusCode < 300 ) {
    var data map[string]interface{}
    bodyBytes, _ := ioutil.ReadAll(resp.Body)
    err := json.Unmarshal(bodyBytes, &data)
    if( err == nil ) {
      fmt.Println(data["sid"])
    }
  } else {
    fmt.Println(resp.Status);
    w.Write([]byte("Hello World!"))
  }
}