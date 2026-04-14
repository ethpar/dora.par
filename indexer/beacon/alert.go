package beacon

import (
	"bytes"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/dora/db"
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/ethpandaops/dora/utils"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

type alertsSender struct {
	indexer *Indexer
}

func newAlertsSender(indexer *Indexer) *alertsSender {
	emailConfig := utils.Config.Email
	indexer.logger.Infof("emailConfig: %v", emailConfig)
	/*return &alertsSender{
		indexer: indexer,
	}*/
	d := &alertsSender{
		indexer: indexer,
	}
	//d1 := 3.344565656565
	//d.sendEmail(uint64(1), uint64(98), float64(d1))
	return d
}

func (al *alertsSender) checkAndSendAlert(tx *sqlx.Tx, epoch phase0.Epoch, blocks []*Block, epochStats *EpochStats, epochVotes *EpochVotes) {

	epochAlert := dbtypes.EpochAlertState{}
	db.GetExplorerState("alert.epoch", &epochAlert)
	al.indexer.logger.Infof("epochAlert: %v", epochAlert)
	al.indexer.logger.Infof("epochAlert: epoch: %v vote: %v", epoch, epochVotes.TargetVotePercent)
	epochAlert.Epoch = uint64(epoch)
	var vote = 100
	if epochVotes.TargetVotePercent <= 90 {
		vote = 90
	} else if epochVotes.TargetVotePercent <= 95 {
		vote = 95
	} else if epochVotes.TargetVotePercent <= 98 {
		vote = 98
	}
	if vote != epochAlert.Percent {
		epochAlert.Percent = vote
		db.SetExplorerState("alert.epoch", epochAlert, tx)
		if vote != 100 {
			al.sendEmail(uint64(epoch), uint64(vote), epochVotes.TargetVotePercent)
			al.sendSlack(uint64(epoch), uint64(vote), epochVotes.TargetVotePercent)
		}
	}
}
func (al *alertsSender) sendSlack(epoch uint64, vote uint64, targetVotePercent float64) {
	if !utils.Config.Slack.Enabled {
		return
	}
	doraUrl := utils.Config.Slack.DoraUrl
	votedString := strconv.FormatFloat(targetVotePercent, 'f', 2, 64)
	subject := utils.Config.Slack.Subject

	circle := ":red_circle:"
	if vote == 95 {
		circle = ":large_orange_circle:"
	}
	if vote == 98 {
		circle = ":large_green_circle:"
	}
	body := circle + "  " + subject + "  Epoch " + strconv.FormatUint(epoch, 10) + " voted " + votedString + "%  " +
		doraUrl + "/epoch/" + strconv.FormatUint(epoch, 10)

	url := utils.Config.Slack.SlackUrl
	data := map[string]string{"text": body}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		al.indexer.logger.Fatal(err)
	}
	defer resp.Body.Close()
}

func (al *alertsSender) sendEmail(epoch uint64, vote uint64, targetVotePercent float64) {
	if !utils.Config.Email.Enabled {
		return
	}

	from := utils.Config.Email.From

	addrList, err := mail.ParseAddressList(utils.Config.Email.To)
	if err != nil {
		al.indexer.logger.Fatal(err)
	}

	var emailsOnly []string
	var formattedNames []string

	for _, addr := range addrList {
		emailsOnly = append(emailsOnly, addr.Address)
		formattedNames = append(formattedNames, addr.String())
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	password := utils.Config.Email.Password

	doraUrl := utils.Config.Email.DoraUrl

	toHeader := "To: " + strings.Join(formattedNames, ", ") + "\r\n"
	subject := utils.Config.Email.Subject
	votedString := strconv.FormatFloat(targetVotePercent, 'f', 2, 64)
	percentLevel := strconv.FormatUint(vote, 10)
	subject = subject + " Epoch " + strconv.FormatUint(epoch, 10) + " voted " + votedString + "% (less than " + percentLevel + "%)\n"
	body := "Epoch " + strconv.FormatUint(epoch, 10) + " voted " + votedString + "%  " +
		doraUrl + "/epoch/" + strconv.FormatUint(epoch, 10) +
		" .\n"
	msg := []byte(toHeader + subject + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, emailsOnly, msg)
	if err != nil {
		al.indexer.logger.Fatal(err)
	}

	al.indexer.logger.Println("Email sent successfully!")
}
