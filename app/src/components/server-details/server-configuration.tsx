import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface ServerConfigurationProps {
	serverId: string;
	region: string;
	ramPlan: string;
	difficulty: string;
	gamemode: string;
	pvpEnabled: boolean;
	onlineMode: boolean;
	worldSeed: string;
	createdAt: string;
}

export function ServerConfiguration({
	serverId,
	region,
	ramPlan,
	difficulty,
	gamemode,
	pvpEnabled,
	onlineMode,
	worldSeed,
	createdAt,
}: ServerConfigurationProps) {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Server Configuration</CardTitle>
				<CardDescription>Current server settings and specifications</CardDescription>
			</CardHeader>
			<CardContent>
				<div className="grid grid-cols-2 md:grid-cols-3 gap-6">
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">Server ID</h3>
						<p className="text-sm font-mono">{serverId}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">Region</h3>
						<p className="text-sm capitalize">{region}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">RAM Plan</h3>
						<p className="text-sm font-semibold">{ramPlan}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">Difficulty</h3>
						<p className="text-sm capitalize">{difficulty}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">Gamemode</h3>
						<p className="text-sm capitalize">{gamemode}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">PvP</h3>
						<p className="text-sm">{pvpEnabled ? "Enabled" : "Disabled"}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">Online Mode</h3>
						<p className="text-sm">{onlineMode ? "Enabled" : "Disabled"}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">World Seed</h3>
						<p className="text-sm font-mono">{worldSeed || "Random"}</p>
					</div>
					<div className="space-y-1">
						<h3 className="text-sm font-medium text-muted-foreground">Created</h3>
						<p className="text-sm">
							{new Date(createdAt).toLocaleDateString()}
						</p>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
