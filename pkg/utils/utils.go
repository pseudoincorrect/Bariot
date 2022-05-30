package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mainflux/senml"
)

// PrettyJsonsString returns a string with pretty JSON representation
func PrettyJsonString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "  "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// PrettyJsons returns a string with pretty JSON representation
func PrettySenmlRecord(r senml.Record) string {
	xmlName, value, stringValue, dataValue, boolValue, sum := "", "", "", "", "", ""

	if r.XMLName != nil {
		xmlName = strconv.FormatBool(*r.XMLName)
	}
	if r.Value != nil {
		value = fmt.Sprintf("%f", *r.Value)
	}
	if r.StringValue != nil {
		stringValue = *r.StringValue
	}
	if r.DataValue != nil {
		dataValue = *r.DataValue
	}
	if r.BoolValue != nil {
		boolValue = strconv.FormatBool(*r.BoolValue)
	}
	if r.Sum != nil {
		sum = fmt.Sprintf("%f", *r.Sum)
	}

	baseTime := time.Unix(int64(r.BaseTime), 0).Format(time.RFC3339)
	theTime := time.Unix(int64(r.Time), 0).Format(time.RFC3339)
	updateTime := time.Unix(int64(r.UpdateTime), 0).Format(time.RFC3339)

	return "SenML record:" + "\n" +
		"  XMLName     = " + xmlName + "\n" +
		"  Link        = " + r.Link + "\n" +
		"  BaseName    = " + r.BaseName + "\n" +
		"  BaseTime    = " + baseTime + "\n" +
		"  BaseUnit    = " + r.BaseUnit + "\n" +
		"  BaseVersion = " + strconv.Itoa(int(r.BaseVersion)) + "\n" +
		"  BaseValue   = " + fmt.Sprintf("%f", r.BaseValue) + "\n" +
		"  BaseSum     = " + fmt.Sprintf("%f", r.BaseSum) + "\n" +
		"  Name        = " + r.Name + "\n" +
		"  Unit        = " + r.Unit + "\n" +
		"  Time        = " + theTime + "\n" +
		"  UpdateTime  = " + updateTime + "\n" +
		"  Value       = " + value + "\n" +
		"  StringValue = " + stringValue + "\n" +
		"  DataValue   = " + dataValue + "\n" +
		"  BoolValue   = " + boolValue + "\n" +
		"  Sum         = " + sum + "\n"
}
