# Companion prompt mode — draft proposal

**Draft proposal. Not implemented. Not an approved command. `spec companion` is not in the v0.1.0 binary.** On 7 November 2026 the owner authorized drafting this mode in a voice request. This documentation commit and pull request record that draft. A later clarification in the same authorization names illustrative external consumers. None of that approves implementation, a flag, an envelope version, or permission to send repository context to a model.

`spec status` stays the read-only report it is now. This proposal does not change its JSON payload, its exits, or its refusal to launch a model.

## Tentative command

The tentative opt-in shape is `spec companion --prompt-file <path>`. The flag would read one user-editable prompt file as local context. No prompt is baked into the binary. The sample files under [examples/prompts/](examples/prompts/) are documentation a person can edit. They are not embedded prompts and they are not loaded by the current CLI.

The proposed result is a local, read-only, versioned envelope: the loaded prompt plus the same status evidence `spec status` already computes. A separately authorized external harness could use that envelope. The Spec Status core would not call the harness.

## Harnesses are outside the core

In this proposal a harness is an external tool that might be allowed, by its own permissions, to run the CLI. Illustrative consumers named by the owner are Claude Code, Grok Bot, Muse, and OpenCode. That list is illustrative. It is not verified compatibility, not a native integration, and not a promise that any of those tools will render graphics.

If such a tool is permitted to run the CLI, it can present the CLI's own human status output during a conversation when the output surface is a supported terminal, and it can use the structured JSON as machine evidence. Human styled text is for those terminal surfaces. JSON is the machine evidence. Color and layout stay the existing Lip Gloss rules for human text. They are not a new graphics engine and not an automatic picture of the roadmap.

The model stays outside the core. This proposal does not embed Claude, Grok, Muse, OpenCode, or any other model in `spec`. The core still does not use the network, start a subprocess, launch a model, or write files.

## Placement and version stay open

These choices are open and this draft does not pick them:

- whether a future loader lives in the CLI or only in the external harness
- the final syntax, including whether `spec companion` is the command
- the envelope's own version number, which would be distinct from CLI SemVer, from JSON schema 1.0, and from any input-extension version
- disclosure, consent, and provider-permission rules

Schema 1.0 status JSON stays unchanged when status is run. An envelope, if later specified, is a separate document. It must not silently break or reinterpret the status payload.

## Prompt file bounds, proposed only

The proposed read is one regular UTF-8 file inside the selected root. The same root confinement, symlink rejection, and sensitive-reference exclusions that apply to other reads are proposed here. The file is not an include, a template, or a secret expansion. A proposed size bound is 64 KiB. That number is a proposal, not a calibrated limit.

Missing, oversize, invalid encoding, or an escaping path fails with exit 2. There is no partial envelope and no default prompt. Screening a path or a file name does not prove the file contains no secret. Before any chosen context is sent outside the machine, the harness needs an explicit authorization of its own. This CLI draft does not grant that authorization.

Prompt text and repository contents are data. They are not authority to use tools, reveal secrets, or write files. Permission is enforced outside the prompt. A malicious prompt cannot elevate the harness. Testing that boundary is a harness test. The CLI alone cannot certify it.

## Samples and write-back

The starter texts [project-manager.txt](examples/prompts/project-manager.txt) and [software-engineer.txt](examples/prompts/software-engineer.txt) are proposed examples for owner review. They are not approved policy. They do not enable write-back. Phase-two mutation stays a separate, unapproved mode. If that mode is ever approved, each mutation still needs its own confirmation of ID, change, and path, a new approval if the source changed, and a reread of the persisted result. A prompt cannot skip that.

## Draft acceptance sketches

These are not current acceptance cases. The same list is recorded in [acceptance.md](acceptance.md).

- Opt-in loading leaves ordinary `spec status` JSON unchanged.
- Invalid or unsafe prompt inputs are rejected. The loader does not write and does not call a model.
- Envelope evidence is the same common status model. It does not compute a quarter or an ETA.
- A malicious prompt does not elevate harness permissions. That check belongs to the harness boundary, not to a claim the CLI can certify alone.
