package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	minimax "github.com/giztoy/minimax-go"
)

const (
	defaultBaseURL = "https://api.minimax.io"
	defaultModel   = "speech-2.6-hd"
	defaultText    = "hello from minimax-go speech stream example"
	defaultOutFile = "speech_stream_output.audio"
)

type options struct {
	apiKey  string
	baseURL string
	text    string
	model   string
	voiceID string
	speed   *float64
	volume  *float64
	timeout time.Duration
	output  string
}

func main() {
	opts, err := parseOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags: %v\n", err)
		os.Exit(2)
	}

	if err := run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "speech stream example failed: %v\n", err)
		os.Exit(1)
	}
}

func parseOptions() (options, error) {
	var opts options

	apiKeyDefault := os.Getenv("MINIMAX_API_KEY")
	baseURLDefault := envOrDefault("MINIMAX_BASE_URL", defaultBaseURL)
	textDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_STREAM_TEXT", "MINIMAX_SPEECH_TEXT"}, defaultText)
	modelDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_STREAM_MODEL", "MINIMAX_SPEECH_MODEL"}, defaultModel)
	voiceDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_STREAM_VOICE_ID", "MINIMAX_SPEECH_VOICE_ID"}, "")

	speedDefault, speedSetByEnv, err := optionalEnvFloat64FromKeys("MINIMAX_SPEECH_STREAM_SPEED", "MINIMAX_SPEECH_SPEED")
	if err != nil {
		return options{}, fmt.Errorf("invalid speech speed env: %w", err)
	}
	if !speedSetByEnv {
		speedDefault = 1
	}

	volumeDefault, volumeSetByEnv, err := optionalEnvFloat64FromKeys("MINIMAX_SPEECH_STREAM_VOLUME", "MINIMAX_SPEECH_VOLUME")
	if err != nil {
		return options{}, fmt.Errorf("invalid speech volume env: %w", err)
	}
	if !volumeSetByEnv {
		volumeDefault = 1
	}

	timeoutDefault := envDurationOrDefaultFromKeys([]string{"MINIMAX_SPEECH_STREAM_TIMEOUT", "MINIMAX_SPEECH_TIMEOUT"}, 30*time.Second)
	outputDefault := envOrDefaultFromKeys([]string{"MINIMAX_SPEECH_STREAM_OUTPUT", "MINIMAX_SPEECH_OUTPUT"}, defaultOutFile)

	speedValue := speedDefault
	volumeValue := volumeDefault

	flag.StringVar(&opts.apiKey, "api-key", apiKeyDefault, "Minimax API key (or env MINIMAX_API_KEY)")
	flag.StringVar(&opts.baseURL, "base-url", baseURLDefault, "Minimax API base URL (env: MINIMAX_BASE_URL)")
	flag.StringVar(&opts.text, "text", textDefault, "Text to synthesize (env: MINIMAX_SPEECH_STREAM_TEXT)")
	flag.StringVar(&opts.model, "model", modelDefault, "Model name (optional, env: MINIMAX_SPEECH_STREAM_MODEL)")
	flag.StringVar(&opts.voiceID, "voice-id", voiceDefault, "Voice ID (optional, env: MINIMAX_SPEECH_STREAM_VOICE_ID)")
	flag.Float64Var(&speedValue, "speed", speedDefault, "Speech speed (optional, env: MINIMAX_SPEECH_STREAM_SPEED)")
	flag.Float64Var(&volumeValue, "volume", volumeDefault, "Speech volume (optional, env: MINIMAX_SPEECH_STREAM_VOLUME)")
	flag.DurationVar(&opts.timeout, "timeout", timeoutDefault, "Request timeout (env: MINIMAX_SPEECH_STREAM_TIMEOUT, e.g. 30s)")
	flag.StringVar(&opts.output, "output", outputDefault, "Output audio file path (env: MINIMAX_SPEECH_STREAM_OUTPUT)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: go run ./examples/speech_stream [flags]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nNotes:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  - API key precedence: -api-key > MINIMAX_API_KEY\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  - stream reads incremental chunks and writes merged bytes to output file\n")
	}

	flag.Parse()

	opts.apiKey = strings.TrimSpace(opts.apiKey)
	opts.baseURL = strings.TrimSpace(opts.baseURL)
	opts.text = strings.TrimSpace(opts.text)
	opts.model = strings.TrimSpace(opts.model)
	opts.voiceID = strings.TrimSpace(opts.voiceID)
	opts.output = strings.TrimSpace(opts.output)

	if opts.timeout <= 0 {
		return options{}, errors.New("timeout must be greater than 0")
	}

	if speedSetByEnv || isFlagSet("speed") {
		speed := speedValue
		opts.speed = &speed
	}

	if volumeSetByEnv || isFlagSet("volume") {
		volume := volumeValue
		opts.volume = &volume
	}

	return opts, nil
}

func run(opts options) (retErr error) {
	if opts.apiKey == "" {
		return errors.New("missing API key: use -api-key or set MINIMAX_API_KEY")
	}

	if opts.baseURL == "" {
		return errors.New("base-url cannot be empty")
	}

	if opts.text == "" {
		return errors.New("text cannot be empty")
	}

	if opts.output == "" {
		return errors.New("output cannot be empty")
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

	speechStream, err := client.Speech.OpenStream(ctx, minimax.SpeechStreamRequest{
		Model:   opts.model,
		Text:    opts.text,
		VoiceID: opts.voiceID,
		Speed:   opts.speed,
		Vol:     opts.volume,
	})
	if err != nil {
		return fmt.Errorf("Speech.OpenStream failed: %w", err)
	}
	defer func() {
		if closeErr := speechStream.Close(); closeErr != nil && retErr == nil {
			retErr = fmt.Errorf("failed to close speech stream: %w", closeErr)
		}
	}()

	if err := ensureOutputDir(opts.output); err != nil {
		return err
	}

	outputFile, err := os.OpenFile(opts.output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil && retErr == nil {
			retErr = fmt.Errorf("failed to close output file: %w", closeErr)
		}
	}()

	totalBytes := 0
	chunkCount := 0

	for {
		chunk, nextErr := speechStream.Next()
		if nextErr != nil {
			if errors.Is(nextErr, io.EOF) {
				break
			}
			return fmt.Errorf("failed to read stream chunk: %w", nextErr)
		}

		if chunk == nil {
			continue
		}

		if len(chunk.Audio) > 0 {
			if _, writeErr := outputFile.Write(chunk.Audio); writeErr != nil {
				return fmt.Errorf("failed to write audio chunk: %w", writeErr)
			}
			totalBytes += len(chunk.Audio)
			chunkCount++
		}

		if chunk.Done {
			break
		}
	}

	if err := outputFile.Sync(); err != nil {
		return fmt.Errorf("failed to flush output file: %w", err)
	}

	if totalBytes == 0 {
		return errors.New("stream synthesis finished but no audio chunk was received")
	}

	fmt.Printf("stream synthesis succeeded, wrote %d bytes from %d chunks to %s\n", totalBytes, chunkCount, opts.output)
	return nil
}

func ensureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	if dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}

func envOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func envOrDefaultFromKeys(keys []string, defaultValue string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
	}

	return defaultValue
}

func optionalEnvFloat64FromKeys(keys ...string) (float64, bool, error) {
	for _, key := range keys {
		raw, ok := os.LookupEnv(key)
		if !ok || strings.TrimSpace(raw) == "" {
			continue
		}

		parsed, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return 0, true, fmt.Errorf("%s: %w", key, err)
		}

		return parsed, true, nil
	}

	return 0, false, nil
}

func isFlagSet(name string) bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})

	return set
}

func envDurationOrDefaultFromKeys(keys []string, defaultValue time.Duration) time.Duration {
	for _, key := range keys {
		raw, ok := os.LookupEnv(key)
		if !ok || strings.TrimSpace(raw) == "" {
			continue
		}

		parsed, err := time.ParseDuration(strings.TrimSpace(raw))
		if err != nil {
			continue
		}

		if parsed <= 0 {
			continue
		}

		return parsed
	}

	return defaultValue
}
