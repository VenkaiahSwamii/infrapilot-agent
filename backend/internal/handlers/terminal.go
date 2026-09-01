package handlers

import (
	"io"
	"log"
	"os/exec"
	"runtime"

	"infrapilot/backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

// HandleTerminal upgrades the connection to WebSocket and binds it to a live shell session
func HandleTerminal(c *gin.Context) {
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[Terminal] Upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe")
	} else {
		cmd = exec.Command("bash")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("[Terminal] StdinPipe error: %v", err)
		return
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[Terminal] StdoutPipe error: %v", err)
		return
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		log.Printf("[Terminal] Start error: %v", err)
		return
	}
	defer cmd.Process.Kill()

	// Goroutine reading shell output and writing to WebSocket
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				err = conn.WriteMessage(1, buf[:n])
				if err != nil {
					break
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
		}
	}()

	// Loop reading WebSocket messages and writing to shell stdin
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		_, err = stdin.Write(msg)
		if err != nil {
			break
		}
	}

	cmd.Wait()
}
