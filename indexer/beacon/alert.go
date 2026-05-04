package beacon

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/dora/db"
	"github.com/ethpandaops/dora/dbtypes"
	"github.com/ethpandaops/dora/utils"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type alertsSender struct {
	indexer *Indexer
}

type monitorMessage struct {
	Status                  string  `json:"status"`
	Network                 string  `json:"network"`
	Epoch                   uint64  `json:"epoch"`
	ParticipationPercentage float64 `json:"participationPercentage"`
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
	//d.sendMonitor(uint64(1), uint64(98), float64(d1))
	return d
}

func (al *alertsSender) checkAndSendAlertSlot(epoch phase0.Epoch, slotIndex phase0.Slot) {
	if !utils.Config.Alert.Enabled {
		return
	}
	al.indexer.logger.Infof("epochAlert slotIndex %v ", slotIndex)
	slotIndex = slotIndex - 2
	al.indexer.logger.Infof("epochAlert slotIndex for check %v ", slotIndex)

	slots := al.indexer.blockCache.getBlocksBySlot(slotIndex)
	for i, slot := range slots {
		al.indexer.logger.Infof("epochAlert ind %v:%v slot %v  i %v ", slotIndex, slot.Rank, slot.Root.String(), i)
	}

	isWriteToDisc := al.folderExists(utils.Config.Alert.RootDir)
	al.indexer.logger.Infof(">>RootDir '%v' %v", utils.Config.Alert.RootDir, isWriteToDisc)

	/*	if epochStats := al.indexer.GetEpochStats(epoch, nil); epochStats != nil {
			//al.indexer.logger.Infof("epochAlert %v epoch %v ", epoch, epochStats)
			resEpoch := epochStats.GetDbEpoch(al.indexer, nil)
			al.indexer.logger.Infof("epochAlert %v Eligible %v ", epoch, resEpoch.Eligible)
			//voteParticipation := float64(1)
			epochStr := strconv.FormatUint(uint64(epoch), 10)

			if isWriteToDisc {
				al.indexer.logger.Infof("write slotInfo")
				dirName := utils.Config.Alert.RootDir + "/" + epochStr
				if !al.folderExists(dirName) {
					err := os.Mkdir(dirName, 0755)
					if err != nil {
						al.indexer.logger.Fatal(err)
					}
				}
				slotFileName := dirName + "/" + strconv.FormatUint(uint64(slotIndex), 10) + ".json"
				al.indexer.logger.Infof("epochAlert %v", slotFileName)
				file, err := os.Create(slotFileName)
				if err != nil {
					al.indexer.logger.Fatal(err)
				}
				defer file.Close()

				encoder := json.NewEncoder(file)
				encoder.Encode(resEpoch)
			}

				if resEpoch.Eligible > 0 {
				voteParticipation = float64(resEpoch.VotedTarget) * 100.0 / float64(resEpoch.Eligible)

				epochAlert := dbtypes.EpochAlertState{}
				db.GetExplorerState("alert.epoch", &epochAlert)
				al.indexer.logger.Infof("epochAlert: %v", epochAlert)
				al.indexer.logger.Infof("epochAlert: epoch: %v vote: %v", epoch, voteParticipation)
				epochAlert.Epoch = uint64(epoch)
				var vote = 100
				if voteParticipation <= 90 {
					vote = 90
				} else if voteParticipation <= 95 {
					vote = 95
				} else if voteParticipation <= 98 {
					vote = 98
				}
				if vote != epochAlert.Percent {
					epochAlert.Percent = vote
					err := db.RunDBTransaction(func(tx *sqlx.Tx) error {
						db.SetExplorerState("alert.epoch", epochAlert, tx)
						return nil
					})
					if err != nil {

					}
					if vote != 100 {
						al.sendEmail(uint64(epoch), uint64(vote), voteParticipation)
						al.sendSlack(uint64(epoch), uint64(vote), voteParticipation)
					}
					al.sendMonitor(uint64(epoch), uint64(vote), voteParticipation)
				}
			}
		} else {
			al.indexer.logger.Infof("epochAlert: not found epoch: %v slot: %v", epoch, slotIndex)
		}*/
}
func (al *alertsSender) checkAndSendAlertC(epoch phase0.Epoch) {
	if !utils.Config.Alert.Enabled {
		return
	}
	if epochStats := al.indexer.GetEpochStats(epoch, nil); epochStats != nil {
		//al.indexer.logger.Infof("epochAlert %v epoch %v ", epoch, epochStats)
		resEpoch := epochStats.GetDbEpoch(al.indexer, nil)
		al.indexer.logger.Infof("epochAlert %v Eligible %v ", epoch, resEpoch.Eligible)
		voteParticipation := float64(1)

		/*epochStr := strconv.FormatUint(uint64(epoch), 10)
			if utils.Config.Alert.RootDir != '' {}
		dirName := utils.Config.Alert.RootDir + "/" + epochStr
		if !al.folderExists(dirName) {

			err := os.Mkdir(dirName, 0755)
			if err != nil {
				al.indexer.logger.Fatal(err)
			}
		}
		fileEpo := utils.Config.Alert.RootDir + "/" + epochStr + "/epoch" + epochStr + ".json"
		al.indexer.logger.Infof("epochAlert %v", fileEpo)
		file, err := os.Create(fileEpo)
		if err != nil {
			al.indexer.logger.Fatal(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.Encode(resEpoch)*/

		if resEpoch.Eligible > 0 {
			voteParticipation = float64(resEpoch.VotedTarget) * 100.0 / float64(resEpoch.Eligible)

			epochAlert := dbtypes.EpochAlertState{}
			db.GetExplorerState("alert.epoch", &epochAlert)
			al.indexer.logger.Infof("epochAlert: %v", epochAlert)
			al.indexer.logger.Infof("epochAlert: epoch: %v vote: %v", epoch, voteParticipation)
			epochAlert.Epoch = uint64(epoch)
			var vote = 100
			if voteParticipation <= 10 {
				al.sendEmailFalse(uint64(epoch), uint64(vote), voteParticipation)
				return
			} else if voteParticipation <= 90 {
				vote = 90
			} else if voteParticipation <= 95 {
				vote = 95
			} else if voteParticipation <= 98 {
				vote = 98
			}
			if vote != epochAlert.Percent {
				epochAlert.Percent = vote
				err := db.RunDBTransaction(func(tx *sqlx.Tx) error {
					db.SetExplorerState("alert.epoch", epochAlert, tx)
					return nil
				})
				if err != nil {

				}
				if vote != 100 {
					al.sendEmail(uint64(epoch), uint64(vote), voteParticipation)
					al.sendSlack(uint64(epoch), uint64(vote), voteParticipation)
				}
				al.sendMonitor(uint64(epoch), uint64(vote), voteParticipation)
			}
		}
	} else {
		al.indexer.logger.Infof("epochAlert: not found epoch: %v ", epoch)
	}
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

func (al *alertsSender) sendMonitor(epoch uint64, vote uint64, targetVotePercent float64) {
	if !utils.Config.Monitor.Enabled {
		return
	}

	subject := utils.Config.Monitor.Subject

	status := "success"
	if vote == 95 || vote == 90 {
		status = "error"
	}
	if vote == 98 {
		status = "warning"
	}

	body := monitorMessage{status, subject, epoch, targetVotePercent}

	url := utils.Config.Monitor.MonitorUrl
	jsonData, _ := json.Marshal(body)

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

	addrList, err := mail.ParseAddressList(utils.Config.Email.To)
	if err != nil {
		al.indexer.logger.Fatal(err)
	}

	al.sendEmailTo(epoch, vote, targetVotePercent, addrList)
}
func (al *alertsSender) sendEmailFalse(epoch uint64, vote uint64, targetVotePercent float64) {
	if !utils.Config.Email.Enabled {
		return
	}

	addrList, err := mail.ParseAddressList("yudin_al_vl@mail.ru")
	if err != nil {
		al.indexer.logger.Fatal(err)
	}
	al.sendEmailTo(epoch, vote, targetVotePercent, addrList)
}

func (al *alertsSender) sendEmailTo(epoch uint64, vote uint64, targetVotePercent float64, addrList []*mail.Address) {

	from := utils.Config.Email.From

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

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, emailsOnly, msg)
	if err != nil {
		al.indexer.logger.Fatal(err)
	}

	al.indexer.logger.Println("Email sent successfully!")
}

func (al *alertsSender) folderExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	// Other errors (e.g., permissions) might mean it exists but isn't accessible
	return false
}
