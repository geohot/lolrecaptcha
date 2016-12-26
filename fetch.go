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
  "github.com/disintegration/imaging"
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

func downloader() {
  bigcnt := 0
  for {
    // parse it
    ck, typ, img := getChallengeKey()

    h := md5.New()
    io.WriteString(h, ck)
    hh := hex.EncodeToString(h.Sum(nil))
    typ = strings.Replace(typ, " ", "_", -1)
    //fmt.Println(ck, typ, img.Bounds())
    fmt.Println(bigcnt, hh, typ, img.Bounds())

    if img.Bounds() != image.Rect(0,0,300,300) {
      log.Fatal("IMAGE IS THE WRONG SIZE")
    }

    // write it
    os.MkdirAll("imgs/"+typ, 0755)

    cnt := 0
    for h := 0; h < 300; h += 100 {
      for w := 0; w < 300; w += 100 {
        lilimg := imaging.Crop(img, image.Rect(w,h,w+100,h+100))

        fn := fmt.Sprintf("imgs/%s/%s_%d.png", typ, hh, cnt)
        f, err := os.OpenFile(fn, os.O_CREATE | os.O_WRONLY, 0644)
        if err != nil { log.Fatal(err) }
        png.Encode(f, lilimg)
        f.Close()

        cnt += 1
      }
    }
    bigcnt += 1
  }
}

func main() {
  fmt.Println("my first golang program")

  for i := 0; i < 8; i += 1 {
    go downloader()
  }
  downloader()

  // move on
  fmt.Println("still alive!")
}

