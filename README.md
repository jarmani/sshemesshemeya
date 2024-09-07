# SSHme SSHme Ya

> Hey, baby, I like it rawww...   
> Shimmy shimmy ya, shimmy yam, shimmy yay  
> *— Ol' Dirty Bastard "Shimmy Shimmy Ya"*

A fully interactive contact form accessible via SSH, designed for secure and efficient communication without a web interface.

## Key features

- Fully configurable using environement variable
- Input validation
- Navigation using Tab/Shift-Tab
- Email submission to the configured address
- Advanced anti-spam using last CAPTCHA technology
- Satisfying animation to celebrate a successful submission

## Demo

![sshemesshemeya](https://github.com/user-attachments/assets/f332a06b-ec58-4330-a106-4851e9398f74)


## Configuration

The configuration can be customized by setting the appropriate values in a .env file or directly through the environment.

BANNER:
Type: string
Default: "WELCOME TO SSH FORM"
This sets the welcome banner that is displayed when a user connects to the form over SSH.
Example:

```sh
BANNER="\n\nWELCOME TO MY CUSTOM SSH FORM\n\n"
```

SERVER_HOST:
Type: string
Default: "localhost"
Defines the SSH server host.
Example:

```sh
SERVER_HOST="0.0.0.0"
```

SERVER_PORT:
Type: integer
Default: 2222
The port on which the SSH server listens.
Example:

```sh
SERVER_PORT=2222
```

SERVER_KEY_PATH:
Type: string
Default: ".ssh/term_info_ed25519"
Path to the SSH server private key. This key is required to authenticate SSH connections.
Example:

```sh
SERVER_KEY_PATH="/path/to/your/ssh_key"
```

*EMAIL_BODY*
Type: string  
Default: "{name} <{email}>\n{content}"  
The format of the email body that will be sent upon form submission. You can use placeholders like {name}, {email}, and {content}.  
Example:  

```sh
EMAIL_BODY="Name: {name}\nEmail: {email}\nMessage: {content}"
```

EMAIL_EXEC:
Type: string
Default: "/usr/sbin/sendmail"
Command to send the email. By default, it uses the sendmail utility, but you can specify another mailer.
Example:


```sh
EMAIL_EXEC="/usr/bin/mail"
```

EMAIL_ARGS:
Type: string
Default: ""
Optional arguments to pass to the email command when sending the message.
Example:

```sh
EMAIL_ARGS="-s 'Contact Form Submission'"
```

Usage Example

Here’s an example .env file that sets up a custom banner, server host, and email body:

```sh
BANNER="\n\nWELCOME TO MY SSH CONTACT FORM\n\n"
SERVER_HOST="0.0.0.0"
SERVER_PORT=2222
SERVER_KEY_PATH="/etc/ssh/my_server_key"
EMAIL_BODY="New message from {name}: {content}"
EMAIL_EXEC="/usr/bin/mail"
EMAIL_ARGS="-s 'New Contact Submission'"
```
