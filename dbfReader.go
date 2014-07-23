package dbfReader

import(
    "os"
    "fmt"
    "strings"
)

type DbfHead struct {
    Version    []byte
    Updatedate string
    Records    int64
    Headerlen  int64
    Recordlen  int64
}

type Field struct {
    Name             string
    Fieldtype        string
    FieldDataaddress []byte
    FieldLen         int64
    DecimalCount     []byte
    Workareaid       []byte
}

type Record struct {
    Delete bool
    Data map[string]string
}

func GetDbfHead(reader *os.File) (dbfhead DbfHead) {
    buf := make([]byte, 16)
    reader.Seek(0, 0)
    _, err := reader.Read(buf)
    if err != nil {
        panic(err)
    }
    dbfhead.Version = buf[0:1]
    dbfhead.Updatedate = fmt.Sprintf("%d", buf[1:4])
    dbfhead.Headerlen = Changebytetoint(buf[8:10])
    dbfhead.Recordlen = Changebytetoint(buf[10:12])
    dbfhead.Records = Changebytetoint(buf[4:8])
    return dbfhead
}

func Changebytetoint(b []byte) (x int64) {
    for i, val := range b {
        if i == 0 {
            x = x + int64(val)
        } else {
            x = x + int64(2<<7*int64(i)*int64(val))
        }
    }

    return
}


func RemoveNullfrombyte(b []byte) (s string) {
    for _, val := range b {
        if val == 0 {
            continue
        }
        s = s + string(val)
    }
    return
}

func GetFields(reader *os.File) []Field {
    dbfhead := GetDbfHead(reader)

//    off := dbfhead.Headerlen - 32 - 264
    off := dbfhead.Headerlen - 32 - 1
    fieldlist := make([]Field, off/32)
    buf := make([]byte, off)
    _, err := reader.ReadAt(buf, 32)
    if err != nil {
        panic(err)
    }
    curbuf := make([]byte, 32)
    for i, val := range fieldlist {
        a := i * 32
        curbuf = buf[a:]
        val.Name = RemoveNullfrombyte(curbuf[0:11])
        val.Fieldtype = fmt.Sprintf("%s", curbuf[11:12])
        val.FieldDataaddress = curbuf[12:16]
        val.FieldLen = Changebytetoint(curbuf[16:17])
        val.DecimalCount = curbuf[17:18]
        val.Workareaid = curbuf[20:21]
        fieldlist[i] = val

    }
    return fieldlist
}

func GetRecords(fp *os.File) (records []map[string]string) {
    dbfhead := GetDbfHead(fp)
    fp.Seek(0, 0)
    fields := GetFields(fp)
    recordlen := dbfhead.Recordlen
    start := dbfhead.Headerlen
    buf := make([]byte, recordlen)

    temp := make([]map[string]string, dbfhead.Records)
    for {
        _, err := fp.ReadAt(buf, start)
        if err != nil {
            return temp
            panic(err)
        }
        record := map[string]string{}
        // * 删除
        if string(buf[0:1]) == " " {
            record["delete"] = "0"
        } else if string(buf[0:1]) == "*" {
            record["delete"] = "1"
        }
        a := int64(1)
        for _, val := range fields {
            fieldlen := val.FieldLen
            record[val.Name] = strings.Trim(fmt.Sprintf("%s", buf[a:a+fieldlen]), " ")
            a = a + fieldlen
        }
        temp = append(temp, record)
        start = start + recordlen
    }
}

/*
func GetRecordbyField(fieldname string, fieldval string, fp *os.File) (record map[int]Record) {
    //GetDbfHead(fp)
    fields := GetFields(fp)
    records := GetRecords(fp)
    temp := map[int]Record{}
    i := 0
    for _, val := range records {
        for _, val1 := range fields {
            if val1.Name == fieldname && val["delete"] == "1" {
                if val[val1.Name] == fieldval || val[val1.Name] == " " {
                    temp[i] = val
                }

            }
        }
        i = i + 1
    }
    return temp
}
*/
