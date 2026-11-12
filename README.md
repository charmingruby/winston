# AIKit

Register configured providers once in a shared Genkit instance, then pass a model
reference when creating an agent or making a direct inference call. Embedders
keep a provider adapter for document/query-specific options.
There is no provider-specific configuration struct, enum, or model-name resolver
in AIKit.

Each provider implements the public `Provider` interface:

- `provider/google.go`: Google plugin authentication and retrieval embedding options.
- `provider/openai.go`: OpenAI plugin authentication and typed SDK references.
- `provider.go`: the interface, without provider-specific registration logic.

```go
// Import github.com/charmingruby/stan/internal/platform/aikit/provider.
// Empty keys use the SDK's environment defaults.
openai := provider.NewOpenAIProvider("")
google := provider.NewGoogleProvider("")
ak := aikit.New(ctx, google, openai)

agent, err := aikit.DefineAgent(
    ak, "atendimento", openai.Model("gpt-5.4-mini"), tools, prompt,
)
if err != nil {
    return err
}

embedder, err := aikit.NewEmbedder(ak, google, "gemini-embedding-001", 768)
if err != nil {
    return err
}
```

- Plugins are registered at initialization. Model references identify the
  registered provider by namespace; no provider argument or identity check is
  needed for agents and direct inference.
- An empty API key delegates to the plugin's environment defaults: Google uses
  `GEMINI_API_KEY` or `GOOGLE_API_KEY`; OpenAI uses `OPENAI_API_KEY`.
- Pass unqualified model names. Each provider creates references using its SDK;
  Genkit and the provider handle model resolution and availability.
- `ak.Infer(ctx, openai.Model("gpt-5.4-mini"), tools, opts...)` uses the same model
  reference API. SDK-created model references with typed configuration also work.
- Google embeddings distinguish retrieval documents from queries. OpenAI uses
  the same configuration for both. Choose embedding models supporting the
  requested dimensions.
- All providers, agents, and embedders share `ak.Genkit`, including its registry
  and Reflection/Dev UI resources. There is no per-provider runtime or cache.
- Genkit rejects duplicate plugin namespaces. Different credentials under the
  same plugin name need separate AIKit instances; they cannot silently replace
  the registered credentials. References resolve using the registered plugin's
  credentials, not the object that created the reference. A missing plugin or
  model is reported by Genkit during execution.
- Plugin initialization follows Genkit's panic-on-invalid-setup behavior; model
  and credential errors during inference are returned by the SDK.

## Direct inference

Use `Infer` for a direct model call without defining an agent, such as
summarization, classification, or extraction. It uses the providers already
registered in the shared AIKit instance:

```go
// Import github.com/firebase/genkit/go/ai.
openai := provider.NewOpenAIProvider("")
ak := aikit.New(ctx, openai)

response, err := ak.Infer(
    ctx,
    openai.Model("gpt-5.4-mini"),
    nil, // No tools for this call.
    ai.WithPrompt("Explain RAG in a few words."),
)
if err != nil {
    return err
}

fmt.Println(response.Text())
```

- `Infer` accepts a model reference, a slice of AIKit tools, and Genkit
  `ai.GenerateOption` values. Genkit resolves the provider from the model
  reference's namespace.
- It delegates generation and tool execution to `genkit.Generate`, with a
  default limit of eight turns. Caller options are applied after the defaults,
  so options such as `ai.WithMaxTurns(12)` can override them.
- It returns `*ai.ModelResponse`. Generation failures and unexpected finish
  reasons wrap `aikit.ErrInference`; `stop` and `interrupted` are accepted
  finish reasons.
- `Infer` does not persist conversation history automatically. Supply history
  with each call, or use an agent for conversations with session state.

## Adding a provider

Implement `Provider` outside this package:

```go
Plugin() api.Plugin
Model(name string) ai.ModelRef
Embedder(name string, dimensions int, query bool) (ai.EmbedderRef, error)
```

Return your configured plugin from `Plugin`.
Use the provider SDK's reference helpers in `Model` and `Embedder`. Those methods
must be safe for concurrent calls and return fresh request configuration.
Generation-only providers can return an error from `Embedder`.

Pass your implementation to `New(ctx, provider)`, then use its `Model` reference
in `DefineAgent` or `Infer`, or pass the provider to `NewEmbedder`. AIKit's generic
initialization needs no changes to support new providers.

Changing the embedding model changes the vector space. Re-embed existing
documents and align the vector store dimension before switching models.
