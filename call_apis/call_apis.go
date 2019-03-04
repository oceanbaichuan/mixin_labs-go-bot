package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"fmt"
	"bufio"
	"io"
	mixin "github.com/MooooonStar/mixin-sdk-go/network"
)

const (
	UserId          = "21042518-85c7-4903-bb19-f311813d1f51"
	PinCode         = "911424"
	SessionId       = "4267b63d-3daa-449e-bc13-970aa0357776"
	PinToken        = "gUUxpm3fPRVkKZNwA/gk10SHHDtR8LmxO+N6KbsZ/jymmwwVitUHKgLbk1NISdN8jBvsYJgF/5hbkxNnCJER5XAZ0Y35gsAxBOgcFN8otsV6F0FAm5TnWN8YYCqeFnXYJnqmI30IXJTAgMhliLj7iZsvyY/3htaHUUuN5pQ5F5s="
	//please delele the blank of PrivateKey the before each line
	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCDXiWJRLe9BzPtXmcVe6acaFTY9Ogb4Hc2VHFjKFsp7QRVCytx
3KC/LRojTFViwwExaANTZQ6ectwpAxIvzeYeHDZCXCh6JRFIYK/ZuREmYPcPQEWD
s92Tv/4XTAdTH8l9UJ4VQY4zwqYMak237N9xEvowT0eR8lpeJG0jAjN97QIDAQAB
AoGADvORLB1hGCeQtmxvKRfIr7aEKak+HaYfi1RzD0kRjyUFwDQkPrJQrVGRzwCq
GzJ8mUXwUvaGgmwqOJS75ir2DL8KPz7UfgQnSsHDUwKqUzULgW6nd/3OdDTYWWaN
cDjbkEpsVchOpcdkywvZhhyGXszpM20Vr8emlBcFUOTfpTUCQQDVVjkeMcpRsImV
U3tPYyiuqADhBTcgPBb+Ownk/87jyKF1CZOPvJAebNmpfJP0RMxUVvT4B9/U/yxZ
WNLhLtCXAkEAnaOEuefUxGdE8/55dUTEb7xrr22mNqykJaax3zFK+hSFBrM3gUY5
fEETtHnl4gEdX4jCPybRVc1JSFY/GWoyGwJBAKoLti95JHkErEXYavuWYEEHLNwv
mgcZnoI6cOKVfEVYEEoHvhTeCkoWHVDZOd2EURIQ1eY18JYIZ0M4Z66R8DUCQCsK
iKTR3dA6eiM8qiEQw6nWgniFscpf3PnCx/Iu3U/m5mNr743GhM+eXSj7136b209I
YfEoQiPxRz8O/W+NBV0CQQDVPNqJlFD34MC9aQN42l3NV1hDsl1+nSkWkXSyhhNR
MpobtV1a7IgJGyt5HxBzgNlBNOayICRf0rRjvCdw6aTP
-----END RSA PRIVATE KEY-----`
	MASTER_UUID     = "0b4f49dc-8fb4-4539-9a89-fb3afc613747"
	BTC_WALLET_ADDR = "14T129GTbXXPGXXvZzVaNLRFPeHXD1C25C"
	AMOUNT          = "0.001"
)

func main() {
	scanner   := bufio.NewScanner(os.Stdin)
	var PromptMsg string
	PromptMsg  = "1: Create user and update PIN\n2: Read Bitcoin balance \n3: Read Bitcoin Address\n4: Read EOS balance\n"
	PromptMsg += "5: Read EOS address\n6: Transfer Bitcoin from bot to new account\n7: Transfer Bitcoin from new account to Master\n"
	PromptMsg += "8: Withdraw bot's Bitcoin\na: Verify Pin\nd: Create Address and Delete it\nr: Create Address and read it\n"
	PromptMsg += "q: Exit \nMake your choose:"
	for  {
	 fmt.Print(PromptMsg)
	 var cmd string
	 scanner.Scan()
	 cmd = scanner.Text()
	 if cmd == "q" {
			 break
	 }
	 if cmd == "1" {
		 fo, err := os.OpenFile("new_users.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		 if err != nil {
			 panic(err)
			 return
		 }

		 user,err := mixin.CreateAppUser("Tom cat", PinCode, UserId, SessionId, PrivateKey)
		 if err != nil {
				 panic(err)
		 }
		 records := [][]string {
												 {user.UserId,user.PrivateKey,user.SessionId,user.PinToken,user.PinCode},
												 }
		 w := csv.NewWriter(fo)
		 if err := w.Error(); err != nil {
			 log.Fatalln("error writing csv:", err)
		 }
		 w.WriteAll(records) // calls Flush internally
		 fo.Close()
		 log.Println(user)
	 }
	 if cmd == "2"|| cmd == "3" {
		 csvFile, err := os.Open("new_users.csv")
		 if err != nil {
            log.Fatal(err)
      }
		 reader := csv.NewReader(bufio.NewReader(csvFile))
		 for {
				 record, err := reader.Read()
				 if err == io.EOF {
					 break
				 }
				 if err != nil {
					 log.Fatal(err)
				 }
				 UserID             := record[0]
				 PrivateKey         := record[1]
				 SessionID     		  := record[2]
				 UserInfoBytes, err := mixin.ReadAsset(mixin.GetAssetId("BTC"),UserID,SessionID,PrivateKey)
				 if err != nil {
								 log.Fatal(err)
				 }
				 var UserInfoMap map[string]interface{}
				 if err := json.Unmarshal(UserInfoBytes, &UserInfoMap); err != nil {
						 panic(err)
				 }
				 fmt.Println("User ID ",UserID, "'s Bitcoin Address is: ",UserInfoMap["data"].(map[string]interface{})["public_key"])
				 fmt.Println("Balance is: ",UserInfoMap["data"].(map[string]interface{})["balance"])
			 }
		 csvFile.Close()
	 }
	 if cmd == "4"|| cmd == "5" {
		 csvFile, err := os.Open("new_users.csv")
		 if err != nil {
            log.Fatal(err)
      }
		 reader := csv.NewReader(bufio.NewReader(csvFile))
		 for {
				 record, err := reader.Read()
				 if err == io.EOF {
					 break
				 }
				 if err != nil {
					 log.Fatal(err)
				 }
				 UserID             := record[0]
				 PrivateKey         := record[1]
				 SessionID     		  := record[2]
				 UserInfoBytes, err := mixin.ReadAsset(mixin.GetAssetId("EOS"),UserID,SessionID,PrivateKey)
				 if err != nil {
								 log.Fatal(err)
				 }
				 var UserInfoMap map[string]interface{}
				 if err := json.Unmarshal(UserInfoBytes, &UserInfoMap); err != nil {
						 panic(err)
				 }
				 fmt.Println("User ID ",UserID, "'s EOS Address is: ",
					 UserInfoMap["data"].(map[string]interface{})["account_name"],
				   UserInfoMap["data"].(map[string]interface{})["account_tag"])
				 fmt.Println("Balance is: ",UserInfoMap["data"].(map[string]interface{})["balance"])
			 }
		 csvFile.Close()
	 }
 }
}
