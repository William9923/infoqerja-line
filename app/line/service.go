package line

import (
	iqq "infoqerja-line/app/event/command"
	iqi "infoqerja-line/app/event/input"
	"infoqerja-line/app/utils/constant"
	"log"
	"regexp"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

type (
	// Command : Interface for Reply service
	Command interface {
		GetReply() []linebot.SendingMessage
	}

	// FinderCommand : interface of searching command service
	FinderCommand interface {
		GetCommand() Command
	}

	// Finder : A service for searching something
	Finder struct {
		Command string
	}
)

// IsValidCommand : Function to check wether user inputs is a command or not
func IsValidCommand(message string) bool {
	re := regexp.MustCompile("^!")
	return re.FindString(message) != ""
}

// GetCommand : get the type of command from user inputs
func (finder *Finder) GetCommand() Command {
	co := strings.TrimSpace(finder.Command)
	switch co {
	case constant.HelpCommandCode:
		return &iqq.IncomingHelp{}
	case constant.AddCommandCode:
		return &iqq.IncomingAdd{}
	case constant.ShowCommandCode:
		return &iqq.IncomingShow{}
	case constant.WelcomeCommandCode: // hard coded code, for retrieving the welcome home page
		return &iqq.Welcome{}
	case constant.UnWelcomeCommandCode:
		return &iqq.UnWelcome{}
	default:
		return &iqq.IncomingInvalid{}
	}
}

type (
	// Job : Interface for Executing a job
	Job interface {
		Execute() error
	}

	// FinderJob : interface of searching job service
	FinderJob interface {
		GetJob() Job
	}

	// JobState : struct representing current state of user
	JobState struct {
		State string
	}
)

// GetState : get the type of command from user inputs
func (state *JobState) GetJob() Job {
	switch state.State {
	case constant.WaitDateInput:
		return &iqi.IncomingAddDateJob{}
	case constant.WaitDescInput:
		return &iqi.IncomingAddDescJob{}
	case constant.WaitTitleInput:
		return &iqi.IncomingAddTitleJob{}
	case constant.NoState:
		return &iqi.IncomingStartInput{}
	default:
		return &iqi.IncomingErrorEvent{}
	}
}

// Service : A Handler for containing the necessary dependencies for all messages
type Service struct {
	Bot   BotClient
	Event linebot.Event
}

// MessageService : interface for injecting messaging service
type MessageService interface {
	MessageServiceReply(command Command) error
}

// JobService : interface for injecting job service
type JobService interface {
	JobServiceExecute(job Job) error
}

// MessageServiceReply : Method service for IncomingAction instance; the service that were going to be injected is the Command interface service
func (service *Service) MessageServiceReply(command Command) error {
	// exec methoda
	_, err := service.Bot.ReplyMessage(service.Event.ReplyToken, command.GetReply()...).Do()
	return err
}

// JobServiceExecute : Method service for IncomingJob instance; the service that were going to be injected is the Job interface service
func (service *Service) JobServiceExecute(job Job) error {
	// executing the method
	if err := job.Execute(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// HandleIncomingCommand : Handler for any incoming event that based on EventTypeMessage
func HandleIncomingCommand(service MessageService, finder FinderCommand) {
	command := finder.GetCommand()
	if command != nil {
		if err := service.MessageServiceReply(command); err != nil {
			log.Print(err)
		}
	}
}

func HandleIncomingService(service JobService, finder FinderJob) {
	job := finder.GetJob()
	// filling job description data
	if job != nil {
		if err := service.JobServiceExecute(job); err != nil {
			log.Print(err)
			finderLocal := &JobState{
				State: "error",
			}
			errJob := finderLocal.GetJob() // handling error
			_ = service.JobServiceExecute(errJob)
		}
	}

}
