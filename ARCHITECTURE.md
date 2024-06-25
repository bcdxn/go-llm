# Architecture.md

The Core LLM Module is installed first. This is the module that a user directly interacts with on the command line. Plugins can be installed allowing a user to interface with specific models. Plugins are installed through the Core LLM Module

```mermaid
flowchart LR

classDef grouping fill:none,stroke:#999,stroke-width:2px

subgraph local
    CoreLlm((Core LLM Module))

    subgraph plugins
        OpenAI(OpenAI)
        Anthropic(Anthropic)
    end
end

subgraph models
    GPT4o([GPT-4o])
    GPT35latest([gpt-3.5-turbo-0125])
    GPT35([gpt-3.5-turbo5])

    claude3haiku([claude-3-haiku-20240307])
    claude35sonnet([claude-3-5-sonnet-20240620])
    claude3opus([claude-3-opus-20240229])
end

CoreLlm -->|HTTP| OpenAI
CoreLlm -->|HTTP| Anthropic

OpenAI -->|HTTPS| GPT4o
OpenAI -->|HTTPS| GPT35latest
OpenAI -->|HTTPS| GPT35

Anthropic -->|HTTPS| claude3haiku
Anthropic -->|HTTPS| claude35sonnet
Anthropic -->|HTTPS| claude3opus

class plugins grouping
class local grouping
class models grouping
```