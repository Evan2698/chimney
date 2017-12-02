package sercurity

import (
	"testing"
)

func Test_En(t *testing.T){
	en := MakeCompressKey("Evan2698")
    if (len(en)  != 32 ){
		t.Error("len is ", len(en))
	}
}

func Test_Compress(t *testing.T) {
	
   salt := MakeSalt()
   srcbyte := []byte (("this is a happy day!")[:])
   key := MakeCompressKey("Evan2698-019283873")
   en, _:= Compress(srcbyte, salt,key)
   unen,_:= Uncompress(en, salt, key)
   hello := string(unen[:])

   t.Log("unencrypt content is: ", hello)
   if (hello != "this is a happy day!"){
	   t.Error("not equal!!")
   }
}