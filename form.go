package main

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormComplete struct {
	name    string
	email   string
	content string
}

type form struct {
	width        int
	height       int
	focused      int
	nameInput    textinput.Model
	emailInput   textinput.Model
	contentText  textarea.Model
	captchaInput textinput.Model
	captcha      captcha
	invalid      [4]bool
}

const captchaLen = 4
const sendIndex = 4
const responsiveBreakpoint = 64

func NewForm() form {
	f := form{}

	f.captcha = NewCaptcha(captchaLen)

	name := textinput.New()
	name.Prompt = ""
	name.Width = 30
	name.Focus()
	f.focused = 0
	f.nameInput = name

	email := textinput.New()
	email.Prompt = ""
	email.Width = 30
	f.emailInput = email

	content := textarea.New()
	content.Prompt = ""
	content.SetWidth(f.width - 4)
	content.ShowLineNumbers = false
	f.contentText = content

	captcha := textinput.New()
	captcha.Prompt = ""
	captcha.Width = captchaLen
	captcha.CharLimit = captchaLen
	f.captchaInput = captcha

	return f
}

func (f *form) FocusPrev() tea.Cmd {
	var cmd tea.Cmd

	switch f.focused {
	case 0:
		f.nameInput.Blur()
	case 1:
		f.emailInput.Blur()
		cmd = f.nameInput.Focus()
	case 2:
		f.contentText.Blur()
		cmd = f.emailInput.Focus()
	case 3:
		f.captchaInput.Blur()
		f.contentText.Focus()
	case 4:
		cmd = f.captchaInput.Focus()
	}

	f.focused--
	if f.focused < 0 {
		f.focused = 4
	}

	return cmd
}

func (f *form) FocusNext() tea.Cmd {
	var cmd tea.Cmd

	switch f.focused {
	case 0:
		f.nameInput.Blur()
		cmd = f.emailInput.Focus()
	case 1:
		f.emailInput.Blur()
		cmd = f.contentText.Focus()
	case 2:
		f.contentText.Blur()
		cmd = f.captchaInput.Focus()
	case 3:
		f.captchaInput.Blur()
	case 4:
		cmd = f.nameInput.Focus()
	}

	f.focused++
	if f.focused > 4 {
		f.focused = 0
	}

	return cmd
}

func (f *form) Validate() ([4]string, bool) {
	isValid := true
	index := 0
	var values [4]string

	values[index] = f.nameInput.Value()
	if len(values[index]) < 2 {
		f.invalid[index] = true
		isValid = false
	} else {
		f.invalid[index] = false
	}
	index++

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	values[index] = f.emailInput.Value()
	if len(values[index]) > 0 && !emailRegex.MatchString(values[index]) {
		f.invalid[index] = true
		isValid = false
	} else {
		f.invalid[index] = false
	}
	index++

	values[index] = f.contentText.Value()
	if len(values[index]) < 2 {
		f.invalid[index] = true
		isValid = false
	} else {
		f.invalid[index] = false
	}
	index++

	values[index] = f.captchaInput.Value()
	if len(values[index]) != 4 || !f.captcha.IsValid(values[index]) {
		f.invalid[index] = true
		isValid = false
	} else {
		f.invalid[index] = false
	}
	index++

	return values, isValid
}

func (f form) Update(msg tea.Msg) (form, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height

		if f.width > responsiveBreakpoint {
			remainder := f.width - (f.width/2)*2
			f.nameInput.Width = f.width/2 - 4 + remainder
			f.emailInput.Width = f.width/2 - 4
		} else {
			f.nameInput.Width = f.width - 5
			f.emailInput.Width = f.width - 5
		}

		f.contentText.SetWidth(f.width - 4)
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab", "up", "left":
			cmd := f.FocusPrev()
			return f, cmd
		case "tab", "down", "right":
			cmd := f.FocusNext()
			return f, cmd
		case "enter", "ctrl+s":
			if f.focused != sendIndex && msg.String() != "ctrl+s" {
				var cmd tea.Cmd
				f.contentText, cmd = f.contentText.Update(msg)
				return f, cmd
			}
			f.focused = sendIndex
			values, isValid := f.Validate()
			if isValid {
				return f, tea.Cmd(func() tea.Msg {
					return FormComplete{name: values[0], email: values[1], content: values[2]}
				})
			}
		}
	}

	var cmds [4]tea.Cmd
	f.nameInput, cmds[0] = f.nameInput.Update(msg)
	f.emailInput, cmds[1] = f.emailInput.Update(msg)
	f.contentText, cmds[2] = f.contentText.Update(msg)
	f.captchaInput, cmds[3] = f.captchaInput.Update(msg)
	return f, tea.Batch(cmds[:]...)
}

func (f form) View() string {

	renderWithBorder := func(index int, view string) string {
		border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
		if index == f.focused {
			border = border.BorderForeground(lipgloss.Color("33"))
		} else if f.invalid[index] {
			border = border.BorderForeground(lipgloss.Color("160"))
		}

		return border.Render(view)
	}

	renderButton := func(index int, label string) string {

		buttonStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("10")).
			Bold(true).
			Margin(1).
			Padding(1, 3)

		if f.focused == sendIndex {
			buttonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("33")).
				Underline(true)
		}

		return buttonStyle.Render(label)
	}

	responsiveJoin := lipgloss.JoinHorizontal
	forceBreak := "•"
	if f.width <= responsiveBreakpoint {
		responsiveJoin = lipgloss.JoinVertical
		forceBreak = "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().MaxWidth(f.width).MaxHeight(8).Foreground(lipgloss.Color("10")).Bold(true).Render(Config.banner),
		responsiveJoin(lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left, "Your name:", renderWithBorder(0, f.nameInput.View())),
			lipgloss.JoinVertical(lipgloss.Left, "Your email (optional):", renderWithBorder(1, f.emailInput.View())),
		),
		lipgloss.JoinVertical(lipgloss.Left, "Your message:", renderWithBorder(2, f.contentText.View())),
		responsiveJoin(lipgloss.Center,
			"Copy the code:\n"+lipgloss.PlaceHorizontal(f.width-16, lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Top, f.captcha.View(), "\n => ", renderWithBorder(3, f.captchaInput.View()))),
			renderButton(4, "Send !"),
		),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf("Tab/↓/→: Next • Shift+Tab/↑/←: Prev %s ⌃s: Send • Esc: Quit • 0.1.0", forceBreak)),
	)
}
