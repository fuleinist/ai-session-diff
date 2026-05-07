package cli

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fuleinist/ai-session-diff/internal/session"
)

func Export(format string, outputPath string) {
	sessions, err := session.ListAll()
	if err != nil {
		fmt.Printf("✗ Failed to list sessions: %v\n", err)
		os.Exit(1)
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions to export.")
		os.Exit(0)
	}

	var data []byte

	switch format {
	case "json":
		data, err = json.MarshalIndent(sessions, "", "  ")
	case "csv":
		data, err = exportCSV(sessions)
	default:
		fmt.Printf("✗ Unknown format: %s (use json or csv)\n", format)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("✗ Failed to format data: %v\n", err)
		os.Exit(1)
	}

	if outputPath == "" || outputPath == "-" {
		os.Stdout.Write(data)
		return
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Printf("✗ Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Exported %d sessions to %s\n", len(sessions), outputPath)
}

func exportCSV(sessions []session.Session) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	rows := [][]string{
		{"id", "startedAt", "commitSha", "branch", "aiProvider", "aiDetected", "linesAdded", "linesRemoved", "message"},
	}
	for _, s := range sessions {
		rows = append(rows, []string{
			s.ID,
			s.StartedAt,
			s.CommitSha,
			s.Branch,
			s.AIProvider,
			fmt.Sprintf("%t", s.AIDetected),
			fmt.Sprintf("%d", s.LinesAdded),
			fmt.Sprintf("%d", s.LinesRemoved),
			s.Message,
		})
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), nil
}
