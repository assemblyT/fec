package fec

import (
	"container/list"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

// http://www.fec.gov/finance/disclosure/metadata/DataDictionaryOperatingExpenditures.shtml
type Oppexp struct {
	CMTE_ID          string // NOTNULL
	AMNDT_IND        byte
	RPT_YR           uint16
	RPT_TP           string
	IMAGE_NUM        string
	LINE_NUM         string
	FORM_TP_CD       string
	SCHED_TP_CD      string
	NAME             string
	CITY             string
	STATE            string
	ZIP_CODE         string
	TRANSACTION_DT   time.Time
	TRANSACTION_AMT  float64
	TRANSACTION_PGI  string
	PURPOSE          string
	CATEGORY         string
	CATEGORY_DESC    string
	MEMO_CD          byte
	MEMO_TEXT        string
	ENTITY_TP        string
	SUB_ID           uint64 // NOTNULL
	FILE_NUM         uint64
	TRAN_ID          string
	BACK_REF_TRAN_ID string
}

func (t Oppexp) Equals(o Oppexp) bool {
	return reflect.DeepEqual(t, o)
}

func (t Oppexp) IsEmpty() bool {
	var nilOppexp Oppexp
	return t.Equals(nilOppexp)
}

func DefaultLocation() *time.Location {
	dc_tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Sprintf("Could not get America/New_York location"))
	}
	return dc_tz
}

func LoadOppexpDatabase(file io.Reader) (*list.List, error) {
	l := list.New()

	rdr := csv.NewReader(file)
	rdr.Comma = '|'
	rdr.TrailingComma = true

	var line []string
	var err error

	for line, err = rdr.Read(); err != io.EOF; line, err = rdr.Read() {
		if err != nil {
			return l, err
		}

		o := new(Oppexp)
		FromStringSlice(o, line)
		l.PushBack(o)
	}

	return l, nil
}

func FromStringSlice(o *Oppexp, s []string) error {
	var amndt_ind, memo_cd byte

	if len(s) < 25 {
		msg := fmt.Sprintf("String array needs to be at least 25 long, but was %d", len(s))
		return errors.New(msg)
	}

	if len(s[1]) < 1 {
		amndt_ind = 0
	} else {
		amndt_ind = s[1][0]
	}

	if len(s[18]) != 1 {
		memo_cd = 0
	} else {
		memo_cd = s[1][0]
	}

	rpt_yr, err := strconv.ParseUint(s[2], 10, 16)
	if err != nil {
		return err
	}

	transaction_amt, err := strconv.ParseFloat(s[13], 64)
	if err != nil {
		return err
	}

	sub_id, err := strconv.ParseUint(s[21], 10, 64)
	if err != nil {
		return err
	}

	file_num, err := strconv.ParseUint(s[22], 10, 64)
	if err != nil {
		return err
	}

	transaction_dt, err := time.ParseInLocation("01/02/2006", s[12], DefaultLocation())
	if err != nil {
		return err
	}

	o.CMTE_ID = s[0]
	o.AMNDT_IND = amndt_ind
	o.RPT_YR = uint16(rpt_yr)
	o.RPT_TP = s[3]
	o.IMAGE_NUM = s[4]
	o.LINE_NUM = s[5]
	o.FORM_TP_CD = s[6]
	o.SCHED_TP_CD = s[7]
	o.NAME = s[8]
	o.CITY = s[9]
	o.STATE = s[10]
	o.ZIP_CODE = s[11]
	o.TRANSACTION_DT = transaction_dt
	o.TRANSACTION_AMT = transaction_amt
	o.TRANSACTION_PGI = s[14]
	o.PURPOSE = s[15]
	o.CATEGORY = s[16]
	o.CATEGORY_DESC = s[17]
	o.MEMO_CD = memo_cd
	o.MEMO_TEXT = s[19]
	o.ENTITY_TP = s[20]
	o.SUB_ID = sub_id
	o.FILE_NUM = file_num
	o.TRAN_ID = s[23]
	o.BACK_REF_TRAN_ID = s[24]

	return nil
}
