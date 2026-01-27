import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface ServerStatsGridProps {
	currentPlayers: number;
	maxPlayers: number;
	tps: number;
	uptime: string;
	worldSize: number;
}

export function ServerStatsGrid({ currentPlayers, maxPlayers, tps, uptime, worldSize }: ServerStatsGridProps) {
	return (
		<div className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-6">
			<Card>
				<CardHeader className="pb-2 md:pb-3">
					<CardDescription className="text-xs">Players Online</CardDescription>
					<CardTitle className="text-lg md:text-xl">
						{currentPlayers}/{maxPlayers}
					</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader className="pb-2 md:pb-3">
					<CardDescription className="text-xs">Performance</CardDescription>
					<CardTitle className="text-lg md:text-xl">
						{tps.toFixed(1)} TPS
					</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader className="pb-2 md:pb-3">
					<CardDescription className="text-xs">Uptime</CardDescription>
					<CardTitle className="text-lg md:text-xl">
						{uptime}
					</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader className="pb-2 md:pb-3">
					<CardDescription className="text-xs">World Size</CardDescription>
					<CardTitle className="text-lg md:text-xl">
						{worldSize.toFixed(1)} GB
					</CardTitle>
				</CardHeader>
			</Card>
		</div>
	);
}
