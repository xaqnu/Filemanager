package clipboard

import (
	"os"
	"os/exec"
	"testing"
)

func TestPoping_empty_clipboard_fails(t *testing.T){
	var c Clipboard
	_, err:=c.pop()
	if err == nil { t.Errorf("pop didn't fail on empty clipboard")}


}
func TestPushing_empty_clipboard_makes_it_one_longer(t *testing.T){

	var c Clipboard
	var e Entry
	l:=len (c)
	c.push(e)
	if l+1 != len(c) { t.Errorf("one push to empty clipboard doesnt give one element")}
}


func TestPushing_end_then_poping_empty_leaves_it_empty(t *testing.T){

	var c Clipboard
	entry:=Entry{"hi",Cut}
	c.push(entry)
	c.pop()

	if 0 != len(c) { t.Errorf("one push and onbe pop to  empty clipboard doesn't leave empty empty")}
}

func TestPoping_returns_last_element(t *testing.T){

	var c Clipboard
	entry:=Entry{"hi",Cut}
	entry2:=Entry{"last",Copy}
	c.push(entry)
	c.push(entry2)
	result,_:= c.pop()
	
	if result != entry2 { t.Errorf("pop did not return last element")}
}

func TestCopy_handler_loads_data_correctly(t *testing.T){
	var c Clipboard
	c.copyhandler("/hi")
	entry, err := c.pop()
	if err != nil {t.Error("copyhandler did not load any data")}
	expected:=Entry{"/hi",Copy} 
	if entry != expected {t.Errorf("loaded wrong data %+v and expected %+v",entry,expected)}
}

func TestCut_handler_loads_data_correctly(t *testing.T){
	var c Clipboard
	c.cuthandler("/hi")
	entry, err := c.pop()
	if err != nil {t.Error("copyhandler did not load any data")}
	expected:=Entry{"/hi",Cut} 
	if entry != expected {t.Errorf("loaded wrong data %+v and expected %+v",entry,expected)}
}


func TestCopiying_files_creates_perfect_copy(t *testing.T){
	os.MkdirAll("/home/unax/t/text", 0750)
	os.Chdir("/home/unax/t/text")
	f,_:=os.Create("testfile.txt")
	f.WriteString("hola este documento testea el copiar \n segimos un poco")
	f.Close()
	os.MkdirAll("/home/unax/t/copy", 0750)
	err:=filecopy("/home/unax/t/text/testfile.txt","/home/unax/t/copy")
	
	if err != nil {t.Errorf("copying failed %s", err)}
	cmd := exec.Command( "diff", "/home/unax/t/text/testfile.txt", "/home/unax/t/copy/testfile.txt")
	output, err := cmd.CombinedOutput()
    defer os.RemoveAll("/home/unax/t")
	if len(output)!=0 || err != nil  {t.Errorf("copied wrong %s", err)}

}
