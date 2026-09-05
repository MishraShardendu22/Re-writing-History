package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

func LogPhaseStart(phase string, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("phase", phase), slog.String("phase_event", "start"))
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "phase start", attrs...)
}

func LogPhaseEnd(phase string, duration time.Duration, attrs ...slog.Attr) {
	attrs = append(attrs,
		slog.String("phase", phase),
		slog.String("phase_event", "end"),
		slog.Duration("duration", duration),
	)
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "phase end", attrs...)
}

func LogPhaseError(phase string, err error, duration time.Duration, attrs ...slog.Attr) {
	attrs = append(attrs,
		slog.String("phase", phase),
		slog.String("phase_event", "error"),
		slog.Duration("duration", duration),
		slog.String("error", err.Error()),
	)
	Logger.LogAttrs(context.Background(), slog.LevelError, "phase error", attrs...)
}

func NewGitCommand(name string, args ...string) *GitCommand {
	return &GitCommand{Name: name, Args: args}
}

func NewGitCommandInDir(dir, name string, args ...string) *GitCommand {
	return &GitCommand{Name: name, Args: args, Dir: dir}
}

func (gc *GitCommand) Run() (string, error) {
	cmd := exec.Command(gc.Name, gc.Args...)
	if gc.Dir != "" {
		cmd.Dir = gc.Dir
	}

	start := time.Now()
	out, err := cmd.CombinedOutput()
	duration := time.Since(start)

	attrs := []slog.Attr{
		slog.String("command", gc.Name),
		slog.String("args", strings.Join(gc.Args, " ")),
		slog.String("dir", gc.Dir),
		slog.Duration("duration", duration),
	}

	if err != nil {
		attrs = append(attrs,
			slog.String("output", string(out)),
			slog.String("error", err.Error()),
		)
		Logger.LogAttrs(context.Background(), slog.LevelError, "command failed", attrs...)
		return string(out), fmt.Errorf("%s %s: %w\n%s", gc.Name, strings.Join(gc.Args, " "), err, string(out))
	}

	attrs = append(attrs, slog.String("output", strings.TrimSpace(string(out))))
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "command succeeded", attrs...)
	return string(out), nil
}

func (gc *GitCommand) RunWithStdout() error {
	cmd := exec.Command(gc.Name, gc.Args...)
	if gc.Dir != "" {
		cmd.Dir = gc.Dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	attrs := []slog.Attr{
		slog.String("command", gc.Name),
		slog.String("args", strings.Join(gc.Args, " ")),
		slog.String("dir", gc.Dir),
		slog.Duration("duration", duration),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		Logger.LogAttrs(context.Background(), slog.LevelError, "command failed", attrs...)
		return fmt.Errorf("%s %s: %w", gc.Name, strings.Join(gc.Args, " "), err)
	}

	Logger.LogAttrs(context.Background(), slog.LevelInfo, "command succeeded", attrs...)
	return nil
}

var activeJobs atomic.Int64

func WorkerStarted(worker string, attrs ...slog.Attr) {
	activeJobs.Add(1)
	attrs = append(attrs,
		slog.String("worker", worker),
		slog.String("worker_event", "started"),
		slog.Int64("active_jobs", activeJobs.Load()),
	)
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "worker started", attrs...)
}

func WorkerFinished(worker string, attrs ...slog.Attr) {
	activeJobs.Add(-1)
	attrs = append(attrs,
		slog.String("worker", worker),
		slog.String("worker_event", "finished"),
		slog.Int64("active_jobs", activeJobs.Load()),
	)
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "worker finished", attrs...)
}

func ActiveJobs() int64 {
	return activeJobs.Load()
}
