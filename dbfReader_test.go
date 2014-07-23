package dbfReader

import(
    "testing"
    "os"
    "fmt"
)

func Test_readdbf(t *testing.T){
    pwd, err := os.Getwd()
    if err != nil {
        t.Errorf("pwd failed", err.Error())
    }

    fmt.Println(pwd)

//    infile := fmt.Sprintf("%s/sjshq.dbf", pwd)
    infile := fmt.Sprintf("%s/show2003.dbf", pwd)

    fp, err := os.Open(infile)
//    fp, err := os.OpenFile(infile, os.O_RDONLY, 0)
    if err != nil {
        t.Errorf("failed to open test file ", infile)
    }

    defer fp.Close()

    head := GetDbfHead(fp)
    fmt.Println("version", head.Version)
    fmt.Println("updatedate", head.Updatedate)
    fmt.Println("records", head.Records)
    fmt.Println("headlen", head.Headerlen)
    fmt.Println("records", head.Recordlen)


    // 取字段名
    fields := GetFields(fp)
    fieldLen := len(fields)
    fmt.Println("filed行数",fieldLen)
    for i := 0; i < fieldLen; i++ {
        fmt.Println(fields[i].Name)
    }

    records := GetRecords(fp)
    recordLen := len(records)
    fmt.Println("记录行数", recordLen)

    for i := 0; i < recordLen; i++ {
        r := records[i]
        if r.Data["S1"] == "122976" ||  r.Data["S1"] == "600644" {
          fmt.Println(r)
        }

    }

    filedRecord := GetRecordbyField("S1", "000001", fp)
    fmt.Println(filedRecord)
}
