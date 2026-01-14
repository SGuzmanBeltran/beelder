# ADR 001: JAR Manager Abstraction

## Context
Beelder needs to download Minecraft server JAR files from mcutils.com to build Docker images. Initially, JARs were bundled in the repository, but this doesn't scale and prevents version flexibility.

## Decision
Implement a `JarManager` interface with a `LocalJarManager` implementation that:
- Downloads JARs on-demand from mcutils.com
- Caches locally to avoid repeated downloads
- Provides an abstraction that can later support S3 or other backends

## Consequences

### Positive
- Flexible versioning - can support multiple Minecraft versions
- Reduced repository size
- Cloud-ready architecture for future AWS migration
- Each worker maintains its own cache (simple, no coordination needed)

### Negative
- First build per server type is slower (download time)
- Each worker node downloads separately (mitigated by local cache)
- Dependency on mcutils.com availability

### Future Considerations
When scaling to multiple workers in AWS:
- Implement `S3JarManager` to share cache across workers
- Reduces redundant downloads and storage costs