package main 

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"
    "io/ioutil"
)

func getCurrentDir() string{
    wd, err := os.Getwd();
    if err != nil {
        user, err := user.Current()
        if err != nil {
            return ""
        }
        wd = user.HomeDir
    }
    return wd
}


type step interface{
    perform() error
}

type createDir struct{
    path string
}

func (this createDir) perform() error{
      return os.Mkdir(this.path, 0777)
}

type createFile struct{
    path string
    content []byte
}

func (this createFile) perform() error{
    return ioutil.WriteFile(this.path, this.content, 0777)
}

type allSteps struct{
    steps []step
}

func (this allSteps) perform() error{
    for _, s := range this.steps {
      err := s.perform()
      if err != nil{
          return err;
      }
    }
    return nil
}

func main() {
    name := os.Args[1];
    wd := getCurrentDir();
    packageDir := filepath.Join(wd, name)
    mainDir := filepath.Join(packageDir, "main")
    mainFilePath := filepath.Join(mainDir, "main.go")
    content := []byte("package main\n\n import (\n \"fmt\"\n)\n func main(){\n fmt.Print(\"Hello\")\n}");

    step1 := createDir { packageDir }
    step2 := createDir { mainDir }
    step3 := createFile { mainFilePath,  content}

    allSteps := allSteps{[]step{step1, step2, step3}}

    err := allSteps.perform()
    if err != nil{
       fmt.Print("Error performing operatiol", err)
       return
    }
    fmt.Print("Successfully created package")

}
