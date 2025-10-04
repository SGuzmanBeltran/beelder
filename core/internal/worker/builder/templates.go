package builder

import "fmt"

const BasicServerTemplate = `FROM alpine:latest

# Install necessary packages (openjdk for Minecraft, bash, curl, etc.)
RUN apk add --no-cache openjdk21-jre bash curl

# Set working directory
WORKDIR /server

# Copy the Minecraft server jar
COPY assets/executables/%s.jar /server/server.jar

# Expose default Minecraft port
EXPOSE 25565

# Accept EULA by default
RUN echo "eula=true" > eula.txt

# Default command to start the Minecraft server
CMD ["java", %s, "-jar", "server.jar", "nogui"]
`

const ForgeServerTemplate = `FROM alpine:latest

# Install necessary packages (openjdk for Minecraft, bash, curl, etc.)
RUN apk add --no-cache openjdk21-jre bash curl

# Set working directory
WORKDIR /server

# Copy the Forge installer
COPY assets/executables/%s.jar /server/forge-installer.jar

# Accept EULA before installation
RUN echo "eula=true" > eula.txt

# Run Forge installer (this creates the server files)
RUN java -jar forge-installer.jar --installServer

# Expose default Minecraft port
EXPOSE 25565

# Start the server using the run.sh script created by Forge installer
CMD ["bash", "run.sh"]
`

func BuildBasicDockerfile(serverType string, memorySettings string) string {
	return fmt.Sprintf(BasicServerTemplate, serverType, memorySettings)
}

func BuildForgeDockerfile(serverType string) string {
	return fmt.Sprintf(ForgeServerTemplate, serverType)
}
