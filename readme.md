# Beelder
The Beelder platform is a project that tries to replicate a minecraft server platform, this is a learning project using Go to manage minecraft servers as a platform.


### Slices
#### First slice
In the first slice I desire that we have the creation server flow, the user is allow to create a server based on a configuration and the system report the progress of
creation.

**Configuration Options:**
- Player count (affects memory allocation)
- Server type (Paper 1.21.x initially)
- Difficulty (Peaceful, Easy, Normal, Hard)
- Server name/MOTD
- Online mode (official Minecraft accounts only)

**The progress states the server should report will be:**
- Created
- Creating
- Aborted
- Running
- Stopped

**User Experience:**
- Server connection details once ready
- Basic server management (start/stop/restart)

**Definition of Done:**
- User can create a server through API
- Server starts successfully and accepts connections
- Progress is tracked
- Basic server info is accessible (IP, port, status)