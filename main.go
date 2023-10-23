package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/joho/godotenv"
)

var Config = struct {
	banner string

	serverHost    string
	serverPort    int
	serverKeyPath string

	emailBody string
	emailExec string
	emailArgs string
}{
	banner: "\n\nWELCOME TO SSH FORM\n\n",

	serverHost:    "localhost",
	serverPort:    2222,
	serverKeyPath: ".ssh/term_info_ed25519",

	emailBody: "{name} <{email}>\n{content}",
	emailExec: "/usr/sbin/sendmail",
	emailArgs: "",
}

type contact struct {
	term     string
	width    int
	height   int
	form     form
	confetti confetti
	ended    bool
}

func newContact() (*contact, error) {

	form := NewForm()
	confetti := NewConfettiModel()

	return &contact{
		form:     form,
		confetti: confetti,
	}, nil
}

func SendMail(name string, email string, content string) error {

	body := strings.ReplaceAll(Config.emailBody, "{name}", name)
	body = strings.ReplaceAll(body, "{email}", email)
	body = strings.ReplaceAll(body, "{content}", content)

	sendmail := exec.Command(Config.emailExec, strings.Split(Config.emailArgs, " ")...)
	sendmail.Stdin = bytes.NewReader([]byte(body))

	return sendmail.Run()
}

func (e contact) Init() tea.Cmd {
	return nil
}

func (e contact) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FormComplete:
		log.Info("New contact", "name", msg.name, "email", msg.email, "content", msg.content)
		SendMail(msg.name, msg.email, msg.content)
		e.ended = true
		e.confetti.Update(tea.WindowSizeMsg{Width: e.width, Height: e.height})
		return e, animate()
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return e, tea.Quit

		}
	}
	var cmd tea.Cmd
	if e.ended {
		e.confetti, cmd = e.confetti.Update(msg)
	} else {
		e.form, cmd = e.form.Update(msg)
	}
	return e, cmd
}

func (e contact) View() string {

	if e.ended {
		return e.confetti.View()
	}

	return e.form.View()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Getenv("BANNER")) > 0 {
		Config.banner = os.Getenv("BANNER")
	}

	if len(os.Getenv("SERVER_HOST")) > 0 {
		Config.serverHost = os.Getenv("SERVER_HOST")
	}

	if len(os.Getenv("SERVER_PORT")) > 0 {
		port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
		if err == nil {
			Config.serverPort = port
		}
	}

	if len(os.Getenv("SERVER_KEY_PATH")) > 0 {
		Config.serverKeyPath = os.Getenv("SERVER_KEY_PATH")
	}

	if len(os.Getenv("EMAIL_BODY")) > 0 {
		Config.emailBody = os.Getenv("EMAIL_BODY")
	}

	if len(os.Getenv("EMAIL_EXEC")) > 0 {
		Config.emailExec = os.Getenv("EMAIL_EXEC")
	}

	if len(os.Getenv("EMAIL_ARGS")) > 0 {
		Config.emailArgs = os.Getenv("EMAIL_ARGS")
	}

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", Config.serverHost, Config.serverPort)),
		wish.WithHostKeyPath(Config.serverKeyPath),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", Config.serverHost, "port", Config.serverPort, "keypath", Config.serverKeyPath)
	log.Info("with config", "port", Config.serverPort, "emailbody", Config.emailBody, "emailexec", Config.emailExec, "emailargs", Config.emailArgs, "banner", Config.banner)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	m, _ := newContact()

	m.term = pty.Term
	m.width = pty.Window.Width
	m.height = pty.Window.Height

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
