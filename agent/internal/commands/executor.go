package commands

import (
    "bytes"
    "os/exec"
    "runtime"
    "strings"

    "infrapilot/agent/internal/config"
)

// executeEngine runs a single command and returns its output and execution status.
// It enforces the allow‑list defined in the agent configuration.
func executeEngine(command string) (string, string) {
    var out bytes.Buffer
    status := "Completed"

    // -----------------------------------------------------------------
    // 1️⃣ Allow‑list enforcement – read from agent config
    // -----------------------------------------------------------------
    allowed := config.GetAllowedCommands()
    allowedOK := false
    for _, a := range allowed {
        if strings.HasPrefix(command, a) {
            allowedOK = true
            break
        }
    }
    if !allowedOK {
        return "Command not allowed by policy", "Failed"
    }

    // -----------------------------------------------------------------
    // 2️⃣ Platform‑specific execution
    // -----------------------------------------------------------------
    var execCmd *exec.Cmd
    if runtime.GOOS == "windows" {
        execCmd = exec.Command("cmd", "/C", command)
    } else {
        execCmd = exec.Command("bash", "-c", command)
    }
    execCmd.Stdout = &out
    execCmd.Stderr = &out

    if err := execCmd.Run(); err != nil {
        status = "Failed"
        out.WriteString("\n")
        out.WriteString(err.Error())
    }

    return out.String(), status
}
