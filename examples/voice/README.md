# Voice Examples

Voice-related examples are organized under `examples/voice/` with three subdirectories:

1. `list/` — list available voices (`Voice.ListVoices`)
2. `design/` — design a custom voice from prompt text (`Voice.DesignVoice`)
3. `clone/` — clone a voice from `audio_url` / `file_id` / local upload (`Voice.CloneVoice`)

## Quick links

- List: `examples/voice/list/README.md`
- Design: `examples/voice/design/README.md`
- Clone: `examples/voice/clone/README.md`

## Notes

- The full snapshot of **official non-cloning voices** is maintained in `examples/voice/list/README.md`.
- In China region (`https://api.minimax.chat`), `audio_url` clone may be unsupported; use `file_id` or local upload flow instead.
