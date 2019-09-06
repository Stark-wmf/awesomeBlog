package services

import (
	"bytes"
	"crypto/rand"
	"math/big"
	//"strings"
	"gopkg.in/gomail.v2"
	//"time"
)
//
//const XXX_MAIL_TEMPLATE = `
//    <div>
//        <h3>123</h3>
//        <p>3456</p>
//        <h3>789</h2>
//        <table style="border-collapse:collapse;border: 1px solid black;">
//            <thead style="border-collapse:collapse;border: 1px solid black;">
//                <tr style="border-collapse:collapse;border: 1px solid black;text-align: center;">
//                    <th style="border-collapse:collapse;border: 1px solid black;">Case Name</th>
//                    <th style="border-collapse:collapse;border: 1px solid black;">Owner</th>
//                    <th style="border-collapse:collapse;border: 1px solid black;">Creator</th>
//                    <th style="border-collapse:collapse;border: 1px solid black;">Status</th>
//                </tr>
//            </thead>
//            <tbody>
//            {{with .Job}}
//{{range $k, $v := .Cases}}
//                <tr style="border-collapse:collapse;border: 1px solid black;text-align: center;">
//                    <td style="border-collapse:collapse;border: 1px solid black;">{{$v.Name}}         </td>
//                    <td style="border-collapse:collapse;border: 1px solid black;">{{$v.IsSuccess}}          </td>
//                    <td style="border-collapse:collapse;border: 1px solid black;">{{$v.Agent}}               </td>
//
//                </tr>
//             {{end}}
//             {{end}}
//            </tbody>
//        </table>
//
//    </div>`
//
//type PageInfo struct {
//	job Job
//}
//
//type Job struct {
//	email string
//}
func CreateRandomString(len int) string  {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < len ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func Sendemail(email string) string {
	s:=CreateRandomString(6)
	//MAIL_TEMPLATE := XXX_MAIL_TEMPLATE
	m := gomail.NewMessage()
	m.SetHeader("From", "675075368@qq.com")
	m.SetHeader("To", email)//send email to multipul persons
	m.SetHeader("Subject", s)
	m.SetBody("text/plain", s)
	buffer := new(bytes.Buffer)
	m.SetBody("text/html", buffer.String())
	d := gomail.Dialer{Host: "smtp.qq.com", Port: 465, Username: "675075368@qq.com", Password: "oponzilcccybbbjf",SSL:true}
	if err := d.DialAndSend(m); err != nil {
		return err.Error()
	}
	return s
}
