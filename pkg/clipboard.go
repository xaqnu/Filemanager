package clipboard

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Action int

const (
	Cut Action = iota
	Copy
)

type Entry struct {
	path   string
	action Action
}

type Clipboard []Entry

func (c *Clipboard) push(entry Entry) {
	*c = append(*c, entry)
}

func (c *Clipboard) pop() (Entry, error) {
	if len(*c)==0 {return Entry{},errors.New("atempted to pop empty slice")}
	last := (*c)[len(*c)-1]
	*c = (*c)[:len(*c)-1]
	return last, nil
}

func (c *Clipboard) copyhandler(path string)  {
	c.push(Entry{path: path, action: Copy})
}

func (c *Clipboard) cuthandler(path string)  {
	c.push(Entry{path: path, action: Cut})
}

func (c *Clipboard) pastehandler (path string) error {
	pastee, err:=c.pop()
	if err!= nil {log.Print("atempted to paste but there was nothing to paste") }
	err=filesystemcopy(pastee.path, path)
	if err !=nil {return err}
	if pastee.action==Cut{
		err=os.RemoveAll(pastee.path)
	}
	return err
}

func filesystemcopy(source string, target string) error {
	filedata, err:= os.Stat(source)
	if err != nil {return err}
	if filedata.IsDir() {
		return dircopy(source,target)
	}else {
		return filecopy(source,target)
	}


}
func dircopy(source string, target string) error{
	return filepath.WalkDir(source,func(abspath string, direntry os.DirEntry, err error)error{
		parent:=filepath.Dir(source)
		path, err:=filepath.Rel(parent, abspath)
		if err!=nil {log.Fatal("parent is wrong what the hell is wrong with you",err)}
		if direntry.IsDir(){
			info, err := direntry.Info()
			if err !=nil {return filepath.SkipDir}
			err=os.Mkdir(filepath.Join(target,path),info.Mode().Perm())	
			if err !=nil {return filepath.SkipDir}
			return nil
		}
		err= filecopy(abspath, filepath.Join(target,filepath.Dir(path))) 
		log.Println("file failed to copy", err)
		return nil

	})
}

func filecopy (source string, target string) error{
	file, err :=os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()
	outcandidate:=filepath.Join(target,filepath.Base(source))
	i:=0
	for{
		if _,err :=os.Stat(outcandidate);errors.Is(err, os.ErrNotExist){
			break
		}
		i+=1
		base:= filepath.Base(source)
		extension:=filepath.Ext(base)
		name:= strings.TrimSuffix(base,extension)
		outcandidate=filepath.Join(target,fmt.Sprintf("%s(%d)%s",name,i,extension))
	}
	out, err:=os.Create(outcandidate)
	if err != nil {return err}
	defer out.Close()
	_, err = io.Copy(out,file)
	return err
}
	


