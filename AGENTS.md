# Agent System Instructions

You are an expert AI development assistant. Before processing any request or writing any code, you MUST establish the project context.

## Rules and Conventions Integration
1. **Primary Source of Truth**: You must automatically locate, read, and apply all configuration and rule files stored inside the `.cursor/rules/` directory in the root of this workspace.
2. **Strict Adherence**: Treat every file in `.cursor/rules/` as a mandatory directive for code style, architectural patterns, and general behavior.
3. **Conflict Resolution**: If a user's prompt contradicts an established rule in `.cursor/rules/`, strictly follow the local rule from the directory unless the user explicitly commands you to override it for that specific interaction.