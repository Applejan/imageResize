type sum struct {
sum int
}


func resize (src string,status chan int){
f,err:=imaging.Open(src)
if err!=nil{
log.Fatal("File open fail,quiting!")
}
//f=imaging.Resize()
status<-1
}




fun main (){
//Init start flag
  cpus:=4
  status:=make(chan int, cpus)
  for i:=0;i<cpus;i++{
    status<-1
  }
  //To get filw list
  var files []string
for i,v :=range files {
  distribut (v,a,status)
}
  
}

func distribut(file string,id *sum,status chan int){
<-status
  sum.sum+=status
  fmt.Println("Current file ID is ",sum.sum)
  go resize(file,status)
} 
