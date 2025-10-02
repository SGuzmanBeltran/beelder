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

func BuildBasicDockerfile (serverType string, memorySettings string) string {
	return fmt.Sprintf(BasicServerTemplate, serverType, memorySettings)
}
