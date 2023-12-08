```mermaid
    sequenceDiagram
        client->>+provider: redirect to provider oauth page
        provider->>provider: do login
        provider->>-client: return oauth token

        client->>+provider: validate oauth token
        provider->>-client: return oauth session token

        client->>+provider: get data with oauth session token
        provider->>-client: return data
```