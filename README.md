# Tide: The Tiny IDE

Tide is a lightweight, AI-powered terminal IDE designed for seamless human-computer collaboration. Inspired by advanced AI coding assistants like Claude Code, Tide leverages a ReAct agent to understand and execute complex development tasks, from writing and debugging code to managing projects and automating workflows.

## Vision

Our vision is to create a powerful, extensible, and intuitive development environment that lives in your terminal. We aim to combine the power of large language models with the flexibility of a command-line interface to create a truly unique and efficient development experience.

## Architecture

The future architecture of Tide is designed to be modular and extensible, allowing for the continuous integration of new tools and capabilities.

```mermaid
graph TD
    subgraph User Interface
        A[Interactive Terminal UI]
    end

    subgraph Core
        B[Main Loop]
        C[ReAct Agent]
        D[Tool Manager]
        K[Plugin System]
    end

    subgraph Tools
        E[File System]
        F[Code Executor]
        G[Version Control]
        H[Debugger]
        L[Web Search]
        M[Multi-Modal]
    end

    subgraph Backend
        I[Language Server Protocol]
        J["AI Model (Multi-LLM Support)"]
    end

    A --> B
    B --> C
    C --> D
    D --> E
    D --> F
    D --> G
    D --> H
    D --> L
    D --> M
    C --> J
    F --> I
    D --> K
```

## Roadmap: Building a Claude-Code-Level Assistant

To achieve our goal of creating a top-tier AI coding assistant, we are focusing on the following key areas. Our roadmap is heavily inspired by the capabilities of Claude Code, and we aim to implement similar features to provide a competitive and powerful tool.

- [ ] **Interactive Terminal UI:** Implement a more interactive and user-friendly terminal UI using a library like `tview` or `bubbletea`. This will provide a more IDE-like experience with features like syntax highlighting, auto-completion, and inline diagnostics.

- [ ] **LSP Integration:** Integrate with language servers via the Language Server Protocol (LSP). This will enable advanced code intelligence features, including:
    - Code completion
    - Go-to-definition
    - Hover information
    - Real-time diagnostics

- [ ] **Debugger Support:** Add a debugging tool that can interact with a debugger (e.g., Delve for Go). This will allow you to set breakpoints, inspect variables, and step through code without leaving the terminal.

- [ ] **Enhanced Version Control:** Enhance the version control tool to support more complex Git operations, such as interactive rebasing, cherry-picking, and managing pull requests.

- [ ] **Plugin System:** Develop a robust plugin system to allow for custom extensions and tools. This will enable the community to contribute to the Tide ecosystem and tailor it to their specific needs.

- [ ] **Web Search:** Implement a tool for searching the web to gather information, read documentation, and stay up-to-date with the latest technologies.

- [ ] **Multi-LLM Support:** Allow configuration of different LLMs (e.g., Claude, Gemini, GPT-4). This will give you the flexibility to choose the model that best suits your needs and preferences.

- [ ] **Enhanced Context Management:** Implement a more sophisticated context management system that maintains a persistent understanding of your project and conversation history. This will enable the agent to provide more relevant and accurate assistance.

- [ ] **Multi-modality Support:** Enable the agent to process and understand images and screenshots. This will be useful for tasks like:
    - Debugging UI issues from a screenshot.
    - Generating code from a wireframe.
    - Understanding diagrams and charts.

- [ ] **Automated Workflow Support:** Automate complex development workflows, such as:
    - Creating and reviewing pull requests.
    - Running and analyzing test suites.
    - Generating documentation.

- [ ] **Advanced Integrations:** Integrate with a wider range of developer tools, including:
    - SSH for remote development.
    - TMUX for session management.
    - Puppeteer for browser automation.

- [ ] **Role-Based Personas:** Implement different AI personas that can be activated for specific tasks. For example, you could have a "debugger" persona for finding and fixing bugs, or a "refactor" persona for improving code quality.