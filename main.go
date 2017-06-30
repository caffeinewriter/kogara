package main

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  
  "github.com/gin-gonic/gin"
  "github.com/go-redis/redis"
  "github.com/pilu/go-base62"
)

type Configuration struct {
  Redis RedisConfig `json:"redis"`
  Bind string `json:"bind"`
}

type RedisConfig struct {
  Addr string `json:"addr"`
  Password string `json:"password"`
  DB int `json:"db"`
}

func LoadConfig() (config Configuration) {
  file, err := ioutil.ReadFile("./config.json")
  if err != nil {
    panic("Unable to read config file.")
  }
  var cfg Configuration
  err = json.Unmarshal(file, &cfg)
  if err != nil {
    panic("Unable to parse config file.")
  }
  return cfg
}

func main() {
  config := LoadConfig()
  r := gin.Default()

  red := redis.NewClient(&redis.Options{
    Addr:     config.Redis.Addr,
    Password: config.Redis.Password,
    DB:       config.Redis.DB,
  })
  
  _, err := red.Ping().Result()
  
  if err != nil {
    panic(err)
  }
  
  r.GET("/", func (ctx *gin.Context) {
    ctx.Header("Content-Type", "text/html")
    ctx.String(200, `<!DOCTYPE html>
                     <html>
                      <head>
                        <title>Kogara - Link Shortening</title>
                      </head>
                      <body>
                        <form action="/" method="POST">
                          <input type="text" name="url"><br>
                          <input type="submit" value="Shorten">
                        </form>
                      </body>
                      </html>`);
  })
  
  r.POST("/", func (ctx *gin.Context) {
    if ctx.PostForm("url") == "" {
      ctx.String(400, "URL cannot be empty.")
    }
    count, err := red.Incr("kogara:counter").Result()
    if err != nil {
      ctx.String(500, "Unable to create link.")
    }
    id := base62.Encode(int(count));
    _, err = red.Set("kogara:links:" + id, ctx.PostForm("url"), 0).Result()
    if err != nil {
      ctx.String(500, "Unable to create link.")
    } else {
      ctx.Header("Content-Type", "text/html")
      ctx.String(200,fmt.Sprintf(`<!DOCTYPE html>
                      <html>
                      <head>
                        <title>Kogara - Link Created</title>
                      </head>
                      <body>
                        Created <a href="/r/%s">/r/%s</a>
                      </body>
                      </html>`, id, id))
    }
  })
  
  r.GET("/r/:id", func (ctx *gin.Context) {
    link, err := red.Get("kogara:links:" + ctx.Params.ByName("id")).Result()
    if err == redis.Nil {
      ctx.String(404, "This link does not seem to exist.")
    } else if err != nil {
      ctx.String(500, "Sorry, something went wrong while trying to fetch this link.")
    } else {
      red.Incr("kogara:analytics:" + ctx.Params.ByName("id"))
      ctx.Redirect(301, link)
    }
  })
  
  r.GET("/+/:id", func (ctx *gin.Context) {
    count, err := red.Get("kogara:analytics:" + ctx.Params.ByName("id")).Result()
    if err == redis.Nil {
      ctx.String(404, "Not found")
    } else if err != nil {
      ctx.String(500, "Something went wrong.")
    } else {
      ctx.String(200, count)
    }
  })
  
  r.GET("/check/:id", func (ctx *gin.Context) {
    exists, err := red.Exists("kogara:links:" + ctx.Params.ByName("id")).Result()
    if err != nil {
      ctx.JSON(500, gin.H{"error": true, "available": nil})
    } else {
      ctx.JSON(200, gin.H{"error": nil, "available": exists != 1})
    }
  })
  
  r.Run(config.Bind)
}
