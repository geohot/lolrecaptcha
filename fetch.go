package main

import (
  "bytes"
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "image"
  "image/jpeg"
  "image/png"
  "io"
  "log"
  "os"
  "golang.org/x/net/html"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
)

const apiKey string = "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-"

func fetchImg(ck string) (image.Image) {
  // fetch the image
  u, err := url.Parse("http://google.com/recaptcha/api2/payload")
  if err != nil { log.Fatal(err) }
  q := u.Query()
  q.Set("c", ck)
  q.Set("k", apiKey)
  u.RawQuery = q.Encode()
  //fmt.Println(u)

  // do fetch
  imgresponse, err := http.Get(u.String())
  if err != nil { log.Fatal(err) }

  img, err := jpeg.Decode(imgresponse.Body)
  if err != nil { log.Fatal(err) }

  return img
}

func getChallengeKey() (string, string, image.Image) {
  // build the request
  u, err := url.Parse("http://google.com/recaptcha/api/fallback")
  if err != nil { log.Fatal(err) }
  q := u.Query()
  q.Set("k", apiKey)
  u.RawQuery = q.Encode()
  //fmt.Println(u)

  // fetch the webpage
  response, err := http.Get(u.String())
  if err != nil { log.Fatal(err) }
  defer response.Body.Close()

  // print it
  bodyBytes, _ := ioutil.ReadAll(response.Body)

  z := html.NewTokenizer(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
  tmparr := []string{}
  ck := ""
  for {
    tt := z.Next()
    switch tt {
      case html.ErrorToken:
        return ck, tmparr[3], fetchImg(ck)
      case html.StartTagToken, html.SelfClosingTagToken:
        tn, attr := z.TagName()
        if string(tn) == "img" && attr {
          for {
            k, v, attr := z.TagAttr()
            if string(k) == "src" {
              //fmt.Println(string(v))
              u, err := url.Parse(string(v))
              if err != nil { log.Fatal(err) }
              q := u.Query()
              //fmt.Println(q)
              if q["k"][0] != apiKey { log.Fatal("apiKey doesn't match") }
              ck = q["c"][0]
            }
            if !attr { break }
          }
        }
      case html.TextToken:
        //fmt.Println(z.Token())
        tmparr = append(tmparr, z.Token().String())
    }
  }
}

func main() {
  fmt.Println("my first golang program")

  for {
    // parse it
    ck, typ, img := getChallengeKey()

    h := md5.New()
    io.WriteString(h, ck)
    hh := hex.EncodeToString(h.Sum(nil))
    typ = strings.Replace(typ, " ", "_", -1)
    //fmt.Println(ck, typ, img.Bounds())
    fmt.Println(hh, typ, img.Bounds())

    // write it
    os.MkdirAll("imgs/"+typ, 0755)
    f, err := os.OpenFile("imgs/"+typ+"/"+hh+".png", os.O_CREATE | os.O_WRONLY, 0644)
    if err != nil { log.Fatal(err) }
    png.Encode(f, img)
    f.Close()
  }

  //if typ == "street signs" { break }

  /*if false {
    //fmt.Println(img.SubImage(image.Rect(0,0,100,100)).Bounds())
    _ = image.Rect(0,0,100,100)
  }
  fmt.Println(img.At(0,0))*/

  /*bb, _ := ioutil.ReadAll(imgresponse.Body)
  if false { fmt.Println(bb) }
  fmt.Println(bb[0:5])*/

  //if err != nil { log.Fatal(err) }
  //defer imgresponse.Body.Close()

  // move on
  fmt.Println("still alive!")
}

