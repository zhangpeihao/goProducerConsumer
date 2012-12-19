package main
import (
  "fmt"
  "math/rand"
  "time"
  "os/signal"
  "syscall"
  "os"
)

type buffer struct {
  data string
}

type Buffer *buffer

var chanBuffer = make(chan Buffer, 100)
var workChan = make(chan Buffer)
var testStrings = map[string]string{
  "Hello":      "world",
  "Foo":        "bar",
  "Sun":        "moon",
  "Zhang":      "peihao",
}

func productSomething(b Buffer) {
  time.Sleep(time.Duration(rand.Int() % 1000) * time.Microsecond)
  len := len(testStrings)
  randIndex := rand.Int() % len
  i := 0
  for k, _ := range(testStrings) {
    if i++; i == randIndex {
      b.data = k
      break
    }
  }
}

func doProcess(b Buffer) {
  v, f := testStrings[b.data]
  if f {
    fmt.Printf("%s:\t%s\n", b.data, v)
  } else {
    fmt.Printf("%s not found in map\n", b.data)
  }
}

func producer() {
  for {
    var b Buffer
    select {
    case b = <- chanBuffer:
    default:
      b = new(buffer)
    }
    productSomething(b)
    workChan <- b
  }
}

func consumer() {
  for {
    b := <- workChan
    doProcess(b)
    select {
    case chanBuffer <- b:
    default:
    }
  }
}

func main() {
  rand.Seed(time.Now().Unix())
  go consumer()
  for i := 0; i < 3; i++ {
    go producer()
  }
  ch := make(chan os.Signal)
  signal.Notify(ch, syscall.SIGINT)
  _ = <-ch
  fmt.Printf("Exit normally\n")
}
