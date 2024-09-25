// package handlers

// import (
// 	"io"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// 	"golang.org/x/crypto/ssh"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Failed to set WebSocket upgrade:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// SSH 연결 설정
// 	sshClient, err := connectSSH("10.10.30.215", "root", "thddlsrbsjdlaak")
// 	if err != nil {
// 		log.Println("SSH connection failed:", err)
// 		conn.WriteMessage(websocket.TextMessage, []byte("Error: SSH connection failed"))
// 		return
// 	}
// 	defer sshClient.Close()

// 	// SSH 세션 생성
// 	session, err := sshClient.NewSession()
// 	if err != nil {
// 		log.Println("Failed to create SSH session:", err)
// 		return
// 	}
// 	defer session.Close()

// 	// PTY 요청
// 	err = session.RequestPty("xterm", 80, 124, ssh.TerminalModes{})
// 	if err != nil {
// 		log.Println("Request PTY failed:", err)
// 		return
// 	}

// 	// 세션에서 명령 실행
// 	stdin, err := session.StdinPipe()
// 	if err != nil {
// 		log.Println("Failed to get StdinPipe:", err)
// 		return
// 	}
// 	stdout, err := session.StdoutPipe()
// 	if err != nil {
// 		log.Println("Failed to get StdoutPipe:", err)
// 		return
// 	}

// 	err = session.Shell()
// 	if err != nil {
// 		log.Println("Failed to start shell:", err)
// 		return
// 	}

// 	// 서버로부터 출력받아 WebSocket을 통해 클라이언트로 전송
// 	go func() {
// 		buf := make([]byte, 1024)
// 		for {
// 			n, err := stdout.Read(buf)
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Println("Read error:", err)
// 				break
// 			}
// 			conn.WriteMessage(websocket.TextMessage, buf[:n])
// 		}
// 	}()

// 	// 클라이언트로부터 입력받아 SSH 세션에 전달
// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("WebSocket read error:", err)
// 			break
// 		}
// 		_, err = stdin.Write(message)
// 		if err != nil {
// 			log.Println("Failed to write to SSH stdin:", err)
// 			break
// 		}
// 	}
// }

// func connectSSH(host, username, password string) (*ssh.Client, error) {
// 	config := &ssh.ClientConfig{
// 		User: username,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(password),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 호스트 키 확인 생략
// 	}

// 	client, err := ssh.Dial("tcp", host+":22", config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client, nil
// }

package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type SSHConnectionInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	// SSH 연결 정보를 WebSocket을 통해 수신
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read WebSocket message:", err)
		return
	}

	// 수신한 메시지를 SSHConnectionInfo 구조체로 변환
	var sshInfo SSHConnectionInfo
	err = json.Unmarshal(msg, &sshInfo)
	if err != nil {
		log.Println("Failed to parse connection info:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Invalid connection information"))
		return
	}

	// SSH 연결 설정
	sshClient, err := connectSSH(sshInfo.Host, sshInfo.Username, sshInfo.Password, sshInfo.Port)
	if err != nil {
		log.Println("SSH connection failed:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: SSH connection failed"))
		return
	}
	defer sshClient.Close()

	// SSH 세션 생성
	session, err := sshClient.NewSession()
	if err != nil {
		log.Println("Failed to create SSH session:", err)
		return
	}
	defer session.Close()

	// PTY 요청
	err = session.RequestPty("xterm", 80, 100, ssh.TerminalModes{})
	if err != nil {
		log.Println("Request PTY failed:", err)
		return
	}

	// 세션에서 명령 실행
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Println("Failed to get StdinPipe:", err)
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Println("Failed to get StdoutPipe:", err)
		return
	}

	err = session.Shell()
	if err != nil {
		log.Println("Failed to start shell:", err)
		return
	}

	// 서버로부터 출력받아 WebSocket을 통해 클라이언트로 전송
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Read error:", err)
				break
			}
			conn.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	// 클라이언트로부터 입력받아 SSH 세션에 전달
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		_, err = stdin.Write(message)
		if err != nil {
			log.Println("Failed to write to SSH stdin:", err)
			break
		}
	}
}

func connectSSH(host, username, password, port string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 호스트 키 확인 생략
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
