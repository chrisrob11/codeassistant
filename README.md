# **Code-Assistant (`ca`) A developer focused AI integration

## **Overview**

**Code-Assistant (**`ca`**)** is a CLI tool designed to assist developers by integrating AI-powered modifications into their code workflow. It supports **file refactoring, summaries, commit tracking, and LLM customization** via `gollm`.

This tool is designed to work seamlessly within a **Git-tracked** or **non-Git** project, allowing developers to iteratively improve their code while maintaining full control over changes.

> ⚠️ **Warning:** Code-Assistant (`ca`) is still under development. Expect changes and improvements as the project evolves.

---

## **Key Features**

✅ **AI-Powered Code Editing** – Modify code using prompts.\
✅ **Per-File or Batch Processing** – Apply changes to one or multiple files.\
✅ **LLM Customization** – Use `gollm` for local or remote LLM models.\
✅ **Git Integration (Optional)** – Track commits and rollback changes.\
✅ **Session-Based Workflow** – Log AI interactions for review and undo.\
✅ **Summarization Mode** – Review code without modifying files.\
✅ **Configurable Defaults** – Customize behavior via environment variables.

---

## **Command Specification**

| Command                                                                        | Description                                      |
| ------------------------------------------------------------------------------ | ------------------------------------------------ |
| `ca start-session "<description>"`                                             | Start a new coding session.                      |
| `ca code "<prompt>" [--files file1 file2] [--per-file] [--dry-run] [--commit]` | Apply AI modifications.                          |
| `ca review`                                                                    | Show session progress and diffs.                 |
| `ca commit "<message>"`                                                        | Commit current changes to Git.                   |
| `ca rollback --step N`                                                         | Undo a specific AI-modified step.                |
| `ca replay-step --step N [--prompt "<new prompt>"] [--files file1 file2]`      | Modify a previous AI change.                     |
| `ca end-session`                                                               | Archive session to historical storage.           |
| `ca code "<prompt>" --summary [--output file] [--store-session]`               | Generate an analysis instead of modifying files. |
| `ca config llm [--set model=gpt-4] [--list]`                                   | Manage LLM configuration.                        |

---

## **Implementation Checklist**

### Basics
- [X] Create stub commands
- [X] Add make file so it can build and do other commands
- [ ] Add start-session and end-session command
- [ ] code command work
  - [ ] Initial work to do simple prompt againt one file
  - [ ] Add ability to do multiple files
  - [ ] Add dry-run
  - [ ] Add ability to do a commit at the end and save that info in file
- [ ] review command work
  - More TBD
- [ ] rollback command work
  - More TBD
- [ ] Add the summary mode into code command

### **File Tracking & Storage**
- [ ] Implement `.ca_session.json` for tracking prompts, modified files, and steps.
- [ ] Store **created, modified, and deleted** files for each AI step.
- [ ] Track Git state (`pre-commit` and `post-commit` hashes).

---

### **Extensions Features**

- [ ] **Code Review & History**
  - [ ] `ca replay-step --step N [--prompt "<new prompt>"] [--files file1 file2]` - Re-run a specific AI step with modifications.

- [ ] **AI-Powered Summaries**
  - [ ] `ca code "<prompt>" --summary` - Generate an AI analysis instead of modifying files.
  - [ ] `--output file` - Save the AI summary to a separate file.
  - [ ] `--store-session` - Store the summary output in the `.ca_session.json` session file.

- [ ] **Git Integration**
  - [ ] Detect and track files under Git.
  - [ ] Store pre-change and post-change commit hashes in `.ca_session.json`.
  - [ ] `ca rollback --step N` - Restore files to their previous Git state.

- [ ] **Configuration Management**
  - [ ] `ca config llm --set model=gpt-4` - Configure the LLM model for AI processing.
  - [ ] `ca config llm --list` - List available LLM models.
  - [ ] Support `CA_LLM_MODEL`, `CA_STORE_SUMMARY`, and other environment variables.

---

### **CLI Behavior & UX**
- [ ] Use `urfave/cli` for handling command-line arguments.
- [ ] Implement global config defaults (`CA_DEFAULT_MODE`, `CA_STORE_SUMMARY`).
- [ ] Ensure `--summary`, `--per-file`, and `--batch` behavior is intuitive.
- [ ] Allow per-command overrides for global config settings.

---

### **LLM Integration (Using `gollm`)**
- [ ] Integrate `gollm` for AI processing.
- [ ] Support local LLM models (e.g., LLaMA, Mistral).
- [ ] Support remote APIs (e.g., OpenAI GPT).
- [ ] Allow per-command model selection (`--llm-model <model>`).
- [ ] Allow configuring temperature and system prompts.

---

### **Error Handling & Safety**
- [ ] Prevent AI modifications on untracked Git files (unless overridden).
- [ ] Display warnings when AI output exceeds token limits.
- [ ] Implement proper rollback in case of errors.
- [ ] Ensure commands fail gracefully when necessary.

---

## **Usage Specification**

### **1️⃣ Start a New Session**

```bash
ca start-session "Refactoring legacy code"
```

- Starts tracking AI modifications in `.ca_session.json`.

### **2️⃣ Apply AI Modifications**

#### **Default (Batch Mode)**

```bash
ca code "Convert functions to use named arguments" --files main.go utils.go
```

- Modifies **all specified files at once**.

#### **Per-File Mode**

```bash
ca code "Refactor functions" --files main.go utils.go --per-file
```

- **Processes each file separately** to avoid token limits.

#### **Dry Run (Preview Changes)**

```bash
ca code "Upgrade to Rails 7" --dry-run
```

- **Shows the AI-generated diff** but does NOT modify files.

#### **Auto-Commit After AI Changes**

```bash
ca code "Improve error handling" --commit
```

- **Commits** the AI-modified files after applying changes.

### **3️⃣ Review AI Changes**

```bash
ca review
```

- Displays session steps and modified files.
- Shows **Git commit history** (if tracking with Git).

### **4️⃣ Undo AI Modifications**

#### **Rollback a Specific Step**

```bash
ca rollback --step 2
```

- Restores files to **pre-step commit**.

#### **Replay a Step with a New Prompt**

```bash
ca replay-step --step 2 --prompt "Use named structs instead of tuples"
```

- Updates the step **without creating a new one**.

### **5️⃣ Generate AI Summaries**

#### **Default (On-Screen)**

```bash
ca code "Review changes for potential bugs" --summary
```

- **Shows an AI-generated analysis** without modifying files.

#### **Save Summary to a File**

```bash
ca code "Analyze security vulnerabilities" --summary --output security_review.txt
```

- Saves output to ``.

#### **Store Summary in Session**

```bash
ca code "Review for bugs" --summary --store-session
```

- Saves **analysis inside **``.

### **6️⃣ Commit & End the Session**

```bash
ca commit "Completed function refactor"
ca end-session
```

- Moves the session to historical storage.

---

## **LLM Integration (**``**)**

### **Configure LLM Defaults**

```bash
ca config llm --set model=llama-3 --set temperature=0.7
```

- Sets the LLM model and temperature globally.

### **List Available Models**

```bash
ca config llm --list
```

- Shows supported local/remote models.

### **Override LLM for a Single Run**

```bash
ca code "Optimize SQL queries" --llm-model gpt-4
```

- Uses `gpt-4` for this command **without changing global settings**.

---

## **Getting Started**

### **1️⃣ Install Code-Assistant**

```bash
go install github.com/chrisrob11/codeassistant@latest
```

#### **Or Clone & Build**

```bash
git clone https://github.com/chrisrob11/codeassistant.git
cd codeassistant
go build -o ca
mv ca /usr/local/bin/
```

### **2️⃣ Configure Your LLM**

```bash
export CA_LLM_MODEL="llama-3"
export CA_LLM_TEMPERATURE=0.7
```

### **3️⃣ Start Using It**

```bash
ca start-session "Refactoring project"
ca code "Improve function names" --per-file
ca review
ca commit "Refactored functions"
ca end-session
```

---

## **Final Thoughts**

💡 **This CLI is flexible, AI-powered, and integrates well with Git.**\
🚀 **It allows both direct file modification and AI-powered analysis.**\
🔧 **LLM settings are fully configurable via **``**.**

📌 **This is currently under development—track progress with the implementation checklist!** 🚀

