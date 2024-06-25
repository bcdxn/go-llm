# Go LLM

_A CLI tool for quickly interacting with various LLMs_

Note: This project was heavily inspired by [Simon Willison's blog post and project](https://simonwillison.net/2024/Jun/17/cli-language-models/).

## Installation

```sh
go install github.com/bcdxn/go-llm
```

```sh
llm -v
# 1.0.0-rc.1
```

## Installing Plugins

```sh
llm install openai
```

You can list the plugins that you have installed and see which one is currently in use:

```sh
llm list plugins
```

You can list the supported models available to you after installing plugins:

```sh
llm list models
```

## Using Models

You can select the model to use:

```sh
llm use gpt-3.5-turbo
```

## Send a prompt to the model

```sh
llm 'hello there'
# Hi! how can I help you today?
```

## Start an interactive session

```sh
llm
> 
```