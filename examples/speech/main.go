package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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
	defaultText    = "hello from minimax-go example"
	defaultOutFile = "speech_output.audio"
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
		fmt.Fprintf(os.Stderr, "speech example failed: %v\n", err)
		os.Exit(1)
	}
}

func parseOptions() (options, error) {
	var opts options

	apiKeyDefault := os.Getenv("MINIMAX_API_KEY")
	baseURLDefault := envOrDefault("MINIMAX_BASE_URL", defaultBaseURL)
	textDefault := envOrDefault("MINIMAX_SPEECH_TEXT", defaultText)
	modelDefault := envOrDefault("MINIMAX_SPEECH_MODEL", defaultModel)
	voiceDefault := os.Getenv("MINIMAX_SPEECH_VOICE_ID")
	speedDefault, speedSetByEnv, err := optionalEnvFloat64("MINIMAX_SPEECH_SPEED")
	if err != nil {
		return options{}, fmt.Errorf("invalid MINIMAX_SPEECH_SPEED: %w", err)
	}
	if !speedSetByEnv {
		speedDefault = 1
	}

	volumeDefault, volumeSetByEnv, err := optionalEnvFloat64("MINIMAX_SPEECH_VOLUME")
	if err != nil {
		return options{}, fmt.Errorf("invalid MINIMAX_SPEECH_VOLUME: %w", err)
	}
	if !volumeSetByEnv {
		volumeDefault = 1
	}

	timeoutDefault := envDurationOrDefault("MINIMAX_SPEECH_TIMEOUT", 30*time.Second)
	outputDefault := envOrDefault("MINIMAX_SPEECH_OUTPUT", defaultOutFile)

	speedValue := speedDefault
	volumeValue := volumeDefault

	flag.StringVar(&opts.apiKey, "api-key", apiKeyDefault, "Minimax API key (or env MINIMAX_API_KEY)")
	flag.StringVar(&opts.baseURL, "base-url", baseURLDefault, "Minimax API base URL (env: MINIMAX_BASE_URL)")
	flag.StringVar(&opts.text, "text", textDefault, "Text to synthesize (env: MINIMAX_SPEECH_TEXT)")
	flag.StringVar(&opts.model, "model", modelDefault, "Model name (optional, env: MINIMAX_SPEECH_MODEL)")
	flag.StringVar(&opts.voiceID, "voice-id", voiceDefault, "Voice ID (optional, env: MINIMAX_SPEECH_VOICE_ID)")
	flag.Float64Var(&speedValue, "speed", speedDefault, "Speech speed (optional, env: MINIMAX_SPEECH_SPEED)")
	flag.Float64Var(&volumeValue, "volume", volumeDefault, "Speech volume (optional, env: MINIMAX_SPEECH_VOLUME)")
	flag.DurationVar(&opts.timeout, "timeout", timeoutDefault, "Request timeout (env: MINIMAX_SPEECH_TIMEOUT, e.g. 30s)")
	flag.StringVar(&opts.output, "output", outputDefault, "Output audio file path (env: MINIMAX_SPEECH_OUTPUT)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: go run ./examples/speech [flags]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nNotes:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  - API key precedence: -api-key > MINIMAX_API_KEY\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  - action_env in .bazelrc.user is not injected into go run automatically; export env manually or use -api-key\n")
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

func run(opts options) error {
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

	response, err := client.Speech.Synthesize(ctx, minimax.SpeechRequest{
		Model:   opts.model,
		Text:    opts.text,
		VoiceID: opts.voiceID,
		Speed:   opts.speed,
		Vol:     opts.volume,
	})
	if err != nil {
		return fmt.Errorf("Speech.Synthesize failed: %w", err)
	}

	if len(response.Audio) == 0 {
		return errors.New("synthesis succeeded but returned empty audio bytes")
	}

	if err := ensureOutputDir(opts.output); err != nil {
		return err
	}

	if err := os.WriteFile(opts.output, response.Audio, 0o644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("synthesis succeeded, wrote %d bytes to %s\n", len(response.Audio), opts.output)
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

func optionalEnvFloat64(key string) (float64, bool, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return 0, false, nil
	}

	parsed, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, true, err
	}

	return parsed, true, nil
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

func envDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return defaultValue
	}

	parsed, err := time.ParseDuration(strings.TrimSpace(raw))
	if err != nil {
		return defaultValue
	}

	if parsed <= 0 {
		return defaultValue
	}

	return parsed
}
