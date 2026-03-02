package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	minimax "github.com/giztoy/minimax-go"
)

const (
	asyncDefaultModel        = "speech-2.6-hd"
	asyncDefaultText         = "hello from minimax-go speech async example"
	asyncDefaultTimeout      = 3 * time.Minute
	asyncDefaultPollInterval = 2 * time.Second
)

type asyncOptions struct {
	apiKey       string
	baseURL      string
	taskID       string
	text         string
	textFileID   string
	model        string
	voiceID      string
	speed        *float64
	volume       *float64
	timeout      time.Duration
	pollInterval time.Duration
	wait         bool
	noWait       bool
}

func runAsyncCommand(args []string, stdout, stderr io.Writer) error {
	opts, err := parseAsyncOptions(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return fmt.Errorf("failed to parse async flags: %w", err)
	}

	return runAsync(opts, stdout)
}

func runTaskAliasCommand(args []string, stdout, stderr io.Writer) error {
	opts, err := parseTaskAliasOptions(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return fmt.Errorf("failed to parse task flags: %w", err)
	}

	return runAsync(opts, stdout)
}

func parseTaskAliasOptions(args []string, out io.Writer) (asyncOptions, error) {
	var opts asyncOptions

	apiKeyDefault := os.Getenv("MINIMAX_API_KEY")
	baseURLDefault := envOrDefault("MINIMAX_BASE_URL", exampleDefaultBaseURL)
	timeoutDefault := envDurationOrDefaultFromKeys(
		[]string{"MINIMAX_SPEECH_TASK_TIMEOUT", "MINIMAX_SPEECH_ASYNC_TIMEOUT", "MINIMAX_SPEECH_TIMEOUT"},
		asyncDefaultTimeout,
	)
	pollIntervalDefault := envDurationOrDefaultFromKeys(
		[]string{"MINIMAX_SPEECH_TASK_POLL_INTERVAL", "MINIMAX_SPEECH_ASYNC_POLL_INTERVAL"},
		asyncDefaultPollInterval,
	)

	waitDefault := true
	if waitFromEnv, waitSetByEnv, waitErr := optionalEnvBool("MINIMAX_SPEECH_ASYNC_WAIT"); waitErr != nil {
		return asyncOptions{}, fmt.Errorf("invalid MINIMAX_SPEECH_ASYNC_WAIT: %w", waitErr)
	} else if waitSetByEnv {
		waitDefault = waitFromEnv
	}

	if noWaitFromEnv, noWaitSetByEnv, noWaitErr := optionalEnvBool("MINIMAX_SPEECH_ASYNC_NO_WAIT"); noWaitErr != nil {
		return asyncOptions{}, fmt.Errorf("invalid MINIMAX_SPEECH_ASYNC_NO_WAIT: %w", noWaitErr)
	} else if noWaitSetByEnv {
		waitDefault = !noWaitFromEnv
	}

	fs := flag.NewFlagSet("task", flag.ContinueOnError)
	fs.SetOutput(out)

	fs.StringVar(&opts.apiKey, "api-key", apiKeyDefault, "Minimax API key (or env MINIMAX_API_KEY)")
	fs.StringVar(&opts.baseURL, "base-url", baseURLDefault, "Minimax API base URL (env: MINIMAX_BASE_URL)")
	fs.StringVar(&opts.taskID, "task-id", "", "Task ID to query (required)")
	fs.DurationVar(&opts.timeout, "timeout", timeoutDefault, "Total timeout for query/watch workflow")
	fs.DurationVar(&opts.pollInterval, "poll-interval", pollIntervalDefault, "Polling interval when wait=true")
	fs.BoolVar(&opts.wait, "wait", waitDefault, "Wait/poll until terminal status")
	fs.BoolVar(&opts.noWait, "no-wait", !waitDefault, "Alias of -wait=false")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: go run ./examples/speech task -task-id <id> [flags]\n\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nNotes:\n")
		fmt.Fprintf(fs.Output(), "  - task is a query-only alias and cannot submit new text\n")
	}

	if err := fs.Parse(args); err != nil {
		return asyncOptions{}, err
	}

	if flagWasSet(fs, "wait") && flagWasSet(fs, "no-wait") {
		return asyncOptions{}, errors.New("wait and no-wait cannot be set together")
	}
	if flagWasSet(fs, "no-wait") {
		opts.wait = !opts.noWait
	}

	opts.apiKey = strings.TrimSpace(opts.apiKey)
	opts.baseURL = strings.TrimSpace(opts.baseURL)
	opts.taskID = strings.TrimSpace(opts.taskID)

	if opts.timeout <= 0 {
		return asyncOptions{}, errors.New("timeout must be greater than 0")
	}

	if opts.pollInterval <= 0 {
		return asyncOptions{}, errors.New("poll-interval must be greater than 0")
	}

	if !flagWasSet(fs, "task-id") || opts.taskID == "" {
		fs.Usage()
		return asyncOptions{}, errors.New("task command requires -task-id")
	}

	// Ensure task alias remains query-only.
	opts.text = ""
	opts.textFileID = ""

	return opts, nil
}

func parseAsyncOptions(args []string, out io.Writer) (asyncOptions, error) {
	var opts asyncOptions

	apiKeyDefault := os.Getenv("MINIMAX_API_KEY")
	baseURLDefault := envOrDefault("MINIMAX_BASE_URL", exampleDefaultBaseURL)
	taskIDDefault := strings.TrimSpace(envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_TASK_ID", "MINIMAX_SPEECH_ASYNC_TASK_ID"}, ""))
	textDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_ASYNC_TEXT", "MINIMAX_SPEECH_TEXT"}, asyncDefaultText)
	textFileIDDefault := os.Getenv("MINIMAX_SPEECH_ASYNC_TEXT_FILE_ID")
	modelDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_ASYNC_MODEL", "MINIMAX_SPEECH_MODEL"}, asyncDefaultModel)
	voiceDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_ASYNC_VOICE_ID", "MINIMAX_SPEECH_VOICE_ID"}, "")

	speedDefault, speedSetByEnv, err := optionalEnvFloat64FromKeys("MINIMAX_SPEECH_ASYNC_SPEED", "MINIMAX_SPEECH_SPEED")
	if err != nil {
		return asyncOptions{}, fmt.Errorf("invalid speech speed env: %w", err)
	}
	if !speedSetByEnv {
		speedDefault = 1
	}

	volumeDefault, volumeSetByEnv, err := optionalEnvFloat64FromKeys("MINIMAX_SPEECH_ASYNC_VOLUME", "MINIMAX_SPEECH_VOLUME")
	if err != nil {
		return asyncOptions{}, fmt.Errorf("invalid speech volume env: %w", err)
	}
	if !volumeSetByEnv {
		volumeDefault = 1
	}

	timeoutDefault := envDurationOrDefaultFromKeys([]string{"MINIMAX_SPEECH_ASYNC_TIMEOUT", "MINIMAX_SPEECH_TIMEOUT"}, asyncDefaultTimeout)
	pollIntervalDefault := envDurationOrDefaultFromKeys([]string{"MINIMAX_SPEECH_ASYNC_POLL_INTERVAL", "MINIMAX_SPEECH_TASK_POLL_INTERVAL"}, asyncDefaultPollInterval)

	waitDefault := true
	if waitFromEnv, waitSetByEnv, waitErr := optionalEnvBool("MINIMAX_SPEECH_ASYNC_WAIT"); waitErr != nil {
		return asyncOptions{}, fmt.Errorf("invalid MINIMAX_SPEECH_ASYNC_WAIT: %w", waitErr)
	} else if waitSetByEnv {
		waitDefault = waitFromEnv
	}

	if noWaitFromEnv, noWaitSetByEnv, noWaitErr := optionalEnvBool("MINIMAX_SPEECH_ASYNC_NO_WAIT"); noWaitErr != nil {
		return asyncOptions{}, fmt.Errorf("invalid MINIMAX_SPEECH_ASYNC_NO_WAIT: %w", noWaitErr)
	} else if noWaitSetByEnv {
		waitDefault = !noWaitFromEnv
	}

	fs := flag.NewFlagSet("async", flag.ContinueOnError)
	fs.SetOutput(out)

	speedValue := speedDefault
	volumeValue := volumeDefault

	fs.StringVar(&opts.apiKey, "api-key", apiKeyDefault, "Minimax API key (or env MINIMAX_API_KEY)")
	fs.StringVar(&opts.baseURL, "base-url", baseURLDefault, "Minimax API base URL (env: MINIMAX_BASE_URL)")
	fs.StringVar(&opts.taskID, "task-id", taskIDDefault, "Query existing task_id instead of submitting new text")
	fs.StringVar(&opts.text, "text", textDefault, "Input text for submit mode (env: MINIMAX_SPEECH_ASYNC_TEXT)")
	fs.StringVar(&opts.textFileID, "text-file-id", textFileIDDefault, "Uploaded text file_id for submit mode (optional)")
	fs.StringVar(&opts.model, "model", modelDefault, "Model name for submit mode (optional)")
	fs.StringVar(&opts.voiceID, "voice-id", voiceDefault, "Voice ID for submit mode (optional; some regions require this)")
	fs.Float64Var(&speedValue, "speed", speedDefault, "Speech speed for submit mode (optional)")
	fs.Float64Var(&volumeValue, "volume", volumeDefault, "Speech volume for submit mode (optional)")
	fs.DurationVar(&opts.timeout, "timeout", timeoutDefault, "Total timeout for submit/query workflow")
	fs.DurationVar(&opts.pollInterval, "poll-interval", pollIntervalDefault, "Polling interval when wait=true")
	fs.BoolVar(&opts.wait, "wait", waitDefault, "Wait/poll until terminal status")
	fs.BoolVar(&opts.noWait, "no-wait", !waitDefault, "Alias of -wait=false")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: go run ./examples/speech async [flags]\n\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nModes:\n")
		fmt.Fprintf(fs.Output(), "  - submit mode: no -task-id, uses -text/-text-file-id\n")
		fmt.Fprintf(fs.Output(), "  - task mode: set -task-id, query existing task\n")
	}

	if err := fs.Parse(args); err != nil {
		return asyncOptions{}, err
	}

	if flagWasSet(fs, "wait") && flagWasSet(fs, "no-wait") {
		return asyncOptions{}, errors.New("wait and no-wait cannot be set together")
	}
	if flagWasSet(fs, "no-wait") {
		opts.wait = !opts.noWait
	}

	opts.apiKey = strings.TrimSpace(opts.apiKey)
	opts.baseURL = strings.TrimSpace(opts.baseURL)
	opts.taskID = strings.TrimSpace(opts.taskID)
	opts.text = strings.TrimSpace(opts.text)
	opts.textFileID = strings.TrimSpace(opts.textFileID)
	opts.model = strings.TrimSpace(opts.model)
	opts.voiceID = strings.TrimSpace(opts.voiceID)

	if opts.timeout <= 0 {
		return asyncOptions{}, errors.New("timeout must be greater than 0")
	}

	if opts.pollInterval <= 0 {
		return asyncOptions{}, errors.New("poll-interval must be greater than 0")
	}

	if opts.taskID == "" && opts.text == "" && opts.textFileID == "" {
		return asyncOptions{}, errors.New("submit mode requires text or text-file-id")
	}

	if speedSetByEnv || flagWasSet(fs, "speed") {
		speed := speedValue
		opts.speed = &speed
	}

	if volumeSetByEnv || flagWasSet(fs, "volume") {
		volume := volumeValue
		opts.volume = &volume
	}

	return opts, nil
}

func runAsync(opts asyncOptions, out io.Writer) error {
	if opts.apiKey == "" {
		return errors.New("missing API key: use -api-key or set MINIMAX_API_KEY")
	}

	if opts.baseURL == "" {
		return errors.New("base-url cannot be empty")
	}

	client, err := minimax.NewClient(minimax.Config{
		BaseURL: opts.baseURL,
		APIKey:  opts.apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create Minimax client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
	defer cancel()

	if opts.taskID != "" {
		return runAsyncTaskMode(ctx, client, opts, out)
	}

	return runAsyncSubmitMode(ctx, client, opts, out)
}

func runAsyncTaskMode(ctx context.Context, client *minimax.Client, opts asyncOptions, out io.Writer) error {
	var (
		statusResp *minimax.SpeechTaskStatusResponse
		err        error
	)

	if opts.wait {
		statusResp, err = waitTask(ctx, client, opts.taskID, opts.pollInterval, out)
		if err != nil {
			return fmt.Errorf("wait async task failed: %w", err)
		}
	} else {
		statusResp, err = client.SpeechAsync.GetAsyncTask(ctx, opts.taskID)
		if err != nil {
			return fmt.Errorf("SpeechAsync.GetAsyncTask failed: %w", err)
		}
	}

	if statusResp.Status == minimax.SpeechTaskStateFailed {
		printTaskResult(out, statusResp)
		return fmt.Errorf("task_id=%s failed: %s", statusResp.TaskID, resolveTaskFailureMessage(statusResp))
	}

	printTaskResult(out, statusResp)
	return nil
}

func runAsyncSubmitMode(ctx context.Context, client *minimax.Client, opts asyncOptions, out io.Writer) error {
	submitResp, err := client.SpeechAsync.SubmitAsync(ctx, minimax.SpeechAsyncSubmitRequest{
		Model:      opts.model,
		Text:       opts.text,
		TextFileID: opts.textFileID,
		VoiceID:    opts.voiceID,
		Speed:      opts.speed,
		Vol:        opts.volume,
	})
	if err != nil {
		return fmt.Errorf("SpeechAsync.SubmitAsync failed: %w", err)
	}

	fmt.Fprintf(out, "async submit succeeded: task_id=%s status=%s file_id=%s\n", submitResp.TaskID, displayTaskState(submitResp.Status), submitResp.FileID)

	if !opts.wait {
		fmt.Fprintln(out, "wait=false, skip polling")
		return nil
	}

	statusResp, err := waitTask(ctx, client, submitResp.TaskID, opts.pollInterval, out)
	if err != nil {
		return fmt.Errorf("wait async task failed: %w", err)
	}

	if statusResp.Status == minimax.SpeechTaskStateFailed {
		return fmt.Errorf("task_id=%s failed: %s", statusResp.TaskID, resolveTaskFailureMessage(statusResp))
	}

	printTaskResult(out, statusResp)
	return nil
}
