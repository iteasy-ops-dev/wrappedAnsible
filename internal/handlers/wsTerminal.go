// TODO:
// 소켓이 끊어질 때 해당 세션도 끊어져야함.
// vim안되는 이유찾아야 함
package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

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

	// WebSocket과 SSH 연결이 닫힐 때 cleanup을 보장하는 채널
	done := make(chan struct{})

	// SSH 연결 정보를 WebSocket을 통해 수신
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read WebSocket message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to read message"))
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
		// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer sshClient.Close()

	// SSH 세션 생성
	session, err := sshClient.NewSession()
	if err != nil {
		log.Println("Failed to create SSH session:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to create SSH session"))
		// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer session.Close()

	// PTY 요청
	err = session.RequestPty("xterm", 80, 100, ssh.TerminalModes{})
	if err != nil {
		log.Println("Request PTY failed:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: PTY request failed"))
		// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	// 세션에서 명령 실행
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Println("Failed to get StdinPipe:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to get stdin pipe"))
		// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Println("Failed to get StdoutPipe:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to get stdout pipe"))
		// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	err = session.Shell()
	if err != nil {
		log.Println("Failed to start shell:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to start shell"))
		// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	// 서버로부터 출력받아 WebSocket을 통해 클라이언트로 전송 (고루틴 사용)
	go func() {
		defer close(done)         // WebSocket 또는 SSH 세션이 끝나면 고루틴 종료
		buf := make([]byte, 4096) // 버퍼 크기 1024 -> 4096으로 증가
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("SSH session closed by the server")
					conn.WriteMessage(websocket.TextMessage, []byte("SSH session closed"))
					break
				}
				log.Println("Read error:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to read from SSH session"))
				// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				break
			}
			// 클라이언트에 전송
			err = conn.WriteMessage(websocket.TextMessage, buf[:n])
			if err != nil {
				log.Println("WebSocket write error:", err)
				// conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				break
			}
		}
	}()

	// 클라이언트로부터 입력받아 SSH 세션에 전달
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("WebSocket read error:", err)
					return
				}
				/**
				웹소켓 입력 메세지 디버깅 용
				*/
				// log.Printf("Websocket Read Byte: %s\n", message)
				// log.Printf("WebSocket Read Byte (Raw): %q\n", message)
				/**
				웹소켓 입력 메세지 디버깅 용
				*/
				_, err = stdin.Write(message)
				if err != nil {
					log.Println("Failed to write to SSH stdin:", err)
					return
				}
			}
		}
	}()

	// WebSocket이 닫히거나 에러가 발생할 때까지 대기
	<-done
}

func connectSSH(host, username, password, port string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         5 * time.Second,             // 타임아웃 추가하여 연결 지연 방지
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 호스트 키 확인 생략 (보안 강화 필요)
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
